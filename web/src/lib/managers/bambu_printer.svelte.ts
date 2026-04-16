import type { PrinterDto } from '$lib/sdk';

export class BambuPrinterManager {
	serial: string = $state('');
	name: string = $state('');
	access: string = $state('');

	constructor(public printer: PrinterDto) {
		this.serial = printer.serial;
	}
}
