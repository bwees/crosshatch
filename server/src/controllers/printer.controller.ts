import {
  Body,
  Controller,
  Delete,
  Get,
  HttpCode,
  Param,
  Patch,
  Post,
  Put,
} from '@nestjs/common';
import { PrinterDto, SetLightDto, UpdatePrinterDto } from 'src/dtos/printer.dto';
import { PrinterService } from 'src/services/printer.service';

@Controller('api/printer')
export class PrinterController {
  constructor(private readonly printerService: PrinterService) {}

  @Get()
  getPrinters(): Promise<PrinterDto[]> {
    return this.printerService.getPrinters();
  }

  @Put()
  @HttpCode(201)
  createPrinter(@Body() dto: PrinterDto): Promise<PrinterDto> {
    return this.printerService.createPrinter(dto);
  }

  @Delete(':serial')
  @HttpCode(204)
  deletePrinter(@Param('serial') serial: string): Promise<void> {
    return this.printerService.deletePrinter(serial);
  }

  @Patch(':serial')
  updatePrinter(
    @Param('serial') serial: string,
    @Body() dto: UpdatePrinterDto,
  ): Promise<PrinterDto> {
    return this.printerService.updatePrinter(serial, dto);
  }

  @Post(':serial/stop')
  @HttpCode(204)
  stopPrint(@Param('serial') serial: string): Promise<void> {
    return this.printerService.stopPrint(serial);
  }

  @Post(':serial/pause')
  @HttpCode(204)
  pausePrint(@Param('serial') serial: string): Promise<void> {
    return this.printerService.pausePrint(serial);
  }

  @Post(':serial/resume')
  @HttpCode(204)
  resumePrint(@Param('serial') serial: string): Promise<void> {
    return this.printerService.resumePrint(serial);
  }

  @Post(':serial/light')
  @HttpCode(204)
  setLight(
    @Param('serial') serial: string,
    @Body() dto: SetLightDto,
  ): Promise<void> {
    return this.printerService.setLight(serial, dto.state);
  }

  @Post(':serial/unload/:amsId')
  @HttpCode(204)
  unloadMaterial(
    @Param('serial') serial: string,
    @Param('amsId') amsId: string,
  ): Promise<void> {
    return this.printerService.unloadMaterial(serial, parseInt(amsId));
  }
}
