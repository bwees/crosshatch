import { getPrinters, type Printer, type PrinterStatus } from '$lib/sdk';
import { io, type Socket } from 'socket.io-client';
import { SvelteMap } from 'svelte/reactivity';
import { BambuPrinterManager } from './bambu_printer.svelte';

class PrinterManager {
	private stateSocket?: Socket;
	private initialized = false;
	private printerManagers: Map<string, BambuPrinterManager> = new Map();

	printers: Map<string, Printer> = new SvelteMap();
	printerState: Map<string, PrinterStatus> = new SvelteMap();

	async refreshPrinters() {
		const printers = await getPrinters();
		this.printers.clear();
		this.printerManagers.clear();
		for (const printer of printers) {
			this.printerManagers.set(printer.serial, new BambuPrinterManager(printer));
			this.printers.set(printer.serial, printer);
		}
	}

	async initialize() {
		if (this.initialized) return;
		this.initialized = true;

		await this.refreshPrinters();

		const socket = io({
			path: '/api/ws',
			transports: ['websocket']
		});
		this.stateSocket = socket;

		socket.on('connect', () => {
			console.log('Connected:', socket.id);
		});

		socket.on('printer.status', (data) => this.handleMqttReport(JSON.parse(data)));
	}

	private handleMqttReport(data: PrinterStatus & { serial: string }) {
		this.printerState.set(data.serial, data);
	}

	reset() {
		this.stateSocket?.disconnect();
		this.stateSocket = undefined;
		this.initialized = false;
		this.printers.clear();
		this.printerManagers.clear();
		this.printerState.clear();
	}
}

export const printerManager = new PrinterManager();
