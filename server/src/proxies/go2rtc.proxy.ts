import {
  Injectable,
  Logger,
  OnApplicationShutdown,
  OnModuleInit,
} from '@nestjs/common';
import httpProxy from 'http-proxy';

@Injectable()
export class Go2RTCProxy implements OnModuleInit, OnApplicationShutdown {
  private readonly logger = new Logger(Go2RTCProxy.name);

  private proxy = httpProxy.createProxyServer({
    target: process.env.GO2RTC_WS_URL ?? 'ws://localhost:1984',
    ws: true,
    changeOrigin: true,
  });

  onModuleInit() {
    this.proxy.on('error', (_, __, socket) => {
      if (socket?.writable) socket.end();
    });
  }

  onApplicationShutdown(signal?: string) {
    this.logger.log(
      `Shutting down Go2RTC proxy${signal ? ` (signal: ${signal})` : ''}`,
    );
    this.proxy.close();
  }

  handleUpgrade(req: any, socket: any, head: any) {
    req.url = req.url.replace(/^\/api\/go2rtc/, '/api/ws');

    // Remove the Origin header to prevent CORS issues with go2rtc
    delete req.headers['origin'];

    socket.on('error', () => socket.destroy());
    this.proxy.ws(req, socket, head);
  }
}
