import { Controller, Post } from "@nestjs/common";
import { ConverterService } from "./converter.service";

@Controller("converter")
export class ConverterController {
  constructor(
    private readonly converterService: ConverterService,
  ) {
  }

  @Post("generate")
  async generate() {
  }
}
