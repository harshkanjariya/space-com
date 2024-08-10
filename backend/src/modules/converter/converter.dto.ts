import { IsNotEmpty, IsString } from "class-validator";

export class EncodeDto {
  @IsNotEmpty()
  @IsString()
  outputFormat: string;

  @IsNotEmpty()
  @IsString()
  message: string;
}