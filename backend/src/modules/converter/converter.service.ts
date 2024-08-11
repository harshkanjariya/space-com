import { Injectable, InternalServerErrorException } from "@nestjs/common";
import { exec, execFile } from "child_process";
import { promisify } from "util";
import { EncodeDto } from "./converter.dto";

const execPromise = promisify(exec);
const execFilePromise = promisify(execFile);

@Injectable()
export class ConverterService {
  async convert(inputFormat: string, outputFormat: string, data: string): Promise<string> {
    const allowedFormats = ["aos", "pus_tm", "pus_tc", "ccsds"];

    if (!allowedFormats.includes(inputFormat) || !allowedFormats.includes(outputFormat)) {
      throw new InternalServerErrorException("Invalid format");
    }

    try {
      const { stdout, stderr } = await execFilePromise("./converter.bin", ["-if", inputFormat, "-of", outputFormat, "-d", data]);

      if (!stdout) {
        throw new InternalServerErrorException(stderr);
      }

      return stdout;
    } catch (error) {
      console.log("Error converting: ", error);
      throw new InternalServerErrorException("Conversion failed");
    }
  }

  async encode(body: EncodeDto): Promise<string> {
    const allowedFormats = ["aos", "pus_tm", "pus_tc", "ccsds"];

    if (!allowedFormats.includes(body.outputFormat)) {
      throw new InternalServerErrorException("Invalid format");
    }

    try {
      const command = `./converter.bin -of=${body.outputFormat} -m="${body.message}"`;
      const { stdout, stderr } = await execPromise(command);

      if (stderr) {
        throw new InternalServerErrorException(stderr);
      }

      return stdout;
    } catch (error) {
      console.log("Error encoding: ", error);
      throw new InternalServerErrorException("Encoding failed");
    }
  }

  async decode(inputFormat: string, data: string): Promise<string> {
    const allowedFormats = ["aos", "pus_tm", "pus_tc", "ccsds"];

    if (!allowedFormats.includes(inputFormat)) {
      throw new InternalServerErrorException("Invalid format");
    }

    try {
      const { stdout, stderr } = await execFilePromise("./converter.bin", ["-if", inputFormat, "-d", data]);

      if (stderr) {
        throw new InternalServerErrorException(stderr);
      }

      return stdout;
    } catch (error) {
      console.log("Error decoding: ", error);
      throw new InternalServerErrorException("Decoding failed");
    }
  }
}
