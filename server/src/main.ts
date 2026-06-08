import { Logger } from '@nestjs/common';
import { NestFactory } from '@nestjs/core';
import { NestExpressApplication } from '@nestjs/platform-express';
import { SwaggerModule } from '@nestjs/swagger';
import { migrate } from 'drizzle-orm/better-sqlite3/migrator';
import type { NextFunction, Request, Response } from 'express';
import { IncomingMessage } from 'http';
import { join } from 'path';
import { Duplex } from 'stream';
import { AppModule } from './app.module';
import { createOpenApiDocument } from './openapi';
import { Go2RTCProxy } from './proxies/go2rtc.proxy';

async function bootstrap() {
  const logger = new Logger('Bootstrap');
  const app = await NestFactory.create<NestExpressApplication>(AppModule);
  app.enableCors();
  app.enableShutdownHooks();

  const openApiDoc = createOpenApiDocument(app);

  SwaggerModule.setup('api', app, openApiDoc);

  const webStaticPath = process.env.WEB_STATIC_PATH;
  if (webStaticPath) {
    const indexFile = join(webStaticPath, 'index.html');
    app.useStaticAssets(webStaticPath, { index: false });
    app.use((req: Request, res: Response, next: NextFunction) => {
      if (req.method !== 'GET' && req.method !== 'HEAD') return next();
      if (req.path.startsWith('/api')) return next();
      res.sendFile(indexFile);
    });
    logger.log(`Serving static web assets from ${webStaticPath}`);
  }

  const server = app.getHttpServer();
  const wsProxy = app.get(Go2RTCProxy);

  logger.log('Running database migrations...');

  migrate(app.get('db'), {
    migrationsFolder: join(__dirname, 'db', 'migrations'),
  });

  logger.log('Database migrations completed.');

  server.on('upgrade', (req: IncomingMessage, socket: Duplex, head: Buffer) => {
    if (req.url?.startsWith('/api/go2rtc')) {
      wsProxy.handleUpgrade(req, socket, head);
    }
  });

  logger.log('WebSocket proxy for Go2RTC initialized');

  const port = process.env.PORT ?? 3000;
  await app.listen(port);
  logger.log(`Server is listening on port ${port}`);
}

bootstrap().catch((error) => {
  new Logger('Bootstrap').error('Failed to start the server', error);
  process.exit(1);
});
