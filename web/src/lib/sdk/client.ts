/**
 * OpenAPI
 * 0.0.1
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
export type Filament = {
    brand: string;
    name: string;
    nozzleTempMax: number;
    nozzleTempMin: number;
    trayInfoIdx: string;
    trayType: string;
};
export type HttpError = {
    /** Human readable error message */
    detail?: string;
    errors?: {
        /** Additional information about the error */
        more?: {
            [key: string]: any;
        };
        /** For example, name of the parameter that caused the error */
        name?: string;
        /** Human readable error message */
        reason?: string;
    }[];
    instance?: string;
    /** HTTP status code */
    status?: number;
    /** Short title of the error */
    title?: string;
    /** URL of the error type. Can be used to lookup the error in a documentation */
    "type"?: string;
};
export type Printer = {
    accessCode: string;
    hostIp: string;
    name: string;
    serial: string;
};
export type CreatePrinterDto = {
    accessCode: string;
    hostIp: string;
    name: string;
    serial: string;
};
export type UnknownInterface = any;
export type UpdatePrinterDto = {
    accessCode?: string;
    hostIp?: string;
    name?: string;
};
export type StartDryingDto = {
    coolingTemp?: number;
    duration: number;
    filament?: string;
    rotateTray?: boolean;
    temperature: number;
};
export type SetFanDto = {
    fan: string;
    speed?: number;
};
export type SetFilamentDto = {
    amsId?: number;
    nozzleTempMax: number;
    nozzleTempMin: number;
    trayColor: string;
    trayId?: number;
    trayInfoIdx: string;
    trayType: string;
};
export type SetLightDto = {
    state?: boolean;
};
export type SetPrintSpeedDto = {
    level: number;
};
export type PrinterStatus = {
    ams: {
        drying: boolean;
        dryingTime: number;
        humidity: number;
        id: number;
        supportsDrying: boolean;
        temperature: number;
        trays: {
            brand?: string;
            color?: string;
            empty: boolean;
            id: number;
            kValue?: number;
            loaded: boolean;
            material?: string;
            nozzleTempMax?: number;
            nozzleTempMin?: number;
            remaining?: number;
            trayInfoIdx?: string;
        }[];
    }[];
    buildPlate: {
        targetTemperature: number;
        temperature: number;
    };
    chamber: {
        controllable: boolean;
        targetTemperature: number;
        temperature: number;
    };
    chamberLight: boolean;
    externalSpool?: {
        brand?: string;
        color?: string;
        empty: boolean;
        id: number;
        kValue?: number;
        loaded: boolean;
        material?: string;
        nozzleTempMax?: number;
        nozzleTempMin?: number;
        remaining?: number;
        trayInfoIdx?: string;
    };
    fans: {
        aux: number;
        chamber: number;
        part: number;
    };
    fileName?: string;
    nozzle: {
        targetTemperature: number;
        temperature: number;
    };
    progress: number;
    speedLevel?: number;
    stage?: number;
    state: string;
    timeRemaining?: number;
};
/**
 * func1
 */
export function getFilaments({ accept }: {
    accept?: string;
} = {}, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 200;
        data: Filament[];
    } | {
        status: 400;
        data: HttpError;
    } | {
        status: 500;
        data: HttpError;
    } | {
        status: number;
    }>("/api/filament/", {
        ...opts,
        headers: oazapfts.mergeHeaders(opts?.headers, {
            Accept: accept
        })
    }));
}
/**
 * func1
 */
export function getPrinters({ accept }: {
    accept?: string;
} = {}, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 200;
        data: Printer[];
    } | {
        status: 400;
        data: HttpError;
    } | {
        status: 500;
        data: HttpError;
    } | {
        status: number;
    }>("/api/printer/", {
        ...opts,
        headers: oazapfts.mergeHeaders(opts?.headers, {
            Accept: accept
        })
    }));
}
/**
 * func3
 */
export function createPrinter(createPrinterDto: CreatePrinterDto, { accept }: {
    accept?: string;
} = {}, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 201;
        data: Printer;
    } | {
        status: 400;
        data: HttpError;
    } | {
        status: 500;
        data: HttpError;
    } | {
        status: number;
    }>("/api/printer/", oazapfts.json({
        ...opts,
        method: "PUT",
        body: createPrinterDto,
        headers: oazapfts.mergeHeaders(opts?.headers, {
            Accept: accept
        })
    })));
}
/**
 * func5
 */
export function deletePrinter(serial: string, { accept }: {
    accept?: string;
} = {}, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 204;
        data: UnknownInterface;
    } | {
        status: 400;
        data: HttpError;
    } | {
        status: 500;
        data: HttpError;
    } | {
        status: number;
    }>(`/api/printer/${encodeURIComponent(serial)}`, {
        ...opts,
        method: "DELETE",
        headers: oazapfts.mergeHeaders(opts?.headers, {
            Accept: accept
        })
    }));
}
/**
 * func4
 */
export function updatePrinter(serial: string, updatePrinterDto: UpdatePrinterDto, { accept }: {
    accept?: string;
} = {}, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 200;
        data: Printer;
    } | {
        status: 400;
        data: HttpError;
    } | {
        status: 500;
        data: HttpError;
    } | {
        status: number;
    }>(`/api/printer/${encodeURIComponent(serial)}`, oazapfts.json({
        ...opts,
        method: "PATCH",
        body: updatePrinterDto,
        headers: oazapfts.mergeHeaders(opts?.headers, {
            Accept: accept
        })
    })));
}
/**
 * func13
 */
export function startDrying(serial: string, amsId: string, startDryingDto: StartDryingDto, { accept }: {
    accept?: string;
} = {}, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 204;
        data: UnknownInterface;
    } | {
        status: 400;
        data: HttpError;
    } | {
        status: 500;
        data: HttpError;
    } | {
        status: number;
    }>(`/api/printer/${encodeURIComponent(serial)}/ams/${encodeURIComponent(amsId)}/dry`, oazapfts.json({
        ...opts,
        method: "POST",
        body: startDryingDto,
        headers: oazapfts.mergeHeaders(opts?.headers, {
            Accept: accept
        })
    })));
}
/**
 * func14
 */
export function stopDrying(serial: string, amsId: string, { accept }: {
    accept?: string;
} = {}, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 204;
        data: UnknownInterface;
    } | {
        status: 400;
        data: HttpError;
    } | {
        status: 500;
        data: HttpError;
    } | {
        status: number;
    }>(`/api/printer/${encodeURIComponent(serial)}/ams/${encodeURIComponent(amsId)}/dry/stop`, {
        ...opts,
        method: "POST",
        headers: oazapfts.mergeHeaders(opts?.headers, {
            Accept: accept
        })
    }));
}
/**
 * func11
 */
export function setFan(serial: string, setFanDto: SetFanDto, { accept }: {
    accept?: string;
} = {}, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 204;
        data: UnknownInterface;
    } | {
        status: 400;
        data: HttpError;
    } | {
        status: 500;
        data: HttpError;
    } | {
        status: number;
    }>(`/api/printer/${encodeURIComponent(serial)}/fan`, oazapfts.json({
        ...opts,
        method: "POST",
        body: setFanDto,
        headers: oazapfts.mergeHeaders(opts?.headers, {
            Accept: accept
        })
    })));
}
/**
 * func12
 */
export function setFilament(serial: string, setFilamentDto: SetFilamentDto, { accept }: {
    accept?: string;
} = {}, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 204;
        data: UnknownInterface;
    } | {
        status: 400;
        data: HttpError;
    } | {
        status: 500;
        data: HttpError;
    } | {
        status: number;
    }>(`/api/printer/${encodeURIComponent(serial)}/filament`, oazapfts.json({
        ...opts,
        method: "POST",
        body: setFilamentDto,
        headers: oazapfts.mergeHeaders(opts?.headers, {
            Accept: accept
        })
    })));
}
/**
 * func9
 */
export function setLight(serial: string, setLightDto: SetLightDto, { accept }: {
    accept?: string;
} = {}, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 204;
        data: UnknownInterface;
    } | {
        status: 400;
        data: HttpError;
    } | {
        status: 500;
        data: HttpError;
    } | {
        status: number;
    }>(`/api/printer/${encodeURIComponent(serial)}/light`, oazapfts.json({
        ...opts,
        method: "POST",
        body: setLightDto,
        headers: oazapfts.mergeHeaders(opts?.headers, {
            Accept: accept
        })
    })));
}
/**
 * func7
 */
export function pausePrint(serial: string, { accept }: {
    accept?: string;
} = {}, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 204;
        data: UnknownInterface;
    } | {
        status: 400;
        data: HttpError;
    } | {
        status: 500;
        data: HttpError;
    } | {
        status: number;
    }>(`/api/printer/${encodeURIComponent(serial)}/pause`, {
        ...opts,
        method: "POST",
        headers: oazapfts.mergeHeaders(opts?.headers, {
            Accept: accept
        })
    }));
}
/**
 * func8
 */
export function resumePrint(serial: string, { accept }: {
    accept?: string;
} = {}, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 204;
        data: UnknownInterface;
    } | {
        status: 400;
        data: HttpError;
    } | {
        status: 500;
        data: HttpError;
    } | {
        status: number;
    }>(`/api/printer/${encodeURIComponent(serial)}/resume`, {
        ...opts,
        method: "POST",
        headers: oazapfts.mergeHeaders(opts?.headers, {
            Accept: accept
        })
    }));
}
/**
 * func10
 */
export function setPrintSpeed(serial: string, setPrintSpeedDto: SetPrintSpeedDto, { accept }: {
    accept?: string;
} = {}, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 204;
        data: UnknownInterface;
    } | {
        status: 400;
        data: HttpError;
    } | {
        status: 500;
        data: HttpError;
    } | {
        status: number;
    }>(`/api/printer/${encodeURIComponent(serial)}/speed`, oazapfts.json({
        ...opts,
        method: "POST",
        body: setPrintSpeedDto,
        headers: oazapfts.mergeHeaders(opts?.headers, {
            Accept: accept
        })
    })));
}
/**
 * func2
 */
export function getPrinterStatus(serial: string, { accept }: {
    accept?: string;
} = {}, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 200;
        data: PrinterStatus;
    } | {
        status: 400;
        data: HttpError;
    } | {
        status: 500;
        data: HttpError;
    } | {
        status: number;
    }>(`/api/printer/${encodeURIComponent(serial)}/status`, {
        ...opts,
        headers: oazapfts.mergeHeaders(opts?.headers, {
            Accept: accept
        })
    }));
}
/**
 * func6
 */
export function stopPrint(serial: string, { accept }: {
    accept?: string;
} = {}, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 204;
        data: UnknownInterface;
    } | {
        status: 400;
        data: HttpError;
    } | {
        status: 500;
        data: HttpError;
    } | {
        status: number;
    }>(`/api/printer/${encodeURIComponent(serial)}/stop`, {
        ...opts,
        method: "POST",
        headers: oazapfts.mergeHeaders(opts?.headers, {
            Accept: accept
        })
    }));
}
/**
 * func15
 */
export function unloadMaterial(serial: string, amsId: string, { accept }: {
    accept?: string;
} = {}, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 204;
        data: UnknownInterface;
    } | {
        status: 400;
        data: HttpError;
    } | {
        status: 500;
        data: HttpError;
    } | {
        status: number;
    }>(`/api/printer/${encodeURIComponent(serial)}/unload/${encodeURIComponent(amsId)}`, {
        ...opts,
        method: "POST",
        headers: oazapfts.mergeHeaders(opts?.headers, {
            Accept: accept
        })
    }));
}
