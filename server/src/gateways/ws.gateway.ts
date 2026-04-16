import { Logger } from '@nestjs/common';
import { OnEvent } from '@nestjs/event-emitter';
import {
  OnGatewayInit,
  WebSocketGateway,
  WebSocketServer,
} from '@nestjs/websockets';
import { Server } from 'socket.io';
import { PrinterStatusDto } from 'src/dtos/printer_status.dto';

@WebSocketGateway({ cors: true, path: '/api/ws' })
export class WebsocketGateway implements OnGatewayInit {
  private readonly logger = new Logger(WebsocketGateway.name);

  @WebSocketServer()
  server!: Server;

  afterInit(server: Server) {
    this.logger.log('WebSocket server initialized');
    this.server = server;
  }

  @OnEvent('printer.status')
  handlePrinterStatus(payload: { serial: string } & PrinterStatusDto) {
    this.server.emit('printer.status', JSON.stringify(payload));
  }
}
