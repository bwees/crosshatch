import { Logger } from '@nestjs/common';
import { NestFactory } from '@nestjs/core';
import { writeFile } from 'node:fs/promises';
import { resolve } from 'node:path';
import { AppModule } from './app.module';
import { createOpenApiDocument } from './openapi';

async function generate() {
  const logger = new Logger('OpenApiGenerator');
  const app = await NestFactory.create(AppModule, { logger: false });

  try {
    const document = createOpenApiDocument(app);
    const outputPath = resolve(process.cwd(), 'openapi.json');

    await writeFile(
      outputPath,
      `${JSON.stringify(document, null, 2)}\n`,
      'utf8',
    );
    logger.log(`Wrote OpenAPI spec to ${outputPath}`);
  } finally {
    await app.close();
  }
}

void generate();
