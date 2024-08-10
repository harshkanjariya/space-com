package main

import (
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

func validateDataFormat(inputFormat Format, data []byte, disableLogs bool) (*MessageData, bool) {
	switch inputFormat {
	case AOS:
		if !disableLogs {
			println("validating AOS Format...")
		}
		return validateAOSFormat(data, disableLogs)
	case PUS_TM:
		return validatePUSTMFormat(data, disableLogs)
	case PUS_TC:
		return validatePUSTCFormat(data, disableLogs)
	case CCSDS:
		return validateCCSDSFormat(data, disableLogs)
	default:
		return nil, false
	}
}

func validateAOSFormat(data []byte, disableLogs bool) (*MessageData, bool) {
	const minHeaderLength = 2 + 1 + 4

	if len(data) < minHeaderLength {
		if !disableLogs {
			println("AOS header length is too short.")
		}
		return nil, false
	}

	header := data[:2+1]
	crcReceived := binary.BigEndian.Uint32(data[len(data)-4:])
	dataWithoutCRC := data[:len(data)-4]

	crcComputed := crc32.ChecksumIEEE(dataWithoutCRC)

	if crcReceived != crcComputed {
		if !disableLogs {
			println("AOS crc is invalid.")
		}
		return nil, false
	}

	if len(dataWithoutCRC) >= 2 && dataWithoutCRC[0] == 0x00 && dataWithoutCRC[1] == 0x7b && dataWithoutCRC[2] == 0x7c {
		dataWithoutCRC = dataWithoutCRC[3:]
	}

	if !disableLogs {
		printAOSFrame(header, dataWithoutCRC, crcReceived)
	}

	return &MessageData{RawData: dataWithoutCRC, Header: header}, true
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
	fmt.Printf("CRC (Hex): %08x\n", crc)
	fmt.Printf("Payload Data (Hex): %s\n", bytesToHex(data))
	fmt.Printf("Payload Data (ASCII): \n")
	println(payloadASCII)
}

func validatePUSTMFormat(data []byte, disableLogs bool) (*MessageData, bool) {
	const CRCLength = 2
	const minHeaderLength = 12

	if len(data) < minHeaderLength+CRCLength {
		if !disableLogs {
			fmt.Println("PUS TM data length is too short.")
		}
		return nil, false
	}

	header := data[:minHeaderLength]
	var crcReceived uint16
	var dataWithoutCRC []byte

	if CRCLength <= len(data) {
		crcReceived = binary.BigEndian.Uint16(data[len(data)-CRCLength:])
		dataWithoutCRC = data[:len(data)-CRCLength]
	} else {
		dataWithoutCRC = data
	}

	crcComputed := CRC16_CCITT_FUNC(dataWithoutCRC)

	if crcReceived != crcComputed {
		if !disableLogs {
			fmt.Printf("PUS TM CRC is invalid. Expected: %08x, Got: %08x\n", crcComputed, crcReceived)
		}
		return nil, false
	}
	dataWithoutCRC = dataWithoutCRC[len(header):]

	serviceID := header[6]
	subserviceID := header[7]

	if !disableLogs {
		printPUSTMFrame(serviceID, subserviceID, dataWithoutCRC, crcReceived)
	}

	return &MessageData{RawData: dataWithoutCRC, Header: header}, true
}

func printPUSTMFrame(serviceID, subserviceID byte, data []byte, crc uint16) {
	fmt.Printf("Service ID: %d\n", serviceID)
	fmt.Printf("Subservice ID: %d\n", subserviceID)
	fmt.Printf("CRC (Hex): %08x\n", crc)
	fmt.Printf("Payload Data (Hex): %s\n", bytesToHex(data))
	fmt.Printf("Payload Data (ASCII): \n")
	payloadASCII := bytesToASCII(data)
	println(payloadASCII)
}

func validatePUSTCFormat(data []byte, disableLogs bool) (*MessageData, bool) {
	CRCLength := 2
	const minHeaderLength = 10

	if len(data) < minHeaderLength+CRCLength {
		if !disableLogs {
			fmt.Println("PUSTC data length is too short.")
		}
		return nil, false
	}

	header := data[:minHeaderLength]
	crcReceived := binary.BigEndian.Uint16(data[len(data)-CRCLength:])
	dataWithoutCRC := data[:len(data)-CRCLength]

	crcComputed := CRC16_CCITT_FUNC(dataWithoutCRC)

	if crcReceived != crcComputed {
		if !disableLogs {
			fmt.Printf("PUSTC CRC is invalid. Expected: %08x, Got: %08x\n", crcComputed, crcReceived)
		}
		return nil, false
	}

	serviceID := header[0]
	subserviceID := header[1]

	if !disableLogs {
		printPUSTCFrame(serviceID, subserviceID, dataWithoutCRC[len(header):], crcReceived)
	}

	return &MessageData{RawData: dataWithoutCRC[len(header):], Header: header}, true
}

func printPUSTCFrame(serviceID, subserviceID byte, data []byte, crc uint16) {
	fmt.Printf("Service ID: %d\n", serviceID)
	fmt.Printf("Subservice ID: %d\n", subserviceID)
	fmt.Printf("CRC (Hex): %08x\n", crc)
	fmt.Printf("Payload Data (Hex): %s\n", bytesToHex(data))
	fmt.Printf("Payload Data (ASCII): \n")
	payloadASCII := bytesToASCII(data)
	println(payloadASCII)
}

func validateCCSDSFormat(data []byte, disableLogs bool) (*MessageData, bool) {
	const headerLength = 9

	if len(data) < headerLength+4 { // 4 bytes for CRC
		if !disableLogs {
			fmt.Println("CCSDS data length is too short.")
		}
		return nil, false
	}

	header := data[:headerLength]
	payload := data[headerLength : len(data)-4] // excluding CRC
	crcReceived := binary.BigEndian.Uint32(data[len(data)-4:])

	crcComputed := crc32.ChecksumIEEE(data[:len(data)-4])
	if crcReceived != crcComputed {
		if !disableLogs {
			fmt.Printf("CCSDS CRC is invalid. Expected: %08x, Got: %08x\n", crcComputed, crcReceived)
		}
		return nil, false
	}

	if !disableLogs {
		fmt.Println("CCSDS data is valid.")
		fmt.Printf("Payload Data (ASCII): %s\n", bytesToASCII(payload))
	}

	return &MessageData{RawData: payload, Header: header}, true
}
