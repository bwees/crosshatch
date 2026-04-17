import { Injectable, Logger } from '@nestjs/common';
import { OnEvent } from '@nestjs/event-emitter';
import {
  OnGatewayConnection,
  OnGatewayInit,
  WebSocketGateway,
  WebSocketServer,
} from '@nestjs/websockets';
import { Server, Socket } from 'socket.io';
import { PrinterStatusDto } from 'src/dtos/printer_status.dto';
import { PrinterService } from 'src/services/printer.service';

@Injectable()
@WebSocketGateway({ cors: true, path: '/api/ws' })
export class WebsocketGateway implements OnGatewayInit, OnGatewayConnection {
  private readonly logger = new Logger(WebsocketGateway.name);

  @WebSocketServer()
  server!: Server;

  constructor(private readonly printerService: PrinterService) {}

  afterInit(server: Server) {
    this.logger.log('WebSocket server initialized');
    this.server = server;
  }

  async handleConnection(client: Socket) {
    this.logger.log(`Client connected: ${client.id}`);
    const state = await this.printerService.getAllPrinterStates();

    for (const [serial, printerState] of Object.entries(state)) {
      if (!printerState) continue;
      const status = printerState;
      client.emit('printer.status', JSON.stringify({ ...status, serial }));
    }
  }

  @OnEvent('printer.status')
  handlePrinterStatus(payload: { serial: string } & PrinterStatusDto) {
    this.server.emit('printer.status', JSON.stringify(payload));
  }
}
