package plugin

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
)

func convertData(outputFormat Format, input *MessageData) []byte {
	switch outputFormat {
	case CCSDS:
		response, err := generateCCSDSPacket(0x00, 0x00, 0x0000, input)
		if err != nil {
			println(err)
		}
		return response
	case AOS:
		response, err := generateAOSPacket(0x00, 0x00, input)
		if err != nil {
			println(err)
		}
		return response
	case PUS_TM:
		response, err := generatePUSTMPacket(0x00, 0x00, input)
		if err != nil {
			println(err)
		}
		return response
	case PUS_TC:
		response, err := generatePUSTCPacket(0x00, 0x00, input)
		if err != nil {
			println(err)
		}
		return response
	default:
		return []byte("Unsupported format")
	}
}

// generateAOSPacket generates an AOS packet
func generateAOSPacket(spacecraftID uint16, vcID byte, message *MessageData) ([]byte, error) {
	const headerLength = 3
	buffer := new(bytes.Buffer)

	header := make([]byte, headerLength)
	binary.BigEndian.PutUint16(header[:2], spacecraftID)
	header[2] = vcID

	buffer.Write(header)
	if _, err := buffer.Write(message.RawData); err != nil {
		return nil, err
	}

	crc := crc32.ChecksumIEEE(buffer.Bytes())
	if err := binary.Write(buffer, binary.BigEndian, crc); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// generatePUSTMPacket generates a PUS TM packet
func generatePUSTMPacket(serviceID, subserviceID byte, message *MessageData) ([]byte, error) {
	const headerLength = 12

	buffer := new(bytes.Buffer)

	header := make([]byte, headerLength)
	header[6] = serviceID
	header[7] = subserviceID

	buffer.Write(header)
	if _, err := buffer.Write(message.RawData); err != nil {
		return nil, err
	}

	crc := CRC16_CCITT_FUNC(buffer.Bytes())
	if err := binary.Write(buffer, binary.BigEndian, crc); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// generatePUSTCPacket generates a PUS TC packet
func generatePUSTCPacket(serviceID, subserviceID byte, message *MessageData) ([]byte, error) {
	const headerLength = 10

	buffer := new(bytes.Buffer)

	header := make([]byte, headerLength)
	header[0] = serviceID
	header[1] = subserviceID

	buffer.Write(header)
	if _, err := buffer.Write(message.RawData); err != nil {
		return nil, err
	}

	crc := CRC16_CCITT_FUNC(buffer.Bytes())
	if err := binary.Write(buffer, binary.BigEndian, crc); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// converter.go
func generateCCSDSPacket(packetType byte, apid byte, sequenceCount uint16, message *MessageData) ([]byte, error) {
	const headerLength = 6
	packetLength := uint16(len(message.RawData) + headerLength + 4)

	buffer := new(bytes.Buffer)

	header := make([]byte, headerLength)
	header[0] = 0x00
	header[1] = packetType
	header[2] = 0x00
	header[3] = 0x00
	binary.BigEndian.PutUint16(header[4:6], sequenceCount)

	buffer.Write(header)
	buffer.WriteByte(apid)
	binary.Write(buffer, binary.BigEndian, packetLength)
	if _, err := buffer.Write(message.RawData); err != nil {
		return nil, err
	}

	crc := crc32.ChecksumIEEE(buffer.Bytes())
	if err := binary.Write(buffer, binary.BigEndian, crc); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
