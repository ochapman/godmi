/*
* File Name:	type18_32bit_memory_error.go
* Description:
* Author:	Chapman Ou <ochapman.cn@gmail.com>
* Created:	2014-08-19
 */
package godmi

import (
	"fmt"
)

type MemoryErrorInformationType byte

const (
	MemoryErrorInformationTypeOther MemoryErrorInformationType = 1 + iota
	MemoryErrorInformationTypeUnknown
	MemoryErrorInformationTypeOK
	MemoryErrorInformationTypeBadread
	MemoryErrorInformationTypeParityerror
	MemoryErrorInformationTypeSingle_biterror
	MemoryErrorInformationTypeDouble_biterror
	MemoryErrorInformationTypeMulti_biterror
	MemoryErrorInformationTypeNibbleerror
	MemoryErrorInformationTypeChecksumerror
	MemoryErrorInformationTypeCRCerror
	MemoryErrorInformationTypeCorrectedsingle_biterror
	MemoryErrorInformationTypeCorrectederror
	MemoryErrorInformationTypeUncorrectableerror
)

func (m MemoryErrorInformationType) String() string {
	types := [...]string{
		"Other",
		"Unknown",
		"OK",
		"Bad read",
		"Parity error",
		"Single-bit error",
		"Double-bit error",
		"Multi-bit error",
		"Nibble error",
		"Checksum error",
		"CRC error",
		"Corrected single-bit error",
		"Corrected error",
		"Uncorrectable error",
	}
	return types[m-1]
}

type MemoryErrorInformationGranularity byte

const (
	MemoryErrorInformationGranularityOther MemoryErrorInformationGranularity = 1 + iota
	MemoryErrorInformationGranularityUnknown
	MemoryErrorInformationGranularityDevicelevel
	MemoryErrorInformationGranularityMemorypartitionlevel
)

func (m MemoryErrorInformationGranularity) String() string {
	grans := [...]string{
		"Other",
		"Unknown",
		"Device level",
		"Memory partition level",
	}
	return grans[m-1]
}

type MemoryErrorInformationOperation byte

const (
	MemoryErrorInformationOperationOther MemoryErrorInformationOperation = 1 + iota
	MemoryErrorInformationOperationUnknown
	MemoryErrorInformationOperationRead
	MemoryErrorInformationOperationWrite
	MemoryErrorInformationOperationPartialwrite
)

func (m MemoryErrorInformationOperation) String() string {
	operations := [...]string{
		"Other",
		"Unknown",
		"Read",
		"Write",
		"Partial write",
	}
	return operations[m-1]
}

type _32BitMemoryErrorInformation struct {
	infoCommon
	Type              MemoryErrorInformationType
	Granularity       MemoryErrorInformationGranularity
	Operation         MemoryErrorInformationOperation
	VendorSyndrome    uint32
	ArrayErrorAddress uint32
	ErrorAddress      uint32
	Resolution        uint32
}

func (m _32BitMemoryErrorInformation) String() string {
	return fmt.Sprintf("32 Bit Memory Error Information\n"+
		"\tError Type: %s\n"+
		"\tError Granularity: %s\n"+
		"\tError Operation: %s\n"+
		"\tVendor Syndrome: %d\n"+
		"\tMemory Array Error Address: %d\n"+
		"\tDevice Error Address: %d\n"+
		"\tError Resoluton: %d",
		m.Type,
		m.Granularity,
		m.Operation,
		m.VendorSyndrome,
		m.ArrayErrorAddress,
		m.ErrorAddress,
		m.Resolution,
	)
}

func new_32BitMemoryErrorInformation(h dmiHeader) dmiTyper {
	data := h.data
	return &_32BitMemoryErrorInformation{
		Type:              MemoryErrorInformationType(data[0x04]),
		Granularity:       MemoryErrorInformationGranularity(data[0x05]),
		Operation:         MemoryErrorInformationOperation(data[0x06]),
		VendorSyndrome:    u32(data[0x07:0x0B]),
		ArrayErrorAddress: u32(data[0x0B:0x0F]),
		ErrorAddress:      u32(data[0x0F:0x13]),
		Resolution:        u32(data[0x13:0x22]),
	}
}

func Get_32BitMemoryErrorInformation() *_32BitMemoryErrorInformation {
	if d, ok := gdmi[SMBIOSStructureType32_bitMemoryError]; ok {
		return d.(*_32BitMemoryErrorInformation)
	}
	return nil
}

func init() {
	addTypeFunc(SMBIOSStructureType32_bitMemoryError, new_32BitMemoryErrorInformation)
}
