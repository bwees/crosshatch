import { createZodDto } from 'nestjs-zod';
import { z } from 'zod';
import { BambuPrintState } from './mqtt.dto';

const PrinterStatusSchema = z.object({
  status: z.union([
    z.literal('idle'),
    z.literal('printing'),
    z.literal('paused'),
    z.literal('error'),
  ]),
  progress: z.number().min(0).max(100),
  fileName: z.string().optional(),
  timeRemaining: z.number().optional(),
});

export class PrinterStatusDto extends createZodDto(PrinterStatusSchema) {}

export function statusFromMQTT(payload: BambuPrintState): PrinterStatusDto {
  return {
    status: 'idle',
    progress: payload.mc_percent,
    fileName: payload.file,
    timeRemaining: payload.mc_remaining_time,
  };
}
