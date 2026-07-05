export type FilamentBrand = 'Generic' | 'Bambu';

export type FilamentPreset = {
	brand: FilamentBrand;
	name: string;
	trayType: string;
	trayInfoIdx: string;
	nozzleTempMin: number;
	nozzleTempMax: number;
};

// Filament presets mirroring OrcaSlicer's Bambu profiles. `trayInfoIdx` is the
// filament id the printer expects; `trayType` is the material family reported to
// the AMS; temperatures are the profile defaults.
export const FILAMENT_PRESETS: FilamentPreset[] = [
	// Generic
	{
		brand: 'Generic',
		name: 'PLA',
		trayType: 'PLA',
		trayInfoIdx: 'GFL99',
		nozzleTempMin: 190,
		nozzleTempMax: 240
	},
	{
		brand: 'Generic',
		name: 'PLA Silk',
		trayType: 'PLA',
		trayInfoIdx: 'GFL96',
		nozzleTempMin: 190,
		nozzleTempMax: 240
	},
	{
		brand: 'Generic',
		name: 'PLA-CF',
		trayType: 'PLA-CF',
		trayInfoIdx: 'GFL98',
		nozzleTempMin: 190,
		nozzleTempMax: 240
	},
	{
		brand: 'Generic',
		name: 'PLA High Speed',
		trayType: 'PLA',
		trayInfoIdx: 'GFL95',
		nozzleTempMin: 190,
		nozzleTempMax: 240
	},
	{
		brand: 'Generic',
		name: 'PETG',
		trayType: 'PETG',
		trayInfoIdx: 'GFG99',
		nozzleTempMin: 220,
		nozzleTempMax: 270
	},
	{
		brand: 'Generic',
		name: 'PETG-CF',
		trayType: 'PETG-CF',
		trayInfoIdx: 'GFG98',
		nozzleTempMin: 240,
		nozzleTempMax: 270
	},
	{
		brand: 'Generic',
		name: 'ABS',
		trayType: 'ABS',
		trayInfoIdx: 'GFB99',
		nozzleTempMin: 240,
		nozzleTempMax: 280
	},
	{
		brand: 'Generic',
		name: 'ASA',
		trayType: 'ASA',
		trayInfoIdx: 'GFB98',
		nozzleTempMin: 240,
		nozzleTempMax: 280
	},
	{
		brand: 'Generic',
		name: 'TPU',
		trayType: 'TPU',
		trayInfoIdx: 'GFU99',
		nozzleTempMin: 200,
		nozzleTempMax: 250
	},
	{
		brand: 'Generic',
		name: 'TPU for AMS',
		trayType: 'TPU',
		trayInfoIdx: 'GFU98',
		nozzleTempMin: 200,
		nozzleTempMax: 250
	},
	{
		brand: 'Generic',
		name: 'PC',
		trayType: 'PC',
		trayInfoIdx: 'GFC99',
		nozzleTempMin: 260,
		nozzleTempMax: 290
	},
	{
		brand: 'Generic',
		name: 'PA',
		trayType: 'PA',
		trayInfoIdx: 'GFN99',
		nozzleTempMin: 240,
		nozzleTempMax: 280
	},
	{
		brand: 'Generic',
		name: 'PA-CF',
		trayType: 'PA-CF',
		trayInfoIdx: 'GFN98',
		nozzleTempMin: 260,
		nozzleTempMax: 300
	},
	{
		brand: 'Generic',
		name: 'PVA',
		trayType: 'PVA',
		trayInfoIdx: 'GFS99',
		nozzleTempMin: 190,
		nozzleTempMax: 240
	},
	{
		brand: 'Generic',
		name: 'HIPS',
		trayType: 'HIPS',
		trayInfoIdx: 'GFS98',
		nozzleTempMin: 220,
		nozzleTempMax: 270
	},
	{
		brand: 'Generic',
		name: 'PP',
		trayType: 'PP',
		trayInfoIdx: 'GFP97',
		nozzleTempMin: 220,
		nozzleTempMax: 250
	},

	// Bambu
	{
		brand: 'Bambu',
		name: 'PLA Basic',
		trayType: 'PLA',
		trayInfoIdx: 'GFA00',
		nozzleTempMin: 190,
		nozzleTempMax: 240
	},
	{
		brand: 'Bambu',
		name: 'PLA Matte',
		trayType: 'PLA',
		trayInfoIdx: 'GFA01',
		nozzleTempMin: 190,
		nozzleTempMax: 240
	},
	{
		brand: 'Bambu',
		name: 'PLA Silk',
		trayType: 'PLA',
		trayInfoIdx: 'GFA05',
		nozzleTempMin: 190,
		nozzleTempMax: 240
	},
	{
		brand: 'Bambu',
		name: 'PLA Galaxy',
		trayType: 'PLA',
		trayInfoIdx: 'GFA15',
		nozzleTempMin: 190,
		nozzleTempMax: 240
	},
	{
		brand: 'Bambu',
		name: 'PLA Metal',
		trayType: 'PLA',
		trayInfoIdx: 'GFA02',
		nozzleTempMin: 190,
		nozzleTempMax: 240
	},
	{
		brand: 'Bambu',
		name: 'PLA-CF',
		trayType: 'PLA-CF',
		trayInfoIdx: 'GFA50',
		nozzleTempMin: 210,
		nozzleTempMax: 250
	},
	{
		brand: 'Bambu',
		name: 'PETG HF',
		trayType: 'PETG',
		trayInfoIdx: 'GFG02',
		nozzleTempMin: 230,
		nozzleTempMax: 270
	},
	{
		brand: 'Bambu',
		name: 'PETG Basic',
		trayType: 'PETG',
		trayInfoIdx: 'GFG00',
		nozzleTempMin: 230,
		nozzleTempMax: 270
	},
	{
		brand: 'Bambu',
		name: 'PETG Translucent',
		trayType: 'PETG',
		trayInfoIdx: 'GFG01',
		nozzleTempMin: 230,
		nozzleTempMax: 270
	},
	{
		brand: 'Bambu',
		name: 'PETG-CF',
		trayType: 'PETG-CF',
		trayInfoIdx: 'GFG50',
		nozzleTempMin: 240,
		nozzleTempMax: 270
	},
	{
		brand: 'Bambu',
		name: 'ABS',
		trayType: 'ABS',
		trayInfoIdx: 'GFB00',
		nozzleTempMin: 240,
		nozzleTempMax: 280
	},
	{
		brand: 'Bambu',
		name: 'ABS-GF',
		trayType: 'ABS-GF',
		trayInfoIdx: 'GFB50',
		nozzleTempMin: 240,
		nozzleTempMax: 280
	},
	{
		brand: 'Bambu',
		name: 'ASA',
		trayType: 'ASA',
		trayInfoIdx: 'GFB01',
		nozzleTempMin: 240,
		nozzleTempMax: 280
	},
	{
		brand: 'Bambu',
		name: 'ASA-CF',
		trayType: 'ASA-CF',
		trayInfoIdx: 'GFB51',
		nozzleTempMin: 250,
		nozzleTempMax: 280
	},
	{
		brand: 'Bambu',
		name: 'TPU 95A',
		trayType: 'TPU',
		trayInfoIdx: 'GFU01',
		nozzleTempMin: 200,
		nozzleTempMax: 250
	},
	{
		brand: 'Bambu',
		name: 'TPU 95A HF',
		trayType: 'TPU',
		trayInfoIdx: 'GFU00',
		nozzleTempMin: 200,
		nozzleTempMax: 250
	},
	{
		brand: 'Bambu',
		name: 'PC',
		trayType: 'PC',
		trayInfoIdx: 'GFC00',
		nozzleTempMin: 260,
		nozzleTempMax: 290
	},
	{
		brand: 'Bambu',
		name: 'PC FR',
		trayType: 'PC',
		trayInfoIdx: 'GFC01',
		nozzleTempMin: 260,
		nozzleTempMax: 290
	},
	{
		brand: 'Bambu',
		name: 'PA-CF',
		trayType: 'PA-CF',
		trayInfoIdx: 'GFN03',
		nozzleTempMin: 260,
		nozzleTempMax: 300
	},
	{
		brand: 'Bambu',
		name: 'PAHT-CF',
		trayType: 'PA-CF',
		trayInfoIdx: 'GFN04',
		nozzleTempMin: 260,
		nozzleTempMax: 300
	},
	{
		brand: 'Bambu',
		name: 'PA6-CF',
		trayType: 'PA6-CF',
		trayInfoIdx: 'GFN05',
		nozzleTempMin: 260,
		nozzleTempMax: 300
	},
	{
		brand: 'Bambu',
		name: 'PVA',
		trayType: 'PVA',
		trayInfoIdx: 'GFS04',
		nozzleTempMin: 210,
		nozzleTempMax: 250
	},
	{
		brand: 'Bambu',
		name: 'Support For PLA',
		trayType: 'PLA',
		trayInfoIdx: 'GFS02',
		nozzleTempMin: 190,
		nozzleTempMax: 240
	},
	{
		brand: 'Bambu',
		name: 'Support For PLA-PETG',
		trayType: 'PLA',
		trayInfoIdx: 'GFS05',
		nozzleTempMin: 190,
		nozzleTempMax: 240
	},
	{
		brand: 'Bambu',
		name: 'Support W',
		trayType: 'PLA',
		trayInfoIdx: 'GFS00',
		nozzleTempMin: 190,
		nozzleTempMax: 240
	}
];

export const FILAMENT_BRANDS: FilamentBrand[] = ['Generic', 'Bambu'];

// The filament color swatches offered by OrcaSlicer's AMS color picker.
export const FILAMENT_COLORS: string[] = [
	'#FFFFFF',
	'#FFF144',
	'#DCF478',
	'#0ACC38',
	'#057748',
	'#0D6284',
	'#0EE2A0',
	'#76D9F4',
	'#46A8F9',
	'#2850E0',
	'#443089',
	'#A03CF7',
	'#F330F9',
	'#D4B1DD',
	'#F95D73',
	'#F72323',
	'#7C4B00',
	'#F98C36',
	'#FCECD6',
	'#D3C5A3',
	'#AF7933',
	'#898989',
	'#BCBCBC',
	'#161616'
];
