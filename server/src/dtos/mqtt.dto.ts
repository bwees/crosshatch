import { z } from 'zod';

const numericString = z
  .string()
  .trim()
  .regex(/^-?\d+(\.\d+)?$/);

const coerceNumber = z
  .union([z.number(), numericString])
  .transform((value) => (typeof value === 'string' ? Number(value) : value));

const coerceInt = coerceNumber.pipe(z.int());

const jsonValue: z.ZodType<unknown> = z.lazy(() =>
  z.union([
    z.string(),
    z.number(),
    z.boolean(),
    z.null(),
    z.array(jsonValue),
    z.record(z.string(), jsonValue),
  ]),
);

const looseRecordSchema = z.record(z.string(), jsonValue);

// --- Array item schemas (fields optional since individual items vary) ---

const BambuAmsTraySchema = z
  .object({
    bed_temp: coerceNumber.optional(),
    bed_temp_type: coerceInt.optional(),
    cali_idx: coerceInt.optional(),
    cols: z.array(z.string()).optional(),
    ctype: coerceInt.optional(),
    drying_temp: coerceNumber.optional(),
    drying_time: coerceNumber.optional(),
    id: coerceInt.optional(),
    nozzle_temp_max: coerceNumber.optional(),
    nozzle_temp_min: coerceNumber.optional(),
    remain: coerceInt.optional(),
    state: coerceInt,
    tag_uid: z.string().optional(),
    total_len: coerceInt.optional(),
    tray_color: z.string().optional(),
    tray_diameter: coerceNumber.optional(),
    tray_id_name: z.string().optional(),
    tray_info_idx: z.string().optional(),
    tray_sub_brands: z.string().optional(),
    tray_type: z.string().optional(),
    tray_uuid: z.string().optional(),
    tray_weight: coerceNumber.optional(),
    xcam_info: z.string().optional(),
  })
  .loose();

const BambuCareSchema = z
  .object({
    id: z.string().optional(),
    info: z.string().optional(),
  })
  .loose();

const BambuDeviceAirductPartSchema = z
  .object({
    func: z.number().optional(),
    id: z.number().optional(),
    range: z.number().optional(),
    state: z.number().optional(),
  })
  .loose();

const BambuDeviceAirductModeSchema = z
  .object({
    ctrl: z.array(z.number()).optional(),
    modeId: z.number().optional(),
    off: z.array(z.number()).optional(),
  })
  .loose();

const BambuDeviceNozzleInfoSchema = z
  .object({
    color_m: z.string().optional(),
    diameter: coerceNumber.optional(),
    fila_id: z.string().optional(),
    id: coerceInt.optional(),
    sn: z.string().optional(),
    stat: coerceInt.optional(),
    tm: coerceInt.optional(),
    type: z.string().optional(),
    wear: coerceNumber.optional(),
  })
  .loose();

const BambuDeviceExtruderInfoSchema = z
  .object({
    filam_bak: z.array(z.unknown()).optional(),
    hnow: coerceInt.optional(),
    hpre: coerceInt.optional(),
    htar: coerceInt.optional(),
    id: coerceInt.optional(),
    info: coerceInt.optional(),
    snow: coerceInt.optional(),
    spre: coerceInt.optional(),
    star: coerceInt.optional(),
    stat: coerceInt.optional(),
    temp: coerceNumber.optional(),
  })
  .loose();

const BambuNetInfoSchema = z
  .object({
    ip: coerceInt.optional(),
    mask: coerceInt.optional(),
  })
  .loose();

const BambuJobStageSchema = z
  .object({
    clock_in: z.boolean().optional(),
    color: z.array(z.string()).optional(),
    diameter: z.array(z.number()).optional(),
    est_time: coerceInt.optional(),
    heigh: coerceNumber.optional(),
    idx: coerceInt.optional(),
    platform: z.string().optional(),
    print_then: z.boolean().optional(),
    proc_list: z.array(z.unknown()).optional(),
    tool: z.array(z.string()).optional(),
    type: coerceInt.optional(),
  })
  .loose();

// --- Singleton sub-schemas (fields required per 1Hz payload) ---

const BambuAmsUnitSchema = z
  .object({
    humidity: coerceInt,
    humidity_raw: coerceInt,
    dry_time: coerceInt,
    temp: coerceNumber,
    id: coerceInt,
    tray: z.array(BambuAmsTraySchema),
  })
  .loose();

const BambuAmsSchema = z
  .object({
    ams: z.array(BambuAmsUnitSchema),
    ams_exist_bits: z.string(),
    ams_exist_bits_raw: z.string(),
    cali_id: coerceInt,
    cali_stat: coerceInt,
    insert_flag: z.boolean(),
    power_on_flag: z.boolean(),
    tray_exist_bits: z.string(),
    tray_is_bbl_bits: z.string(),
    tray_now: z.string(),
    tray_pre: z.string(),
    tray_read_done_bits: z.string(),
    tray_reading_bits: z.string(),
    tray_tar: z.string(),
    unbind_ams_stat: coerceInt,
    version: coerceInt,
  })
  .loose();

const BambuDeviceAirductSchema = z
  .object({
    modeCur: z.number(),
    modeFunc: z.number(),
    modeList: z.array(BambuDeviceAirductModeSchema),
    modeVisable: z.number(),
    parts: z.array(BambuDeviceAirductPartSchema),
    subFunc: z.number(),
    subMode: z.number(),
    subVisable: z.number(),
    version: z.number(),
  })
  .loose();

const BambuDeviceTempInfoSchema = z
  .object({
    temp: coerceNumber,
  })
  .loose();

const BambuDeviceSchema = z
  .object({
    airduct: BambuDeviceAirductSchema,
    bed: z
      .object({
        info: BambuDeviceTempInfoSchema,
        state: z.number(),
      })
      .loose(),
    bed_temp: coerceNumber,
    cam: z
      .object({
        laser: z
          .object({
            cond: coerceInt,
            state: coerceInt,
          })
          .loose(),
        timelapse_path: z.string(),
      })
      .loose(),
    ctc: z
      .object({
        info: BambuDeviceTempInfoSchema,
        state: coerceInt,
      })
      .loose(),
    ext_tool: z
      .object({
        calib: coerceInt,
        low_prec: z.boolean(),
        mount: coerceInt,
        mount_3d: coerceInt,
        th_temp: coerceInt,
        type: z.string(),
      })
      .loose(),
    extruder: z
      .object({
        info: z.array(BambuDeviceExtruderInfoSchema),
        state: coerceInt,
      })
      .loose(),
    fan: coerceInt,
    holder: z.unknown(),
    laser: z
      .object({
        power: coerceInt,
      })
      .loose(),
    nozzle: z
      .object({
        exist: coerceInt,
        info: z.array(BambuDeviceNozzleInfoSchema),
        src_id: coerceInt,
        state: coerceInt,
        tar_id: coerceInt,
      })
      .loose(),
    plate: z
      .object({
        base: coerceInt,
        cali2d_id: z.string(),
        cur_id: z.string(),
        mat: coerceInt,
        tar_id: z.string(),
      })
      .loose(),
    type: coerceInt,
  })
  .loose();

const BambuIpcamSchema = z
  .object({
    agora_service: z.string(),
    brtc_service: z.string(),
    bs_state: coerceInt,
    cap_pic_enable: z.string(),
    ipcam_dev: z.string(),
    ipcam_record: z.string(),
    laser_preview_res: coerceInt,
    mode_bits: coerceInt,
    resolution: z.string(),
    rtsp_url: z.string(),
    timelapse: z.string(),
    tl_store_hpd_type: coerceInt,
    tl_store_path_type: coerceInt,
    tutk_server: z.string(),
  })
  .loose();

const BambuJobSchema = z
  .object({
    cur_stage: z
      .object({
        idx: coerceInt,
        state: coerceInt,
      })
      .loose(),
    stage: z.array(BambuJobStageSchema),
  })
  .loose();

const BambuUpgradeStateSchema = z
  .object({
    ahb_new_version_number: z.string(),
    ams_new_version_number: z.string(),
    consistency_request: z.boolean(),
    dis_state: coerceInt,
    err_code: coerceInt,
    ext_new_version_number: z.string(),
    force_upgrade: z.boolean(),
    idx: coerceInt,
    idx2: coerceInt,
    lower_limit: z.string(),
    message: z.string(),
    module: z.string(),
    new_version_state: coerceInt,
    ota_new_version_number: z.string(),
    progress: coerceInt,
    sequence_id: coerceInt,
    sn: z.string(),
    status: z.string(),
  })
  .loose();

const BambuUploadSchema = z
  .object({
    file_size: coerceInt,
    finish_size: coerceInt,
    message: z.string(),
    oss_url: z.string(),
    progress: coerceInt,
    sequence_id: z.string(),
    speed: coerceInt,
    status: z.string(),
    task_id: z.string(),
    time_remaining: coerceInt,
    trouble_id: z.string(),
  })
  .loose();

const BambuXcamSchema = z
  .object({
    allow_skip_parts: z.boolean(),
    buildplate_marker_detector: z.boolean(),
    cfg: coerceInt,
    first_layer_inspector: z.boolean(),
    halt_print_sensitivity: z.string(),
    print_halt: z.boolean(),
    printing_monitor: z.boolean(),
    spaghetti_detector: z.boolean(),
  })
  .loose();

const BambuPrintInfoSchema = z
  .object({
    temp: z.union([coerceNumber, z.number()]),
  })
  .loose();

const BambuPrint3DSchema = z
  .object({
    layer_num: coerceInt,
    total_layer_num: coerceInt,
  })
  .passthrough();

// --- Main print state schema (1Hz fields are required) ---

export const BambuPrintStateSchema = z
  .object({
    '3D': BambuPrint3DSchema,
    ams: BambuAmsSchema,
    ams_rfid_status: coerceInt,
    ams_status: coerceInt,
    ap_err: coerceInt,
    aux: z.string(),
    aux_part_fan: z.boolean(),
    batch_id: coerceInt,
    bed_temper: coerceNumber,
    bed_target_temper: coerceNumber,
    big_fan1_speed: coerceInt,
    big_fan2_speed: coerceInt,
    cali_version: coerceInt,
    canvas_id: coerceInt,
    care: z.array(BambuCareSchema),
    cfg: z.string(),
    command: z.string(),
    cooling_fan_speed: coerceInt,
    design_id: z.string(),
    device: BambuDeviceSchema,
    err: z.string(),
    fail_reason: z.string(),
    fan_gear: coerceInt,
    file: z.string(),
    force_upgrade: z.boolean(),
    fun: z.string(),
    fun2: z.string(),
    gcode_file: z.string(),
    chamber_temper: coerceNumber.optional(),
    gcode_state: z.string(),
    heatbreak_fan_speed: coerceInt,
    hms: z.array(z.unknown()),
    home_flag: coerceInt,
    hw_switch_state: coerceInt,
    info: BambuPrintInfoSchema,
    ipcam: BambuIpcamSchema,
    job: BambuJobSchema,
    job_attr: coerceInt,
    job_id: z.string(),
    lan_task_id: z.string(),
    nozzle_temper: coerceNumber,
    nozzle_target_temper: coerceNumber,
    layer_num: coerceInt,
    lights_report: z.array(
      z
        .object({
          mode: z.string().optional(),
          node: z.string().optional(),
        })
        .loose(),
    ),
    mapping: z.array(z.number()),
    mc_action: coerceInt,
    mc_err: coerceInt,
    total_layer_num: coerceInt,
    mc_print_error_code: z.string(),
    mc_print_stage: coerceInt,
    mc_print_sub_stage: coerceInt,
    mc_stage: coerceInt,
    model_id: z.string(),
    msg: coerceInt,
    net: z
      .object({
        conf: coerceInt,
        info: z.array(BambuNetInfoSchema),
      })
      .loose(),
    gcode_file_prepare_percent: coerceInt,
    nozzle_diameter: coerceNumber,
    nozzle_type: z.string(),
    online: z
      .object({
        ahb: z.boolean(),
        version: coerceInt,
      })
      .loose(),
    percent: coerceInt,
    plate_cnt: coerceInt,
    plate_id: coerceInt,
    plate_idx: coerceInt,
    prepare_per: coerceInt,
    sequence_id: coerceInt,
    print_gcode_action: coerceInt,
    spd_mag: coerceInt,
    print_real_action: coerceInt,
    print_error: coerceInt,
    profile_id: z.string(),
    project_id: z.string(),
    queue: coerceInt,
    queue_est: coerceInt,
    queue_number: coerceInt,
    queue_sts: coerceInt,
    queue_total: coerceInt,
    remain_time: coerceInt,
    s_obj: z.array(z.unknown()),
    sdcard: z.boolean(),
    spd_lvl: coerceInt,
    stat: z.string(),
    state: coerceInt,
    stg: z.array(z.number()),
    stg_cd: coerceInt,
    stg_cur: coerceInt,
    subtask_id: z.string(),
    subtask_name: z.string(),
    support_chamber_temp_edit: z.boolean().optional(),
    task_id: z.string(),
    upgrade_state: BambuUpgradeStateSchema,
    upload: BambuUploadSchema,
    ver: z.string(),
    vir_slot: z.array(BambuAmsTraySchema.omit({ state: true })),
    xcam: BambuXcamSchema,
    xcam_status: z.string(),
    mc_percent: coerceInt,
    mc_remaining_time: coerceInt,
    print_type: z.string(),
    wifi_signal: z.string(),
  })
  .loose();

// Partial version for parsing incoming MQTT messages (all top-level fields optional)
const BambuPrintMessageSchema = BambuPrintStateSchema.partial();

// Outer wrapper for parsing raw MQTT payloads
export const BambuPrinterMessageSchema = z
  .object({
    print: BambuPrintMessageSchema.optional(),
    info: looseRecordSchema.optional(),
    system: looseRecordSchema.optional(),
    upgrade: looseRecordSchema.optional(),
    pushing: looseRecordSchema.optional(),
  })
  .loose();

export type BambuPrintState = z.infer<typeof BambuPrintStateSchema>;
export type BambuPrintMessage = z.infer<typeof BambuPrintMessageSchema>;
export type BambuPrinterMessage = z.infer<typeof BambuPrinterMessageSchema>;
