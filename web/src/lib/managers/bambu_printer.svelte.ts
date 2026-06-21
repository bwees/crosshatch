import type { Printer } from '$lib/sdk';

export class BambuPrinterManager {
	serial: string = $state('');
	name: string = $state('');
	access: string = $state('');

	constructor(public printer: Printer) {
		this.serial = printer.serial;
	}
}
