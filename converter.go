package main

func convertData(inputFormat, outputFormat Format, input []byte) []byte {
	if inputFormat == outputFormat {
		return input
	}

	switch inputFormat {
	case AOS:
		return convertFromAOS(outputFormat, input)
	case PUS_TM:
		return convertFromTM(outputFormat, input)
	case CCSDS:
		return convertFromCCSDS(outputFormat, input)
	default:
		return []byte("Unsupported format")
	}
}

func convertFromAOS(outputFormat Format, input []byte) []byte {
	switch outputFormat {
	case PUS_TM:
		return convertAOSToTM(input)
	case CCSDS:
		return convertAOSToCCSDS(input)
	default:
		return input
	}
}

func convertFromTM(outputFormat Format, input []byte) []byte {
	switch outputFormat {
	case AOS:
		return convertTMToAOS(input)
	case CCSDS:
		return convertTMToCCSDS(input)
	default:
		return input
	}
}

func convertFromCCSDS(outputFormat Format, input []byte) []byte {
	switch outputFormat {
	case AOS:
		return convertCCSDSToAOS(input)
	case PUS_TM:
		return convertCCSDSToTM(input)
	default:
		return input
	}
}

func convertAOSToTM(input []byte) []byte {
	return append([]byte("AOS to TM: "), input...)
}

func convertAOSToCCSDS(input []byte) []byte {
	return append([]byte("AOS to CCSDS: "), input...)
}

func convertTMToAOS(input []byte) []byte {
	return append([]byte("TM to AOS: "), input...)
}

func convertTMToCCSDS(input []byte) []byte {
	return append([]byte("TM to CCSDS: "), input...)
}

func convertCCSDSToAOS(input []byte) []byte {
	return append([]byte("CCSDS to AOS: "), input...)
}

func convertCCSDSToTM(input []byte) []byte {
	return append([]byte("CCSDS to TM: "), input...)
}
