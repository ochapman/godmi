/*
* File Name:	type17_memory_device.go
* Description:
* Author:	Chapman Ou <ochapman.cn@gmail.com>
* Created:	2014-08-19
 */
package godmi

import (
	"fmt"
)

type MemoryDeviceFormFactor byte

const (
	MemoryDeviceFormFactorOther MemoryDeviceFormFactor = 1 + iota
	MemoryDeviceFormFactorUnknown
	MemoryDeviceFormFactorSIMM
	MemoryDeviceFormFactorSIP
	MemoryDeviceFormFactorChip
	MemoryDeviceFormFactorDIP
	MemoryDeviceFormFactorZIP
	MemoryDeviceFormFactorProprietaryCard
	MemoryDeviceFormFactorDIMM
	MemoryDeviceFormFactorTSOP
	MemoryDeviceFormFactorRowofchips
	MemoryDeviceFormFactorRIMM
	MemoryDeviceFormFactorSODIMM
	MemoryDeviceFormFactorSRIMM
	MemoryDeviceFormFactorFB_DIMM
)

func (m MemoryDeviceFormFactor) String() string {
	factors := [...]string{
		"Other",
		"Unknown",
		"SIMM",
		"SIP",
		"Chip",
		"DIP",
		"ZIP",
		"Proprietary Card",
		"DIMM",
		"TSOP",
		"Row of chips",
		"RIMM",
		"SODIMM",
		"SRIMM",
		"FB-DIMM",
	}
	return factors[m-1]
}

type MemoryDeviceType byte

const (
	MemoryDeviceTypeOther MemoryDeviceType = 1 + iota
	MemoryDeviceTypeUnknown
	MemoryDeviceTypeDRAM
	MemoryDeviceTypeEDRAM
	MemoryDeviceTypeVRAM
	MemoryDeviceTypeSRAM
	MemoryDeviceTypeRAM
	MemoryDeviceTypeROM
	MemoryDeviceTypeFLASH
	MemoryDeviceTypeEEPROM
	MemoryDeviceTypeFEPROM
	MemoryDeviceTypeEPROM
	MemoryDeviceTypeCDRAM
	MemoryDeviceType3DRAM
	MemoryDeviceTypeSDRAM
	MemoryDeviceTypeSGRAM
	MemoryDeviceTypeRDRAM
	MemoryDeviceTypeDDR
	MemoryDeviceTypeDDR2
	MemoryDeviceTypeDDR2FB_DIMM
	MemoryDeviceTypeReserved
	MemoryDeviceTypeDDR3
	MemoryDeviceTypeFBD2
)

func (m MemoryDeviceType) String() string {
	types := [...]string{
		"Other",
		"Unknown",
		"DRAM",
		"EDRAM",
		"VRAM",
		"SRAM",
		"RAM",
		"ROM",
		"FLASH",
		"EEPROM",
		"FEPROM",
		"EPROM",
		"CDRAM",
		"3DRAM",
		"SDRAM",
		"SGRAM",
		"RDRAM",
		"DDR",
		"DDR2",
		"DDR2 FB-DIMM",
		"Reserved",
		"DDR3",
		"FBD2",
	}
	return types[m-1]
}

type MemoryDeviceTypeDetail byte

const (
	MemoryDeviceTypeDetailReserved MemoryDeviceTypeDetail = 1 + iota
	MemoryDeviceTypeDetailOther
	MemoryDeviceTypeDetailUnknown
	MemoryDeviceTypeDetailFast_paged
	MemoryDeviceTypeDetailStaticcolumn
	MemoryDeviceTypeDetailPseudo_static
	MemoryDeviceTypeDetailRAMBUS
	MemoryDeviceTypeDetailSynchronous
	MemoryDeviceTypeDetailCMOS
	MemoryDeviceTypeDetailEDO
	MemoryDeviceTypeDetailWindowDRAM
	MemoryDeviceTypeDetailCacheDRAM
	MemoryDeviceTypeDetailNon_volatile
	MemoryDeviceTypeDetailRegisteredBuffered
	MemoryDeviceTypeDetailUnbufferedUnregistered
	MemoryDeviceTypeDetailLRDIMM
)

func (m MemoryDeviceTypeDetail) String() string {
	details := [...]string{
		"Reserved",
		"Other",
		"Unknown",
		"Fast-paged",
		"Static column",
		"Pseudo-static",
		"RAMBUS",
		"Synchronous",
		"CMOS",
		"EDO",
		"Window DRAM",
		"Cache DRAM",
		"Non-volatile",
		"Registered (Buffered)",
		"Unbuffered (Unregistered)",
		"LRDIMM",
	}
	return details[m-1]
}

type MemoryDevice struct {
	infoCommon
	PhysicalMemoryArrayHandle  uint16
	ErrorInformationHandle     uint16
	TotalWidth                 uint16
	DataWidth                  uint16
	Size                       uint16
	FormFactor                 MemoryDeviceFormFactor
	DeviceSet                  byte
	DeviceLocator              string
	BankLocator                string
	Type                       MemoryDeviceType
	TypeDetail                 MemoryDeviceTypeDetail
	Speed                      uint16
	Manufacturer               string
	SerialNumber               string
	AssetTag                   string
	PartNumber                 string
	Attributes                 byte
	ExtendedSize               uint32
	ConfiguredMemoryClockSpeed uint16
	MinimumVoltage             uint16
	MaximumVoltage             uint16
	ConfiguredVoltage          uint16
}

func (m MemoryDevice) String() string {
	return fmt.Sprintf("Memory Device\n"+
		"\tPhysical Memory Array Handle: %d\n"+
		"\tMemory Error Information Handle: %d\n"+
		"\tTotal Width: %d\n"+
		"\tData Width: %d\n"+
		"\tSize: %d\n"+
		"\tForm Factor: %s\n"+
		"\tDevice Set: %d\n"+
		"\tDevice Locator: %s\n"+
		"\tBank Locator: %s\n"+
		"\tMemory Type: %s\n"+
		"\tType Detail: %s\n"+
		"\tSpeed: %d\n"+
		"\tManufacturer: %s\n"+
		"\tSerial Number: %s\n"+
		"\tAsset Tag: %s\n"+
		"\tPart Number: %s\n"+
		"\tAttributes: %s\n"+
		"\tExtended Size: %s\n"+
		"\tConfigured Memory Clock Speed: %d\n"+
		"\tMinimum voltage: %d\n"+
		"\tMaximum voltage: %d\n"+
		"\tConfigured voltage: %d",
		m.PhysicalMemoryArrayHandle,
		m.ErrorInformationHandle,
		m.TotalWidth,
		m.DataWidth,
		m.Size,
		m.FormFactor,
		m.DeviceSet,
		m.DeviceLocator,
		m.BankLocator,
		m.Type,
		m.Speed,
		m.Manufacturer,
		m.SerialNumber,
		m.AssetTag,
		m.PartNumber,
		m.Attributes,
		m.ExtendedSize,
		m.ConfiguredMemoryClockSpeed,
		m.MinimumVoltage,
		m.MaximumVoltage,
		m.ConfiguredVoltage,
	)
}

func newMemoryDevice(h dmiHeader) dmiTyper {
	data := h.data
	return &MemoryDevice{
		PhysicalMemoryArrayHandle:  u16(data[0x04:0x06]),
		ErrorInformationHandle:     u16(data[0x06:0x08]),
		TotalWidth:                 u16(data[0x08:0x0A]),
		DataWidth:                  u16(data[0x0A:0x0C]),
		Size:                       u16(data[0x0C:0x0e]),
		FormFactor:                 MemoryDeviceFormFactor(data[0x0E]),
		DeviceSet:                  data[0x0F],
		DeviceLocator:              h.FieldString(int(data[0x10])),
		BankLocator:                h.FieldString(int(data[0x11])),
		Type:                       MemoryDeviceType(data[0x12]),
		TypeDetail:                 MemoryDeviceTypeDetail(u16(data[0x13:0x15])),
		Speed:                      u16(data[0x15:0x17]),
		Manufacturer:               h.FieldString(int(data[0x17])),
		SerialNumber:               h.FieldString(int(data[0x18])),
		PartNumber:                 h.FieldString(int(data[0x1A])),
		Attributes:                 data[0x1B],
		ExtendedSize:               u32(data[0x1C:0x20]),
		ConfiguredMemoryClockSpeed: u16(data[0x20:0x22]),
		MinimumVoltage:             u16(data[0x22:0x24]),
		MaximumVoltage:             u16(data[0x24:0x26]),
		ConfiguredVoltage:          u16(data[0x26:0x28]),
	}
}

func GetMemoryDevice() *MemoryDevice {
	if d, ok := gdmi[SMBIOSStructureTypeMemoryDevice]; ok {
		return d.(*MemoryDevice)
	}
	return nil
}

func init() {
	addTypeFunc(SMBIOSStructureTypeMemoryDevice, newMemoryDevice)
}
