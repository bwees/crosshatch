import { Injectable } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { EventEmitter2 } from '@nestjs/event-emitter';
import { Go2RTCRepository } from '../repositories/go2rtc.repository';
import { PrinterRepository } from '../repositories/printer.repository';

@Injectable()
export class BaseService {
  constructor(
    // NestJS services
    protected readonly configService: ConfigService,
    protected readonly eventEmitter: EventEmitter2,

    protected printerRepository: PrinterRepository,
    protected go2rtcRepository: Go2RTCRepository,
  ) {}
}
