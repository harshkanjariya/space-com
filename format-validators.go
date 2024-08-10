package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
)

func isValidFormat(format string) bool {
	switch Format(format) {
	case AOS, PUS_TM, PUS_TC, CCSDS:
		return true
	default:
		return false
	}
}

func validateDataFormat(inputFormat Format, data []byte) bool {
	switch inputFormat {
	case AOS:
		println("validating AOS Format...")
		return validateAOSFormat(data)
	case PUS_TM:
		return validatePUSTMFormat(data)
	case PUS_TC:
		return validatePUSTCFormat(data)
	case CCSDS:
		return validateCCSDSFormat(data)
	default:
		return false
	}
}

func validateAOSFormat(data []byte) bool {
	// Define minimum header length: spacecraft ID (2 bytes) + virtual channel ID (1 byte) + CRC (4 bytes)
	const minHeaderLength = 2 + 1 + 4

	// Check if data is long enough to contain at least the header and CRC
	if len(data) < minHeaderLength {
		println("AOS header length is too short.")
		return false
	}

	// Extract the header and CRC from the data
	header := data[:2+1]
	crcReceived := binary.BigEndian.Uint32(data[len(data)-4:])
	dataWithoutCRC := data[:len(data)-4]

	// Compute the CRC for the data (excluding the CRC itself)
	crcComputed := crc32.ChecksumIEEE(dataWithoutCRC)

	// Validate the computed CRC against the received CRC
	if crcReceived != crcComputed {
		println("AOS crc is invalid.")
		return false
	}

	if len(dataWithoutCRC) >= 2 && dataWithoutCRC[0] == 0x00 && dataWithoutCRC[1] == 0x7b && dataWithoutCRC[2] == 0x7c {
		dataWithoutCRC = dataWithoutCRC[3:]
	}

	printAOSFrame(header, dataWithoutCRC, crcReceived)

	return true
}

func printAOSFrame(header, data []byte, crc uint32) {
	if len(header) < 3 {
		fmt.Println("Error: Header length is too short")
		return
	}

	spacecraftID := binary.BigEndian.Uint16(header[:2])
	virtualChannelID := header[2]

	payloadASCII := bytesToASCII(data)
	fmt.Printf("Spacecraft ID: %d\n", spacecraftID)
	fmt.Printf("Virtual Channel ID: %d\n", virtualChannelID)
	fmt.Printf("Payload Data (Hex): %s\n", bytesToHex(data))
	fmt.Printf("Payload Data (ASCII): %s\n", payloadASCII)
	fmt.Printf("CRC (Hex): %08x\n", crc)
}

func validatePUSTMFormat(data []byte) bool {
	const CRCLength = 2
	const minHeaderLength = 12

	// Check if data length is sufficient
	if len(data) < minHeaderLength+CRCLength {
		fmt.Println("PUS TM data length is too short.")
		return false
	}

	// Extract header
	header := data[:minHeaderLength]

	// Initialize variables for CRC
	var crcReceived uint16
	var dataWithoutCRC []byte

	if CRCLength <= len(data) {
		// Extract CRC and data without CRC
		crcReceived = binary.BigEndian.Uint16(data[len(data)-CRCLength:])
		dataWithoutCRC = data[:len(data)-CRCLength]
	} else {
		// If no CRC length, the whole data is considered as payload
		dataWithoutCRC = data
	}

	// Calculate CRC over the header and payload (excluding CRC if present)
	crcComputed := CRC16_CCITT_FUNC(dataWithoutCRC)

	// Validate CRC if CRC length is non-zero
	if CRCLength > 0 && crcReceived != crcComputed {
		fmt.Printf("PUS TM CRC is invalid. Expected: %08x, Got: %08x\n", crcComputed, crcReceived)
		return false
	}
	dataWithoutCRC = dataWithoutCRC[len(header):]

	// Extract and print details from header
	serviceID := header[6]
	subserviceID := header[7]
	timestamp := header[2 : 2+1]

	printPUSTMFrame(serviceID, subserviceID, timestamp, dataWithoutCRC, crcReceived)

	return true
}

func printPUSTMFrame(serviceID, subserviceID byte, timestamp []byte, data []byte, crc uint16) {
	fmt.Printf("Service ID: %d\n", serviceID)
	fmt.Printf("Subservice ID: %d\n", subserviceID)
	fmt.Printf("Payload Data (Hex): %s\n", bytesToHex(data))
	fmt.Printf("Payload Data (ASCII): %s\n", bytesToASCII(data))
	fmt.Printf("CRC (Hex): %08x\n", crc)
}

func validatePUSTCFormat(data []byte) bool {
	CRCLength := 2
	// Calculate the minimum header length
	const minHeaderLength = 10

	// Check if the data is long enough to contain the header and CRC
	if len(data) < minHeaderLength+CRCLength {
		fmt.Println("PUSTC data length is too short.")
		return false
	}

	// Extract the header and CRC from the data
	header := data[:minHeaderLength]
	crcReceived := binary.BigEndian.Uint16(data[len(data)-CRCLength:])
	dataWithoutCRC := data[:len(data)-CRCLength]

	// Compute the CRC for the data (excluding the CRC itself)
	crcComputed := CRC16_CCITT_FUNC(dataWithoutCRC)

	// Validate the computed CRC against the received CRC
	if crcReceived != crcComputed {
		fmt.Printf("PUSTC CRC is invalid. Expected: %08x, Got: %08x\n", crcComputed, crcReceived)
		return false
	}

	// Extract fields from the header
	serviceID := header[0]
	subserviceID := header[1]

	// Print the parsed PUSTC packet details
	printPUSTCFrame(serviceID, subserviceID, dataWithoutCRC[len(header):], crcReceived)

	return true
}

func printPUSTCFrame(serviceID, subserviceID byte, data []byte, crc uint16) {
	fmt.Printf("Service ID: %d\n", serviceID)
	fmt.Printf("Subservice ID: %d\n", subserviceID)
	fmt.Printf("Payload Data (Hex): %s\n", bytesToHex(data))
	fmt.Printf("Payload Data (ASCII): %s\n", bytesToASCII(data))
	fmt.Printf("CRC (Hex): %08x\n", crc)
}

func validateCCSDSFormat(data []byte) bool {
	return bytes.HasPrefix(data, []byte("CCSDS:"))
}
