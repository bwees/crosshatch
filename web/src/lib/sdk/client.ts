/**
 * Hatch API
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
export type CreatePrinterDto = {
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
export function getPrinters(opts?: Oazapfts.RequestOpts) {
    return oazapfts.fetchJson<{
        status: 200;
        data: PrinterDto[];
    }>("/api/printer", {
        ...opts
    });
}
export function createPrinter(createPrinterDto: CreatePrinterDto, opts?: Oazapfts.RequestOpts) {
    return oazapfts.fetchJson<{
        status: 201;
        data: PrinterDto;
    }>("/api/printer", oazapfts.json({
        ...opts,
        method: "PUT",
        body: createPrinterDto
    }));
}
export function deletePrinter(serial: string, opts?: Oazapfts.RequestOpts) {
    return oazapfts.fetchText(`/api/printer/${encodeURIComponent(serial)}`, {
        ...opts,
        method: "DELETE"
    });
}
export function updatePrinter(serial: string, updatePrinterDto: UpdatePrinterDto, opts?: Oazapfts.RequestOpts) {
    return oazapfts.fetchJson<{
        status: 200;
        data: PrinterDto;
    }>(`/api/printer/${encodeURIComponent(serial)}`, oazapfts.json({
        ...opts,
        method: "PATCH",
        body: updatePrinterDto
    }));
}
