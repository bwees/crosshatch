/**
 * Crosshatch API
 * 1.0.0
 * DO NOT MODIFY - This file has been generated using oazapfts.
 * See https://www.npmjs.com/package/oazapfts
 */
import * as Oazapfts from "@oazapfts/runtime";
import * as QS from "@oazapfts/runtime/query";
export const defaults: Oazapfts.Defaults<Oazapfts.CustomHeaders> = {
    headers: {},
    baseUrl: "/"
};
const oazapfts = Oazapfts.runtime(defaults);
export const servers = {};
export type PrinterDto = {
    serial: string;
    name: string;
    hostIp: string;
    accessCode: string;
};
export type UpdatePrinterDto = {
    name?: string;
    hostIp?: string;
    accessCode?: string;
};
export type SetLightDto = {
    state: boolean;
};
export type PrinterStatusDto = {
    state: "IDLE" | "RUNNING" | "PAUSE" | "PREPARE" | "SLICING" | "FINISH" | "FAILED";
    stage?: -1 | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12 | 13 | 14 | 15 | 16 | 17 | 18 | 19 | 20 | 21 | 22 | 23 | 24 | 25 | 26 | 27 | 28 | 29 | 30 | 31 | 32 | 33 | 34 | 35 | 36 | 37 | 38 | 39 | 40 | 41 | 42 | 43 | 44 | 45 | 46 | 47 | 48 | 49 | 50 | 51 | 52 | 53 | 54 | 55 | 56 | 57 | 58 | 65;
    progress: number;
    fileName?: string;
    timeRemaining?: number;
    buildPlate: {
        temperature: number;
        targetTemperature: number;
    };
    nozzle: {
        temperature: number;
        targetTemperature: number;
    };
    chamber: {
        temperature: number;
        targetTemperature: number;
        controllable: boolean;
    };
    ams: {
        id: number;
        humidity: number;
        temperature: number;
        dryingTime: number;
        trays: {
            id: number;
            empty: boolean;
            loaded: boolean;
            material?: string;
            brand?: string;
            color?: string;
            kValue?: number;
            nozzleTempMin?: number;
            nozzleTempMax?: number;
            remaining?: number;
        }[];
    }[];
    externalSpool?: {
        id: number;
        empty: boolean;
        loaded: boolean;
        material?: string;
        brand?: string;
        color?: string;
        kValue?: number;
        nozzleTempMin?: number;
        nozzleTempMax?: number;
        remaining?: number;
    };
    chamberLight: boolean;
};
export function getPrinters(opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 200;
        data: PrinterDto[];
    }>("/api/printer", {
        ...opts
    }));
}
export function createPrinter(printerDto: PrinterDto, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 201;
        data: PrinterDto;
    }>("/api/printer", oazapfts.json({
        ...opts,
        method: "PUT",
        body: printerDto
    })));
}
export function deletePrinter(serial: string, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchText(`/api/printer/${encodeURIComponent(serial)}`, {
        ...opts,
        method: "DELETE"
    }));
}
export function updatePrinter(serial: string, updatePrinterDto: UpdatePrinterDto, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 200;
        data: PrinterDto;
    }>(`/api/printer/${encodeURIComponent(serial)}`, oazapfts.json({
        ...opts,
        method: "PATCH",
        body: updatePrinterDto
    })));
}
export function stopPrint(serial: string, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchText(`/api/printer/${encodeURIComponent(serial)}/stop`, {
        ...opts,
        method: "POST"
    }));
}
export function pausePrint(serial: string, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchText(`/api/printer/${encodeURIComponent(serial)}/pause`, {
        ...opts,
        method: "POST"
    }));
}
export function resumePrint(serial: string, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchText(`/api/printer/${encodeURIComponent(serial)}/resume`, {
        ...opts,
        method: "POST"
    }));
}
export function setLight(serial: string, setLightDto: SetLightDto, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchText(`/api/printer/${encodeURIComponent(serial)}/light`, oazapfts.json({
        ...opts,
        method: "POST",
        body: setLightDto
    })));
}
export function unloadMaterial(serial: string, amsId: string, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchText(`/api/printer/${encodeURIComponent(serial)}/unload/${encodeURIComponent(amsId)}`, {
        ...opts,
        method: "POST"
    }));
}
