import { createZodDto } from 'nestjs-zod';
import { z } from 'zod';

const PrinterSchema = z.object({
  serial: z.string(),
  name: z.string(),
  hostIp: z.string(),
  accessCode: z.string(),
});

export class PrinterDto extends createZodDto(PrinterSchema) {}
export class UpdatePrinterDto extends createZodDto(
  PrinterSchema.partial().omit({ serial: true }),
) {}

const SetLightSchema = z.object({
  state: z.boolean(),
});

export class SetLightDto extends createZodDto(SetLightSchema) {}
