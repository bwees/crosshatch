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

  private state: BambuPrintState | null = null;

  // MQTT topic for receiving printer reports
  get reportTopic() {
    return `device/${this.connectionConfig.serial}/report`;
  }

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
}
