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
	inputFormat := flag.String("input-format", "ccsds", "Input format: aos, tm, or ccsds (default: ccsds)")
	outputFormat := flag.String("output-format", "ccsds", "Output format: aos, tm, or ccsds (default: ccsds)")
	hexInput := flag.String("data", "", "Data to be processed (required)")

	flag.StringVar(inputFormat, "if", "ccsds", "Short flag for --input-format")
	flag.StringVar(outputFormat, "of", "ccsds", "Short flag for --output-format")
	flag.StringVar(hexInput, "d", "", "Short flag for --data")

	flag.Parse()

	if *hexInput == "" {
		fmt.Println("Error: data must be provided")
		return
	}

	data, err := hex.DecodeString(*hexInput)
	if err != nil {
		log.Fatalf("Failed to decode hex string: %v", err)
	}

	if !isValidFormat(*inputFormat) {
		fmt.Println("Error: input-format must be one of: aos, tm, ccsds")
		return
	}

	if !isValidFormat(*outputFormat) {
		fmt.Println("Error: output-format must be one of: aos, tm, ccsds")
		return
	}

	if !validateDataFormat(Format(*inputFormat), data) {
		fmt.Println("Error: The data format does not match the specified input format")
		return
	}

	outputData := convertData(Format(*inputFormat), Format(*outputFormat), data)

	fmt.Println("Converted Data:")
	fmt.Println(outputData)
}
