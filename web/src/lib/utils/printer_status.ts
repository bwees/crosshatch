import type { PrinterStatusDto } from '$lib/sdk';

export function stateMessage(status: PrinterStatusDto['state'] | 'UNKNOWN') {
	switch (status) {
		case 'IDLE':
		case 'FINISH':
			return 'Idle';
		case 'RUNNING':
			return 'Printing';
		case 'PAUSE':
			return 'Paused';
		case 'PREPARE':
			return 'Preparing';
		case 'FAILED':
			return 'Failed';
		case 'SLICING':
			return 'Slicing';
		default:
			return 'Unknown';
	}
}

export function stateColor(status: PrinterStatusDto['state'] | 'UNKNOWN') {
	switch (status) {
		case 'IDLE':
		case 'FINISH':
			return `bg-green-500`;
		case 'RUNNING':
			return `bg-blue-500`;
		case 'PAUSE':
			return `bg-yellow-500`;
		case 'PREPARE':
			return `bg-purple-500`;
		case 'FAILED':
			return `bg-red-500`;
		case 'SLICING':
			return `bg-cyan-500`;
		default:
			return `bg-gray-600`;
	}
}

export function stageMessage(stage: number | null) {
	switch (stage) {
		case -1:
			return 'Idle';
		case 0:
			return 'Printing';
		case 1:
			return 'Bed Leveling';
		case 2:
			return 'Heatbed Preheating';
		case 3:
			return 'Vibration Compensation';
		case 4:
			return 'Changing Filament';
		case 5:
			return 'M400 Pause';
		case 6:
			return 'Pause: Filament Runout';
		case 7:
			return 'Heating Nozzle';
		case 8:
			return 'Dynamic Flow Calibration';
		case 9:
			return 'Bed Scan';
		case 10:
			return 'First Layer Inspection';
		case 11:
			return 'Build Plate Detection';
		case 12:
			return 'Lidar Calibration';
		case 13:
			return 'Homing';
		case 14:
			return 'Nozzle Cleaning';
		case 15:
			return 'Extruder Temp Check';
		case 16:
			return 'Pause: User';
		case 17:
			return 'Pause: Cover Fell';
		case 18:
			return 'Lidar Calibration 2';
		case 19:
			return 'Flow Ratio Calibration';
		case 20:
			return 'Pause: Nozzle Temp';
		case 21:
			return 'Pause: Heatbed Temp';
		case 22:
			return 'Filament Unloading';
		case 23:
			return 'Pause: Step Loss';
		case 24:
			return 'Filament Loading';
		case 25:
			return 'Motor Noise Cancellation';
		case 26:
			return 'Pause: AMS Offline';
		case 27:
			return 'Pause: Fan Speed';
		case 28:
			return 'Pause: Chamber Temp';
		case 29:
			return 'Cooling Chamber';
		case 30:
			return 'Pause: G-code User';
		case 31:
			return 'Motor Noise Showoff';
		case 32:
			return 'Pause: Nozzle Clumping';
		case 33:
			return 'Pause: Cutter Error';
		case 34:
			return 'Pause: First Layer Error';
		case 35:
			return 'Pause: Nozzle Clog';
		case 36:
			return 'Motion Precision Measurement';
		case 37:
			return 'Motion Precision Enhancement';
		case 38:
			return 'Motion Accuracy Measurement';
		case 39:
			return 'Nozzle Offset Calibration';
		case 40:
			return 'High Temp Bed Leveling';
		case 41:
			return 'Auto Check Lever';
		case 42:
			return 'Auto Check Door';
		case 43:
			return 'Laser Calibration';
		case 44:
			return 'Auto Check Platform';
		case 45:
			return 'Confirming BirdsEye Camera Location';
		case 46:
			return 'BirdsEye Calibration';
		case 47:
			return 'Bed Leveling Phase 1';
		case 48:
			return 'Bed Leveling Phase 2';
		case 49:
			return 'Heating Chamber';
		case 50:
			return 'Cooling Heatbed';
		case 51:
			return 'Printing Calibration Lines';
		case 52:
			return 'Auto Check Material';
		case 53:
			return 'Live View Calibration';
		case 54:
			return 'Waiting for Heatbed Temp';
		case 55:
			return 'Auto Check Material Position';
		case 56:
			return 'Cutting Module Calibration';
		case 57:
			return 'Surface Measurement';
		case 58:
			return 'Thermal Preconditioning';
		case 65:
			return 'Clumping Calibration';

		default:
			return 'Unknown Stage';
	}
}
