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

const PrinterStatusSchema = z.object({
  state: PrinterStateSchema,
  stage: PrinterStageSchema.optional(),
  progress: z.number().min(0).max(100),
  fileName: z.string().optional(),
  timeRemaining: z.number().optional(),
});

export class PrinterStatusDto extends createZodDto(PrinterStatusSchema) {}

export function statusFromMQTT(payload: BambuPrintState): PrinterStatusDto {
  return {
    state: PrinterStateSchema.parse(payload.gcode_state),
    stage: PrinterStageSchema.parse(payload.stg_cur),
    progress: payload.mc_percent,
    fileName: payload.file,
    timeRemaining: payload.mc_remaining_time,
  };
}
