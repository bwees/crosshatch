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
export type Printer = {
    accessCode: string;
    hostIp: string;
    name: string;
    serial: string;
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
export type SetLightDto = {
    state?: boolean;
};
export type PrinterStatus = {
    ams: {
        dryingTime: number;
        humidity: number;
        id: number;
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
    };
    fileName?: string;
    nozzle: {
        targetTemperature: number;
        temperature: number;
    };
    progress: number;
    stage?: number;
    state: string;
    timeRemaining?: number;
};
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
 * func10
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
