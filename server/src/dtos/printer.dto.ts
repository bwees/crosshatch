import { createZodDto } from 'nestjs-zod';
import { z } from 'zod';

const CreatePrinterSchema = z.object({
  serial: z.string(),
  name: z.string(),
  hostIp: z.string(),
  accessCode: z.string(),
});

const PrinterSchema = CreatePrinterSchema;

export class CreatePrinterDto extends createZodDto(CreatePrinterSchema) {}
export class UpdatePrinterDto extends createZodDto(
  CreatePrinterSchema.partial().omit({ serial: true }),
) {}

export class PrinterDto extends createZodDto(PrinterSchema) {}
