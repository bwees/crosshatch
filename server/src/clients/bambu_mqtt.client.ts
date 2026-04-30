import { Logger } from '@nestjs/common';
import { mergeWith } from 'lodash';
import mqtt, { MqttClient } from 'mqtt';
import {
  BambuPrinterMessageSchema,
  BambuPrintMessage,
  BambuPrintState,
} from 'src/dtos/mqtt.dto';
import { MaybePromise } from 'src/utils/utils';

type PrinterConnectionConfig = {
  hostIp: string;
  accessCode: string;
  serial: string;
};

export class BambuMQTTClient {
  private connectionConfig: PrinterConnectionConfig;
  private client: MqttClient;
  private readonly logger: Logger;
  private onStatusUpdate: (
    serial: string,
    status: BambuPrintState,
  ) => MaybePromise<void>;

  state: BambuPrintState | null = null;

  // MQTT topic for receiving printer reports
  get reportTopic() {
    return `device/${this.connectionConfig.serial}/report`;
  }

  get commandTopic() {
    return `device/${this.connectionConfig.serial}/request`;
  }

  commandSequence = 0;

  constructor(
    connection: PrinterConnectionConfig,
    onStatusUpdate: (
      serial: string,
      status: BambuPrintState,
    ) => MaybePromise<void>,
    logger: Logger = new Logger(`${BambuMQTTClient.name}-${connection.serial}`),
  ) {
    this.logger = logger;
    this.connectionConfig = connection;
    this.onStatusUpdate = onStatusUpdate;

    this.client = mqtt.connect({
      host: connection.hostIp,
      username: 'bblp',
      password: connection.accessCode,
      reconnectPeriod: 5000,
      protocol: 'mqtts',
      port: 8883,
      keepalive: 60,
      rejectUnauthorized: false,
    });

    this.client.on('connect', () => {
      this.logger.log(`Connected to Bambu Printer at ${connection.hostIp}`);
      this.client.subscribe(this.reportTopic, (err) => {
        if (err) {
          this.logger.error(
            `Failed to subscribe to "${this.reportTopic}": ${err.message}`,
          );
        } else {
          this.logger.log(
            `Subscribed to topic "${this.reportTopic}" successfully`,
          );
        }
      });
    });

    this.client.on('message', (_, payload) => {
      const message = payload.toString();

      const jsonPayload = JSON.parse(message);
      const parsedPayload = BambuPrinterMessageSchema.safeParse(jsonPayload);

      if (!parsedPayload.success) {
        this.logger.warn(
          `Received invalid MQTT message on topic "${this.reportTopic}": ${parsedPayload.error}`,
        );
        return;
      }

      if (!('print' in parsedPayload.data)) {
        return;
      }

      this.state = this.mergePrintState(this.state, parsedPayload.data!.print!);
      void this.onStatusUpdate(this.connectionConfig.serial, this.state);
    });

    this.client.on('error', (err) => {
      this.logger.error(`MQTT client error: ${err.message}`);
    });
  }

  disconnect() {
    this.client.end(() => {
      this.logger.debug('MQTT client disconnected');
    });
  }

  private mergePrintState(
    currentState: BambuPrintState | null,
    nextMessage: BambuPrintMessage,
  ): BambuPrintState {
    const base = currentState ?? {};
    return mergeWith({}, base, nextMessage, (_, sourceValue) => {
      if (Array.isArray(sourceValue)) {
        return sourceValue;
      }

      return undefined;
    }) as BambuPrintState;
  }

  private async sendCommand(payload: Record<string, any>): Promise<void> {
    for (const key of ['print', 'info', 'system'] as const) {
      if (payload[key]) {
        payload[key].sequence_id = (this.commandSequence++).toString();
      }
    }

    this.client.publish(this.commandTopic, JSON.stringify(payload), (err) => {
      if (err) {
        this.logger.error('Failed to publish command to MQTT:', err);
      }
    });
  }

  // Control Methods

  async stopPrint(): Promise<void> {
    this.logger.debug(
      `Sending stop command to printer ${this.connectionConfig.serial}`,
    );

    await this.sendCommand({ print: { command: 'stop' } });
  }

  async pausePrint(): Promise<void> {
    this.logger.debug(
      `Sending pause command to printer ${this.connectionConfig.serial}`,
    );
    await this.sendCommand({ print: { command: 'pause' } });
  }

  async resumePrint(): Promise<void> {
    this.logger.debug(
      `Sending resume command to printer ${this.connectionConfig.serial}`,
    );
    await this.sendCommand({ print: { command: 'resume' } });
  }

  async setLight(on: boolean): Promise<void> {
    this.logger.debug(
      `Sending light ${on ? 'on' : 'off'} command to printer ${
        this.connectionConfig.serial
      }`,
    );

    await this.sendCommand({
      system: {
        command: 'ledctrl',
        led_node: 'chamber_light',
        led_mode: on ? 'on' : 'off', // "on" | "off" | "flashing"
        led_on_time: 500,
        led_off_time: 500,
        loop_times: 1,
        interval_time: 1000,
      },
    });
  }

  async unloadMaterial(amsID: number): Promise<void> {
    this.logger.debug(
      `Sending unload material command to printer ${this.connectionConfig.serial}`,
    );
    await this.sendCommand({
      print: {
        command: 'ams_change_filament',
        curr_temp: 255,
        tar_temp: 255,
        ams_id: amsID,
        target: 255,
        slot_id: 255,
      },
    });
  }
}
