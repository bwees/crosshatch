import { createZodDto } from 'nestjs-zod';
import { z } from 'zod';
import { BambuPrintState } from './mqtt.dto';

enum PrinterState {
  IDLE = 'IDLE',
  RUNNING = 'RUNNING',
  PAUSED = 'PAUSE',
  PREPARING = 'PREPARE',
  SLICING = 'SLICING',
  FINISHED = 'FINISH',
  FAILED = 'FAILED',
}

const PrinterStateSchema = z.enum(PrinterState);

export enum PrinterStage {
  IDLE = -1,
  PRINTING = 0,
  BED_LEVELING = 1,
  HEATBED_PREHEATING = 2,
  VIBRATION_COMPENSATION = 3,
  CHANGING_FILAMENT = 4,
  M400_PAUSE = 5,
  PAUSE_FILAMENT_RUNOUT = 6,
  HEATING_NOZZLE = 7,
  DYNAMIC_FLOW_CALIBRATION = 8,
  BED_SCAN = 9,
  FIRST_LAYER_INSPECTION = 10,
  BUILD_PLATE_DETECTION = 11,
  LIDAR_CALIBRATION = 12,
  HOMING = 13,
  NOZZLE_CLEANING = 14,
  EXTRUDER_TEMP_CHECK = 15,
  PAUSE_USER = 16,
  PAUSE_COVER_FELL = 17,
  LIDAR_CALIBRATION_2 = 18,
  FLOW_RATIO_CALIBRATION = 19,
  PAUSE_NOZZLE_TEMP = 20,
  PAUSE_HEATBED_TEMP = 21,
  FILAMENT_UNLOADING = 22,
  PAUSE_STEP_LOSS = 23,
  FILAMENT_LOADING = 24,
  MOTOR_NOISE_CANCELLATION = 25,
  PAUSE_AMS_OFFLINE = 26,
  PAUSE_FAN_SPEED = 27,
  PAUSE_CHAMBER_TEMP = 28,
  COOLING_CHAMBER = 29,
  PAUSE_GCODE_USER = 30,
  MOTOR_NOISE_SHOWOFF = 31,
  PAUSE_NOZZLE_CLUMPING = 32,
  PAUSE_CUTTER_ERROR = 33,
  PAUSE_FIRST_LAYER_ERROR = 34,
  PAUSE_NOZZLE_CLOG = 35,
  MOTION_PRECISION_MEASUREMENT = 36,
  MOTION_PRECISION_ENHANCEMENT = 37,
  MOTION_ACCURACY_MEASUREMENT = 38,
  NOZZLE_OFFSET_CALIBRATION = 39,
  HIGH_TEMP_BED_LEVELING = 40,
  AUTO_CHECK_LEVER = 41,
  AUTO_CHECK_DOOR = 42,
  LASER_CALIBRATION = 43,
  AUTO_CHECK_PLATFORM = 44,
  BIRDSEYE_CONFIRM_LOCATION = 45,
  BIRDSEYE_CALIBRATION = 46,
  BED_LEVELING_PHASE_1 = 47,
  BED_LEVELING_PHASE_2 = 48,
  HEATING_CHAMBER = 49,
  COOLING_HEATBED = 50,
  PRINTING_CALIBRATION_LINES = 51,
  AUTO_CHECK_MATERIAL = 52,
  LIVE_VIEW_CALIBRATION = 53,
  WAITING_HEATBED_TEMP = 54,
  AUTO_CHECK_MATERIAL_POSITION = 55,
  CUTTING_MODULE_CALIBRATION = 56,
  SURFACE_MEASUREMENT = 57,
  THERMAL_PRECONDITIONING = 58,
  CLUMPING_CALIBRATION = 65,
}

const PrinterStageSchema = z.enum(PrinterStage);

const AMSTraySchema = z.object({
  id: z.number(),
  empty: z.boolean(),
  loaded: z.boolean(),
  material: z.string().optional(),
  brand: z.string().optional(),
  color: z.string().optional(),
  kValue: z.number().optional(),
  nozzleTempMin: z.number().optional(),
  nozzleTempMax: z.number().optional(),
  remaining: z.number().optional(),
});

const AMSUnitSchema = z.object({
  id: z.number(),
  humidity: z.number(),
  temperature: z.number(),
  dryingTime: z.number(),
  trays: z.array(AMSTraySchema),
});

const PrinterStatusSchema = z.object({
  state: PrinterStateSchema,
  stage: PrinterStageSchema.optional(),
  progress: z.number().min(0).max(100),
  fileName: z.string().optional(),
  timeRemaining: z.number().optional(),
  buildPlate: z.object({
    temperature: z.number(),
    targetTemperature: z.number(),
  }),
  nozzle: z.object({
    temperature: z.number(),
    targetTemperature: z.number(),
  }),
  chamber: z.object({
    temperature: z.number(),
    targetTemperature: z.number(),
    controllable: z.boolean(),
  }),
  ams: z.array(AMSUnitSchema),
  externalSpool: AMSTraySchema.optional(),
  chamberLight: z.boolean(),
});

export class PrinterStatusDto extends createZodDto(PrinterStatusSchema) {}

export function statusFromMQTT(payload: BambuPrintState): PrinterStatusDto {
  // console.log(payload);
  const trayNow = payload.ams.tray_now;
  return {
    state: PrinterStateSchema.parse(payload.gcode_state),
    stage: PrinterStageSchema.parse(payload.stg_cur),
    progress: payload.mc_percent,
    fileName: payload.file,
    timeRemaining: payload.mc_remaining_time,
    buildPlate: {
      temperature: payload.bed_temper,
      targetTemperature: payload.bed_target_temper,
    },
    nozzle: {
      temperature: payload.nozzle_temper,
      targetTemperature: payload.nozzle_target_temper,
    },
    chamber: {
      temperature: payload.device.ctc.info.temp & 0xffff,
      targetTemperature: (payload.device.ctc.info.temp >> 16) & 0xffff,
      controllable: payload.support_chamber_temp_edit ?? false,
    },
    ams: payload.ams.ams.map((unit) => ({
      id: unit.id,
      humidity: unit.humidity_raw,
      dryingTime: unit.dry_time,
      temperature: unit.temp,
      trays: unit.tray.map((tray) => {
        const trayId = tray.id ?? 0;
        const globalId = unit.id * 4 + trayId;
        return {
          id: trayId,
          empty: !tray.tray_type,
          loaded: trayNow === String(globalId),
          material: tray.tray_type || undefined,
          brand: tray.tray_sub_brands || undefined,
          color: tray.tray_color || undefined,
          nozzleTempMin: tray.nozzle_temp_min,
          nozzleTempMax: tray.nozzle_temp_max,
          remaining: tray.remain,
        };
      }),
    })),
    externalSpool: payload.vir_slot?.[0]
      ? {
          id: payload.vir_slot[0].id ?? 254,
          empty: !payload.vir_slot[0].tray_type,
          loaded: trayNow === '254',
          material: payload.vir_slot[0].tray_type || undefined,
          brand: payload.vir_slot[0].tray_sub_brands || undefined,
          color: payload.vir_slot[0].tray_color || undefined,
          nozzleTempMin: payload.vir_slot[0].nozzle_temp_min,
          nozzleTempMax: payload.vir_slot[0].nozzle_temp_max,
          remaining: payload.vir_slot[0].remain,
        }
      : undefined,
    chamberLight:
      payload.lights_report.find((light) => light.node === 'chamber_light')
        ?.mode === 'on',
  };
}
