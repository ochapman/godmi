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
