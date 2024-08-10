import { Controller, Post, Body } from "@nestjs/common";
import { ConverterService } from "./converter.service";
import { EncodeDto } from "./converter.dto";

@Controller("converter")
export class ConverterController {
  constructor(
    private readonly converterService: ConverterService,
  ) {}

  @Post("convert")
  async convert(
    @Body() body: { inputFormat: string; outputFormat: string; data: string }
  ): Promise<string> {
    const { inputFormat, outputFormat, data } = body;
    return this.converterService.convert(inputFormat, outputFormat, data);
  }

  @Post("encode")
  async encode(
    @Body() body: EncodeDto
  ): Promise<string> {
    return this.converterService.encode(body);
  }

  @Post("decode")
  async decode(
    @Body() body: { inputFormat: string; data: string }
  ): Promise<string> {
    const { inputFormat, data } = body;
    return this.converterService.decode(inputFormat, data);
  }
}
