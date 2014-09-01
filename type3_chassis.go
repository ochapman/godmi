/*
* File Name:	type3_chassis.go
* Description:
* Author:	Chapman Ou <ochapman.cn@gmail.com>
* Created:	2014-08-18 23:07:35
 */

package godmi

import (
	"fmt"
)

type ChassisType byte

const (
	ChssisTypeOther ChassisType = 1 + iota
	ChssisTypeUnknown
	ChssisTypeDesktop
	ChssisTypeLowProfileDesktop
	ChssisTypePizzaBox
	ChssisTypeMiniTower
	ChssisTypeTower
	ChssisTypePortable
	ChssisTypeLaptop
	ChssisTypeNotebook
	ChssisTypeHandHeld
	ChssisTypeDockingStation
	ChssisTypeAllinOne
	ChssisTypeSubNotebook
	ChssisTypeSpaceSaving
	ChssisTypeLunchBox
	ChssisTypeMainServerChassis
	ChssisTypeExpansionChassis
	ChssisTypeSubChassis
	ChssisTypeBusExpansionChassis
	ChssisTypePeripheralChassis
	ChssisTypeRAIDChassis
	ChssisTypeRackMountChassis
	ChssisTypeSealedcasePC
	ChssisTypeMultiSystem
	ChssisTypeCompactPCI
	ChssisTypeAdvancedTCA
	ChssisTypeBlade
	ChssisTypeBladeEnclosure
)

func (c ChassisType) String() string {
	types := [...]string{
		"Other",
		"Unknown",
		"Desktop",
		"LowProfileDesktop",
		"PizzaBox",
		"MiniTower",
		"Tower",
		"Portable",
		"Laptop",
		"Notebook",
		"HandHeld",
		"DockingStation",
		"AllinOne",
		"SubNotebook",
		"SpaceSaving",
		"LunchBox",
		"MainServerChassis",
		"ExpansionChassis",
		"SubChassis",
		"BusExpansionChassis",
		"PeripheralChassis",
		"RAIDChassis",
		"RackMountChassis",
		"SealedcasePC",
		"MultiSystem",
		"CompactPCI",
		"AdvancedTCA",
		"Blade",
		"BladeEnclosure",
	}
	c &= 0x7F
	if c >= 0x01 && c < 0x1D {
		return types[c-1]
	}
	return OUT_OF_SPEC
}

type ChassisLock byte

func (c ChassisLock) String() string {
	locks := [...]string{
		"Not Present", /* 0x00 */
		"Present",     /* 0x01 */
	}
	return locks[c]
}

type ChassisState byte

const (
	ChassisStateOther ChassisState = 1 + iota
	ChassisStateUnknown
	ChassisStateSafe
	ChassisStateWarning
	ChassisStateCritical
	ChassisStateNonRecoverable
)

func (c ChassisState) String() string {
	states := [...]string{
		"Other",
		"Unknown",
		"Safe",
		"Warning",
		"Critical",
		"NonRecoverable",
	}
	return states[c-1]
}

type ChassisContainedElementType byte

type ChassisContainedElements struct {
	Type    ChassisContainedElementType
	Minimum byte
	Maximum byte
}

type ChassisSecurityStatus byte

const (
	ChassisSecurityStatusOther ChassisSecurityStatus = 1 + iota
	ChassisSecurityStatusUnknown
	ChassisSecurityStatusNone
	ChassisSecurityStatusExternalInterfaceLockedOut
	ChassisSecurityStatusExternalInterfaceEnabled
)

func (s ChassisSecurityStatus) String() string {
	status := [...]string{
		"Other",
		"Unknown",
		"None",
		"ExternalInterfaceLockedOut",
		"ExternalInterfaceEnabled",
	}
	return status[s-1]
}

type ChassisHeight byte

type ChassisInformation struct {
	infoCommon
	Manufacturer                 string
	Type                         ChassisType
	Lock                         ChassisLock
	Version                      string
	AssetTag                     string
	SerialNumber                 string
	BootUpState                  ChassisState
	PowerSupplyState             ChassisState
	ThermalState                 ChassisState
	SecurityStatus               ChassisSecurityStatus
	OEMdefined                   uint16
	Height                       ChassisHeight
	NumberOfPowerCords           byte
	ContainedElementCount        byte
	ContainedElementRecordLength byte
	ContainedElements            ChassisContainedElements
	SKUNumber                    string
}

func (c ChassisInformation) String() string {
	return fmt.Sprintf("Chassis Information\n"+
		"\tManufacturer: %s\n"+
		"\tType: %s\n"+
		"\tLock: %s\n"+
		"\tVersion: %s\n"+
		"\tSerial Number: %s\n"+
		"\tAsset Tag: %s\n"+
		"\tBoot-up State: %s\n"+
		"\tPower Supply State: %s\n"+
		"\tThermal State: %s\n"+
		"\tSecurity Status: %s",
		c.Manufacturer,
		c.Type,
		c.Lock,
		c.Version,
		c.SerialNumber,
		c.AssetTag,
		c.BootUpState,
		c.PowerSupplyState,
		c.ThermalState,
		c.SecurityStatus)
}

func newChassisInformation(h dmiHeader) dmiTyper {
	data := h.data
	return &ChassisInformation{
		Manufacturer:                 h.FieldString(int(data[0x04])),
		Type:                         ChassisType(data[0x05]),
		Lock:                         ChassisLock(data[0x05] >> 7),
		Version:                      h.FieldString(int(data[0x06])),
		SerialNumber:                 h.FieldString(int(data[0x07])),
		AssetTag:                     h.FieldString(int(data[0x08])),
		BootUpState:                  ChassisState(data[0x09]),
		PowerSupplyState:             ChassisState(data[0xA]),
		ThermalState:                 ChassisState(data[0x0B]),
		SecurityStatus:               ChassisSecurityStatus(data[0x0C]),
		OEMdefined:                   u16(data[0x0D : 0x0D+4]),
		Height:                       ChassisHeight(data[0x11]),
		NumberOfPowerCords:           data[0x12],
		ContainedElementCount:        data[0x13],
		ContainedElementRecordLength: data[0x14],
		// TODO: 7.4.4
		//ci.ContainedElements:
		SKUNumber: h.FieldString(int(data[0x15])),
	}
}

func GetChassisInformation() *ChassisInformation {
	if d, ok := gdmi[SMBIOSStructureTypeChassis]; ok {
		return d.(*ChassisInformation)
	}
	return nil
}

func init() {
	addTypeFunc(SMBIOSStructureTypeChassis, newChassisInformation)
}
