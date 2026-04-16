import { getPrinters, type PrinterDto, type PrinterStatusDto } from '$lib/sdk';
import { io, type Socket } from 'socket.io-client';
import { SvelteMap } from 'svelte/reactivity';
import { BambuPrinterManager } from './bambu_printer.svelte';

class PrinterManager {
	private stateSocket!: Socket;
	private printerManagers: Map<string, BambuPrinterManager> = new Map();

	printers: Map<string, PrinterDto> = new SvelteMap();
	printerState: Map<string, PrinterStatusDto> = new SvelteMap();

	async initialize() {
		const printers = await getPrinters();
		for (const printer of printers) {
			this.printerManagers.set(printer.serial, new BambuPrinterManager(printer));
			this.printers.set(printer.serial, printer);
		}

		this.stateSocket = io(`http://${window.location.host}`, {
			path: '/api/ws',
			transports: ['websocket']
		});

		this.stateSocket.on('connect', () => {
			console.log('Connected:', this.stateSocket.id);
		});

		this.stateSocket.on('printer.status', (data) => this.handleMqttReport(JSON.parse(data)));
	}

	private handleMqttReport(data: PrinterStatusDto & { serial: string }) {
		this.printerState.set(data.serial, data);
	}
}

export const printerManager = new PrinterManager();
