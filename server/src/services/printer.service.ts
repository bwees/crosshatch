import { Logger } from '@nestjs/common';
import { BambuMQTTClient } from 'src/clients/bambu_mqtt.client';
import { BambuPrintState } from 'src/dtos/mqtt.dto';
import { PrinterDto, UpdatePrinterDto } from 'src/dtos/printer.dto';
import { statusFromMQTT } from 'src/dtos/printer_status.dto';
import { streamURLForPrinter } from 'src/utils/utils';
import { BaseService } from './base.service';

export class PrinterService extends BaseService {
  private mqttClients: Map<string, BambuMQTTClient> = new Map();
  private readonly logger = new Logger(PrinterService.name);

  getPrinters(): Promise<PrinterDto[]> {
    return this.printerRepository.getPrinters();
  }

  async onPrinterStatusUpdate(serial: string, status: BambuPrintState) {
    const statusDto = statusFromMQTT(status);
    console.log(`Received status update for printer ${serial}:`, statusDto);
    this.eventEmitter.emit('printer.status', {
      serial,
      ...statusDto,
    });
  }

  async createPrinter(dto: PrinterDto) {
    const printer = (await this.printerRepository.createPrinter(dto))[0];

    const client = new BambuMQTTClient(
      printer,
      this.onPrinterStatusUpdate.bind(this),
    );
    this.mqttClients.set(printer.serial, client);

    await this.syncCameraStreams();

    return printer;
  }

  async deletePrinter(serial: string) {
    await this.printerRepository.deletePrinter(serial);

    const client = this.mqttClients.get(serial);
    if (client) {
      client.disconnect();
      this.mqttClients.delete(serial);
      this.logger.log(`Deleted printer ${serial} and disconnected MQTT client`);
    }

    await this.syncCameraStreams();
  }

  async updatePrinter(serial: string, dto: UpdatePrinterDto) {
    const existingPrinter =
      await this.printerRepository.getPrinterBySerial(serial);

    if (!existingPrinter) {
      throw new Error(`Printer with serial ${serial} not found`);
    }

    const updatedPrinter = await this.printerRepository.updatePrinter(
      serial,
      dto,
    );
    await this.syncCameraStreams();

    return updatedPrinter ?? existingPrinter;
  }

  private async syncCameraStreams() {
    const printers = await this.printerRepository.getPrinters();

    const desiredStreams = new Map<string, string>(
      printers.map((printer) => [printer.serial, streamURLForPrinter(printer)]),
    );

    // Add or update needed streams
    const existingStreams = await this.go2rtcRepository.getStreams();
    for (const [serial, url] of desiredStreams) {
      if (!existingStreams[serial]) {
        // Stream doesn't exist, create it
        await this.go2rtcRepository.createStream(serial, url);
      } else if (existingStreams[serial]?.producers[0].url !== url) {
        // Stream exists but URL is different, update it
        await this.go2rtcRepository.updateStream(serial, url);
      }

      this.logger.debug(`Installed stream for printer ${serial}`);
    }

    // Remove streams that are no longer needed
    for (const serial of Object.keys(existingStreams)) {
      if (!desiredStreams.has(serial)) {
        await this.go2rtcRepository.deleteStream(serial);

        this.logger.debug(`Deleted orphaned stream for printer ${serial}`);
      }
    }
  }

  async onModuleInit() {
    const printers = await this.printerRepository.getPrinters();

    for (const printer of printers) {
      const client = new BambuMQTTClient(
        printer,
        this.onPrinterStatusUpdate.bind(this),
      );
      this.mqttClients.set(printer.serial, client);
    }
  }

  async onModuleDestroy() {
    for (const client of this.mqttClients.values()) {
      client.disconnect();
    }
    this.mqttClients.clear();
  }
}
