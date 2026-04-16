import type { INestApplication } from '@nestjs/common';
import {
  DocumentBuilder,
  SwaggerDocumentOptions,
  SwaggerModule,
} from '@nestjs/swagger';
import { cleanupOpenApiDoc } from 'nestjs-zod';
import { PrinterStatusDto } from './dtos/printer_status.dto';

export function createOpenApiDocument(app: INestApplication) {
  const documentOptions: SwaggerDocumentOptions = {
    operationIdFactory: (_controllerKey: string, methodKey: string) =>
      methodKey,
    extraModels: [PrinterStatusDto],
  };

  const document = SwaggerModule.createDocument(
    app,
    new DocumentBuilder()
      .setTitle('Crosshatch API')
      .setDescription('Crosshatch server API')
      .setVersion('1.0.0')
      .build(),
    documentOptions,
  );

  return cleanupOpenApiDoc(document);
}
