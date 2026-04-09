import { Inject, Injectable } from '@nestjs/common';
import { eq } from 'drizzle-orm';
import { Printer, schema, type Database } from '../db/schema';

@Injectable()
export class PrinterRepository {
  constructor(@Inject('db') private db: Database) {}

  getPrinters(): Promise<Printer[]> {
    return this.db.select().from(schema.printer);
  }

  async getPrinterBySerial(serial: string): Promise<Printer | null> {
    const printers = await this.db
      .select()
      .from(schema.printer)
      .where(eq(schema.printer.serial, serial));

    return printers.length > 0 ? printers[0] : null;
  }

  createPrinter(data: {
    serial: string;
    name: string;
    hostIp: string;
    accessCode: string;
  }): Promise<Printer[]> {
    return this.db.insert(schema.printer).values(data).returning();
  }

  deletePrinter(serial: string): Promise<void> {
    return this.db
      .delete(schema.printer)
      .where(eq(schema.printer.serial, serial));
  }

  async updatePrinter(
    serial: string,
    data: {
      name?: string;
      hostIp?: string;
      accessCode?: string;
    },
  ): Promise<Printer | null> {
    const printers = await this.db
      .update(schema.printer)
      .set(data)
      .where(eq(schema.printer.serial, serial))
      .returning();

    return printers[0] ?? null;
  }
}
