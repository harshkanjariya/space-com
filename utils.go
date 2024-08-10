package main

import (
	"encoding/binary"
	"fmt"
)

func bytesToHex(data []byte) string {
	return fmt.Sprintf("%x", data)
}

func bytesToASCII(data []byte) string {
	return string(data)
}

func uint32ToUint16Slice(data []byte) ([]uint16, error) {
	// Ensure the data length is even, since we're converting pairs of bytes
	if len(data)%2 != 0 {
		return nil, fmt.Errorf("data length must be even to convert to uint16")
	}

	// Create a slice to hold the converted uint16 values
	uint16Slice := make([]uint16, len(data)/2)

	// Convert each pair of bytes into a uint16 value
	for i := 0; i < len(data); i += 2 {
		uint16Slice[i/2] = binary.BigEndian.Uint16(data[i : i+2])
	}

	return uint16Slice, nil
}

func CRC16_CCITT_FUNC(data []byte) uint16 {
	// CRC16-CCITT polynomial
	const polynomial = 0x1021
	var crc uint16 = 0xFFFF

	for _, b := range data {
		crc ^= uint16(b) << 8
		for i := 0; i < 8; i++ {
			if crc&0x8000 != 0 {
				crc = (crc << 1) ^ polynomial
			} else {
				crc <<= 1
			}
		}
	}
	return crc & 0xFFFF
}
