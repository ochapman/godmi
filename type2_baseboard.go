/*
* File Name:	type2_baseboard.go
* Description:
* Author:	Chapman Ou <ochapman.cn@gmail.com>
* Created:	2014-08-18 22:58:31
 */

package godmi

import (
	"fmt"
)

type BaseboardFeatureFlags byte

// Baseboard feature flags
const (
	BaseboardFeatureFlagsHostingBoard BaseboardFeatureFlags = 1 << iota
	BaseboardFeatureFlagsAtLeastOneDaughter
	BaseboardFeatureFlagsRemovable
	BaseboardFeatureFlagsRepleaceable
	BaseboardFeatureFlagsHotSwappable
	//FeatureFlagsReserved = 000b
)

func (f BaseboardFeatureFlags) String() string {
	features := [...]string{
		"Board is a hosting board", /* 0 */
		"Board requires at least one daughter board",
		"Board is removable",
		"Board is replaceable",
		"Board is hot swappable", /* 4 */
	}
	var s string
	for i := uint32(0); i < 5; i++ {
		if f&(1<<i) != 0 {
			s += "\n\t\t" + features[i]
		}
	}
	return s
}

type BaseboardType byte

const (
	BaseboardTypeUnknown BaseboardType = 1 + iota
	BaseboardTypeOther
	BaseboardTypeServerBlade
	BaseboardTypeConnectivitySwitch
	BaseboardTypeSystemManagementModule
	BaseboardTypeProcessorModule
	BaseboardTypeIOModule
	BaseboardTypeMemModule
	BaseboardTypeDaughterBoard
	BaseboardTypeMotherboard
	BaseboardTypeProcessorMemmoryModule
	BaseboardTypeProcessorIOModule
	BaseboardTypeInterconnectBoard
)

func (b BaseboardType) String() string {
	types := [...]string{
		"Unknown", /* 0x01 */
		"Other",
		"Server Blade",
		"Connectivity Switch",
		"System Management Module",
		"Processor Module",
		"I/O Module",
		"Memory Module",
		"Daughter Board",
		"Motherboard",
		"Processor+Memory Module",
		"Processor+I/O Module",
		"Interconnect Board", /* 0x0D */
	}
	if b > BaseboardTypeUnknown && b < BaseboardTypeInterconnectBoard {
		return types[b-1]
	}
	return "Out Of Spec"
}

type BaseboardInformation struct {
	infoCommon
	Manufacturer                   string
	ProductName                    string
	Version                        string
	SerialNumber                   string
	AssetTag                       string
	FeatureFlags                   BaseboardFeatureFlags
	LocationInChassis              string
	ChassisHandle                  uint16
	BoardType                      BaseboardType
	NumberOfContainedObjectHandles byte
	ContainedObjectHandles         []byte
}

func (b BaseboardInformation) String() string {
	return fmt.Sprintf("Base Board Information\n"+
		"\tManufacturer: %s\n"+
		"\tProduct Name: %s\n"+
		"\tVersion: %s\n"+
		"\tSerial Number: %s\n"+
		"\tAsset Tag: %s\n"+
		"\tFeatures:%s\n"+
		"\tLocation In Chassis: %s\n"+
		"\tType: %s",
		b.Manufacturer,
		b.ProductName,
		b.Version,
		b.SerialNumber,
		b.AssetTag,
		b.FeatureFlags,
		b.LocationInChassis,
		b.BoardType)
}

func newBaseboardInformation(h dmiHeader) dmiTyper {
	data := h.data
	return &BaseboardInformation{
		Manufacturer:      h.FieldString(int(data[0x04])),
		ProductName:       h.FieldString(int(data[0x05])),
		Version:           h.FieldString(int(data[0x06])),
		SerialNumber:      h.FieldString(int(data[0x07])),
		AssetTag:          h.FieldString(int(data[0x08])),
		FeatureFlags:      BaseboardFeatureFlags(data[0x09]),
		LocationInChassis: h.FieldString(int(data[0x0A])),
		BoardType:         BaseboardType(data[0x0D]),
	}
}

func GetBaseboardInformation() *BaseboardInformation {
	if d, ok := gdmi[SMBIOSStructureTypeBaseBoard]; ok {
		return d.(*BaseboardInformation)
	}
	return nil
}

func init() {
	addTypeFunc(SMBIOSStructureTypeBaseBoard, newBaseboardInformation)
}
