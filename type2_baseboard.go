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

type FeatureFlags byte

// Baseboard feature flags
const (
	FeatureFlagsHostingBoard = 1 << iota
	FeatureFlagsAtLeastOneDaughter
	FeatureFlagsRemovable
	FeatureFlagsRepleaceable
	FeatureFlagsHotSwappable
	//FeatureFlagsReserved = 000b
)

func (f FeatureFlags) String() string {
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

type BoardType byte

const (
	BoardTypeUnknown BoardType = 1 + iota
	BoardTypeOther
	BoardTypeServerBlade
	BoardTypeConnectivitySwitch
	BoardTypeSystemManagementModule
	BoardTypeProcessorModule
	BoardTypeIOModule
	BoardTypeMemModule
	BoardTypeDaughterBoard
	BoardTypeMotherboard
	BoardTypeProcessorMemmoryModule
	BoardTypeProcessorIOModule
	BoardTypeInterconnectBoard
)

func (b BoardType) String() string {
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
	if b > BoardTypeUnknown && b < BoardTypeInterconnectBoard {
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
	FeatureFlags                   FeatureFlags
	LocationInChassis              string
	ChassisHandle                  uint16
	BoardType                      BoardType
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
