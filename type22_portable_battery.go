/*
* File Name:	type22_portable_battery.go
* Description:
* Author:	Chapman Ou <ochapman.cn@gmail.com>
* Created:	2014-08-19
 */
package godmi

import (
	"fmt"
)

type PortableBatteryDeviceChemistry byte

const (
	PortableBatteryDeviceChemistryOther PortableBatteryDeviceChemistry = 1 + iota
	PortableBatteryDeviceChemistryUnknown
	PortableBatteryDeviceChemistryLeadAcid
	PortableBatteryDeviceChemistryNickelCadmium
	PortableBatteryDeviceChemistryNickelmetalhydride
	PortableBatteryDeviceChemistryLithium_ion
	PortableBatteryDeviceChemistryZincair
	PortableBatteryDeviceChemistryLithiumPolymer
)

func (p PortableBatteryDeviceChemistry) String() string {
	chems := [...]string{
		"Other",
		"Unknown",
		"Lead Acid",
		"Nickel Cadmium",
		"Nickel metal hydride",
		"Lithium-ion",
		"Zinc air",
		"Lithium Polymer",
	}
	return chems[p-1]
}

type PortableBattery struct {
	infoCommon
	Location                  string
	Manufacturer              string
	ManufacturerDate          string
	SerialNumber              string
	DeviceName                string
	DeviceChemistry           PortableBatteryDeviceChemistry
	DesignCapacity            uint16
	DesignVoltage             uint16
	SBDSVersionNumber         string
	MaximumErrorInBatteryData byte
	SBDSSerialNumber          uint16
	SBDSManufactureDate       uint16
	SBDSDeviceChemistry       string
	DesignCapacityMultiplier  byte
	OEMSepecific              uint32
}

func (p PortableBattery) String() string {
	return fmt.Sprintf("Portable Battery\n"+
		"\tLocation: %s\n"+
		"\tManufacturer: %s\n"+
		"\tManufacturer Date: %s\n"+
		"\tSerial Number: %s\n"+
		"\tDevice Name: %s\n"+
		"\tDevice Chemistry: %s\n"+
		"\tDesign Capacity: %d\n"+
		"\tDesign Voltage: %d\n"+
		"\tSBDS Version Number: %s\n"+
		"\tMaximum Error in Battery Data: %d\n"+
		"\tSBDS Serial Numberd: %d\n"+
		"\tSBDS Manufacturer Date: %d\n"+
		"\tSBDS Device Chemistry: %s\n"+
		"\tDesign Capacity Multiplier: %d\n"+
		"\tOEM-specific: %d",
		p.Location,
		p.Manufacturer,
		p.ManufacturerDate,
		p.SerialNumber,
		p.DeviceName,
		p.DeviceChemistry,
		p.DesignCapacity,
		p.DesignVoltage,
		p.SBDSVersionNumber,
		p.MaximumErrorInBatteryData,
		p.SBDSSerialNumber,
		p.SBDSManufactureDate,
		p.SBDSDeviceChemistry,
		p.DesignCapacityMultiplier,
		p.OEMSepecific,
	)
}

func newPortableBattery(h dmiHeader) dmiTyper {
	data := h.data
	return &PortableBattery{
		Location:                  h.FieldString(int(data[0x04])),
		Manufacturer:              h.FieldString(int(data[0x05])),
		ManufacturerDate:          h.FieldString(int(data[0x06])),
		SerialNumber:              h.FieldString(int(data[0x07])),
		DeviceName:                h.FieldString(int(data[0x08])),
		DeviceChemistry:           PortableBatteryDeviceChemistry(data[0x09]),
		DesignCapacity:            u16(data[0x0A:0x0C]),
		DesignVoltage:             u16(data[0x0C:0x0E]),
		SBDSVersionNumber:         h.FieldString(int(data[0x0E])),
		MaximumErrorInBatteryData: data[0x0F],
		SBDSSerialNumber:          u16(data[0x10:0x12]),
		SBDSManufactureDate:       u16(data[0x12:0x14]),
		SBDSDeviceChemistry:       h.FieldString(int(data[0x14])),
		DesignCapacityMultiplier:  data[0x15],
		OEMSepecific:              u32(data[0x16:0x1A]),
	}
}

func GetPortableBattery() *PortableBattery {
	if d, ok := gdmi[SMBIOSStructureTypePortableBattery]; ok {
		return d.(*PortableBattery)
	}
	return nil
}

func init() {
	addTypeFunc(SMBIOSStructureTypePortableBattery, newPortableBattery)
}
