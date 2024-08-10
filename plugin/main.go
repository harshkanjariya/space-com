package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
)

type Format string

const (
	AOS    Format = "aos"
	PUS_TM Format = "pus_tm"
	PUS_TC Format = "pus_tc"
	CCSDS  Format = "ccsds"
)

func main() {
	inputFormat := flag.String("input-format", "", "Input format: aos, pus_tm, pus_tc, or ccsds (required for format conversion)")
	outputFormat := flag.String("output-format", "", "Output format: aos, pus_tm, pus_tc, or ccsds (required)")
	hexInput := flag.String("data", "", "Hex data to be processed (required for format conversion)")
	message := flag.String("message", "", "String message to be encoded (required for message conversion)")

	flag.StringVar(inputFormat, "if", "", "Short flag for --input-format")
	flag.StringVar(outputFormat, "of", "", "Short flag for --output-format")
	flag.StringVar(hexInput, "d", "", "Short flag for --data")
	flag.StringVar(message, "m", "", "Short flag for --message")

	flag.Parse()

	if *inputFormat != "" && *hexInput != "" && *outputFormat == "" && *message == "" {
		data, err := hex.DecodeString(*hexInput)
		if err != nil {
			log.Fatalf("Failed to decode hex string: %v", err)
		}

		if !isValidFormat(*inputFormat) {
			fmt.Println("Error: input-format must be one of: aos, pus_tm, pus_tc, ccsds")
			return
		}

		payload, isValid := validateDataFormat(Format(*inputFormat), data, true)
		if !isValid {
			fmt.Println("Invalid input")
			return
		}
		fmt.Println(bytesToASCII(payload.RawData))
		return
	}

	// Case 1: Message conversion to a specific output format
	if *message != "" && *inputFormat == "" && *hexInput == "" {
		data := []byte(*message)
		outputData := convertData(Format(*outputFormat), &MessageData{RawData: data})

		// Print the converted data as a hex string
		fmt.Println("Converted Data:")
		fmt.Println(bytesToHex(outputData))
		return
	}

	// Case 2: Format conversion from input format to output format
	if *hexInput != "" && *inputFormat != "" {
		data, err := hex.DecodeString(*hexInput)
		if err != nil {
			log.Fatalf("Failed to decode hex string: %v", err)
		}

		// Validate formats
		if !isValidFormat(*inputFormat) {
			fmt.Println("Error: input-format must be one of: aos, pus_tm, pus_tc, ccsds")
			return
		}

		if !isValidFormat(*outputFormat) {
			fmt.Println("Error: output-format must be one of: aos, pus_tm, pus_tc, ccsds")
			return
		}

		// Validate data format (Implement validateDataFormat according to your needs)
		payload, isValid := validateDataFormat(Format(*inputFormat), data, false)
		if !isValid {
			fmt.Println("Error: The data format does not match the specified input format")
			return
		}

		// Convert the data using the appropriate generator function
		outputData := convertData(Format(*outputFormat), payload)

		// Print the converted data as a hex string
		fmt.Println("Converted Data:")
		fmt.Println(bytesToHex(outputData))
		return
	}

	// If neither valid combination is provided
	fmt.Println("Error: Invalid combination of flags provided. Use either:")
	fmt.Println("1. --message and --output-format for message conversion")
	fmt.Println("2. --input-format, --output-format, and --data for format conversion")
}
