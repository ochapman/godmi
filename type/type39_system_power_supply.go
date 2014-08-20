/*
* File Name:	type39_system_power_supply.go
* Description:	
* Author:	Chapman Ou <ochapman.cn@gmail.com>
* Created:	2014-08-20
*/
package godmi

import (
	"fmt"
)

type SystemPowerSupplyType byte

const (
	SystemPowerSupplyTypeOther SystemPowerSupplyType = 1 + iota
	SystemPowerSupplyTypeUnknown
	SystemPowerSupplyTypeLinear
	SystemPowerSupplyTypeSwitching
	SystemPowerSupplyTypeBattery
	SystemPowerSupplyTypeUPS
	SystemPowerSupplyTypeConverter
	SystemPowerSupplyTypeRegulator
	SystemPowerSupplyTypeReserved
)

func (s SystemPowerSupplyType) String() string {
	types := [...]string{
		"Other",
		"Unknown",
		"Linear",
		"Switching",
		"Battery",
		"UPS",
		"Converter",
		"Regulator",
		"Reserved",
	}
	if s <= 7 {
		return types[s]
	}
	return types[8]
}

type SystemPowerSupplyStatus byte

const (
	SystemPowerSupplyStatusOther SystemPowerSupplyStatus = 1 + iota
	SystemPowerSupplyStatusUnknown
	SystemPowerSupplyStatusOK
	SystemPowerSupplyStatusNonCritical
	SystemPowerSupplyStatusCritical
)

func (s SystemPowerSupplyStatus) String() string {
	status := [...]string{
		"Other",
		"Unknown",
		"OK",
		"Non-critical",
		"Critical",
	}
	return status[s-1]
}

type SystemPowerSupplyInputVoltageSwitching byte

const (
	SystemPowerSupplyInputVoltageSwitchingOther SystemPowerSupplyInputVoltageSwitching = 1 + iota
	SystemPowerSupplyInputVoltageSwitchingUnknown
	SystemPowerSupplyInputVoltageSwitchingManual
	SystemPowerSupplyInputVoltageSwitchingAutoSwitch
	SystemPowerSupplyInputVoltageSwitchingWiderange
	SystemPowerSupplyInputVoltageSwitchingNotApplicable
	SystemPowerSupplyInputVoltageSwitchingReserved
)

func (s SystemPowerSupplyInputVoltageSwitching) String() string {
	switches := [...]string{
		"Other",
		"Unknown",
		"Manual",
		"Auto-switch",
		"Wide range",
		"Not applicable",
		"Reserved",
	}
	if s < 6 {
		return switches[s-1]
	}
	return switches[6]
}

type SystemPowerSupplyCharacteristics struct {
	DMTFPowerSupplyType       SystemPowerSupplyType
	Status                    SystemPowerSupplyStatus
	DMTFInputVoltageSwitching SystemPowerSupplyInputVoltageSwitching
	IsUnpluggedFromWall       bool
	IsPresent                 bool
	IsHotRepleaceable         bool
}

func newSystemPowerSupplyCharacteristics(ch uint16) SystemPowerSupplyCharacteristics {
	var sp SystemPowerSupplyCharacteristics
	sp.DMTFPowerSupplyType = SystemPowerSupplyType((ch & 0x3c00) >> 10)
	sp.Status = SystemPowerSupplyStatus((ch & 0x380) >> 7)
	sp.DMTFInputVoltageSwitching = SystemPowerSupplyInputVoltageSwitching((ch & 0x70) >> 3)
	sp.IsUnpluggedFromWall = (ch&0x04 != 0)
	sp.IsPresent = (ch&0x02 != 0)
	sp.IsHotRepleaceable = (ch&0x01 != 0)
	return sp
}

func (s SystemPowerSupplyCharacteristics) String() string {
	return fmt.Sprintf("System Power Supply Characteristics:\n"+
		"\t\t\tDMTF Power Supply Type: %s\n"+
		"\t\t\tStatus: %s\n"+
		"\t\t\tDMTF Input Voltage Switching: %s\n"+
		"\t\t\tIs Unplugged From Wall: %t\n"+
		"\t\t\tIs Present: %t\n"+
		"\t\t\tIs Hot Repleaceable: %t\n",
		s.DMTFPowerSupplyType,
		s.Status,
		s.DMTFInputVoltageSwitching,
		s.IsUnpluggedFromWall,
		s.IsPresent,
		s.IsHotRepleaceable)
}

type SystemPowerSupply struct {
	infoCommon
	PowerUnitGroup             byte
	Location                   string
	DeviceName                 string
	Manufacturer               string
	SerialNumber               string
	AssetTagNumber             string
	ModelPartNumber            string
	RevisionLevel              string
	MaxPowerCapacity           uint16
	PowerSupplyCharacteristics SystemPowerSupplyCharacteristics
	InputVoltageProbeHandle    uint16
	CoolingDeviceHandle        uint16
	InputCurrentProbeHandle    uint16
}

func (s SystemPowerSupply) String() string {
	return fmt.Sprintf("System Power Supply\n"+
		"\tPower Unit Group: %d\n"+
		"\tLocation: %s\n"+
		"\tDevice Name: %s\n"+
		"\tManufacturer: %s\n"+
		"\tSerial Number: %s\n"+
		"\tAsset Tag Number: %s\n"+
		"\tModel Part Number: %s\n"+
		"\tRevision Level: %s\n"+
		"\tMax Power Capacity: %d\n"+
		"\tPower Supply Characteristics: %s\n"+
		"\tInput Voltage Probe Handle: %d\n"+
		"\tCooling Device Handle: %d\n"+
		"\tInput Current Probe Handle: %d",
		s.PowerUnitGroup,
		s.Location,
		s.DeviceName,
		s.Manufacturer,
		s.SerialNumber,
		s.AssetTagNumber,
		s.ModelPartNumber,
		s.RevisionLevel,
		s.MaxPowerCapacity,
		s.PowerSupplyCharacteristics,
		s.InputVoltageProbeHandle,
		s.CoolingDeviceHandle,
		s.InputCurrentProbeHandle)
}
