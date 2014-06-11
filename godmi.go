/*
*
 */
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"syscall"
)

type Type byte
type Handle uint16
type InfoCommon struct {
	Type   Type
	Length byte
	Handle Handle
}

type SMBIOS_EPS struct {
	Anchor        []byte //4
	Checksum      byte
	Length        byte
	MajorVersion  byte
	MinorVersion  byte
	MaxSize       uint16
	Revision      byte
	FormattedArea []byte // 5
	InterAnchor   []byte // 5
	InterChecksum byte
	TableLength   uint16
	TableAddress  uint32
	NumberOfSM    uint16
	BCDRevision   byte
}

type DMIHeader struct {
	Type   byte
	Length byte
	Handle uint16
	data   []byte
}

type SMBIOS_Structure struct {
}

type Characteristics uint64
type CharacteristicsExt1 byte
type CharacteristicsExt2 byte

type BIOSInformation struct {
	Type                                   byte
	Length                                 byte
	Handle                                 uint16
	Vendor                                 string
	BIOSVersion                            string
	StartingAddressSegment                 uint16
	ReleaseDate                            string
	RomSize                                byte
	Characteristics                        Characteristics
	CharacteristicsExt1                    CharacteristicsExt1
	CharacteristicsExt2                    CharacteristicsExt2
	SystemBIOSMajorRelease                 byte
	SystemBIOSMinorRelease                 byte
	EmbeddedControllerFirmwareMajorRelease byte
	EmbeddedControllerFirmawreMinorRelease byte
}

type WakeUpType byte

const (
	Reserved WakeUpType = iota
	Other
	Unknown
	APM_Timer
	Modem_Ring
	LAN_Remote
	Power_Switch
	PCI_PME
	AC_Power_Restored
)

func (w WakeUpType) String() string {
	types := [...]string{
		"Reserved", /* 0x00 */
		"Other",
		"Unknown",
		"APM Timer",
		"Modem Ring",
		"LAN Remote",
		"Power Switch",
		"PCI PME#",
		"AC Power Restored", /* 0x08 */
	}
	return types[w]
}

type SystemInformation struct {
	Type         byte
	Length       byte
	Handle       uint16
	Manufacturer string
	ProductName  string
	Version      string
	SerialNumber string
	UUID         string
	WakeUpType   WakeUpType
	SKUNumber    string
	Family       string
}

func (si SystemInformation) String() string {
	return fmt.Sprintf("SystemInformation:\n\tManufacturer: %s\n\tProduct Name: %s\n\tVersion: %s\n\tSerial Number: %s\n\tUUID: %s\n\tWake-up Type: %s\n\tSKU Number: %s\n\tFamily: %s\n\t", si.Manufacturer, si.ProductName, si.Version, si.SerialNumber, si.UUID, si.WakeUpType, si.SKUNumber, si.Family)
}

type FeatureFlags byte

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
		//fmt.Printf("F%08b\nI%08b\n", f, 1<<i)
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
	Type                           byte
	Length                         byte
	Handle                         uint16
	Manufacturer                   string
	Product                        string
	Version                        string
	SerailNumber                   string
	AssetTag                       string
	FeatureFlags                   FeatureFlags
	LocationInChassis              string
	ChassisHandle                  uint16
	BoardType                      BoardType
	NumberOfContainedObjectHandles byte
	ContainedObjectHandles         []byte
}

func (bi BaseboardInformation) String() string {
	return fmt.Sprintf("BaseboardInformation:\n\tManufacturer: %s\n\tProduct: %s\n\tVersion: %s\n\tSerial Number: %s\n\tAsset Tag: %s\n\tFeature Flags: %s\n\tLocation In Chassis: %s\n\tBoard Type: %s\n\t", bi.Manufacturer, bi.Product, bi.Version, bi.SerailNumber, bi.AssetTag, bi.FeatureFlags, bi.LocationInChassis, bi.BoardType)
}

// BIOS Characteristics
const (
	BIOSCharacteristicsReserved0 = 1 << iota
	BIOSCharacteristicsReserved1
	BIOSCharacteristicsUnknown
	BIOSCharacteristicsNotSupported
	BIOSCharacteristicsISASupported
	BIOSCharacteristicsMCASupported
	BIOSCharacteristicsEISASupported
	BIOSCharacteristicsPCISupported
	BIOSCharacteristicsPCMCIASupported
	BIOSCharacteristicsPlugPlaySupported
	BIOSCharacteristicsAPMSupported
	BIOSCharacteristicsUpgradeable
	BIOSCharacteristicsShadowingIsAllowed
	BIOSCharacteristicsVLVESASupported
	BIOSCharacteristicsESCDSupported
	BIOSCharacteristicsBootFromCDSupported
	BIOSCharacteristicsSelectableBootSupported
	BIOSCharacteristicsBIOSROMIsSockectd
	BIOSCharacteristicsBootFromPCMCIASupported
	BIOSCharacteristicsEDDSupported
	BIOSCharacteristicsJPFloppyNECSupported
	BIOSCharacteristicsJPFloppyToshibaSupported
	BIOSCharacteristics525_360KBFloppySupported
	BIOSCharacteristics525_1_2MBFloppySupported
	BIOSCharacteristics35_720KBFloppySupported
	BIOSCharacteristics35_2_88MBFloppySupported
	BIOSCharacteristicsPrintScreenSupported
	BIOSCharacteristics8042KeyboardSupported
	BIOSCharacteristicsSerialSupported
	BIOSCharacteristicsPrinterSupported
	BIOSCharacteristicsCGAMonoSupported
	BIOSCharacteristicsNECPC98
	//Bit32:47 Reserved for BIOS vendor
	//Bit47:63 Reserved for system vendor
)

var sBIOSCharacteristics = [...]string{
	"BIOS characteristics not supported", /* 3 */
	"ISA is supported",
	"MCA is supported",
	"EISA is supported",
	"PCI is supported",
	"PC Card (PCMCIA) is supported",
	"PNP is supported",
	"APM is supported",
	"BIOS is upgradeable",
	"BIOS shadowing is allowed",
	"VLB is supported",
	"ESCD support is available",
	"Boot from CD is supported",
	"Selectable boot is supported",
	"BIOS ROM is socketed",
	"Boot from PC Card (PCMCIA) is supported",
	"EDD is supported",
	"Japanese floppy for NEC 9800 1.2 MB is supported (int 13h)",
	"Japanese floppy for Toshiba 1.2 MB is supported (int 13h)",
	"5.25\"/360 kB floppy services are supported (int 13h)",
	"5.25\"/1.2 MB floppy services are supported (int 13h)",
	"3.5\"/720 kB floppy services are supported (int 13h)",
	"3.5\"/2.88 MB floppy services are supported (int 13h)",
	"Print screen service is supported (int 5h)",
	"8042 keyboard services are supported (int 9h)",
	"Serial services are supported (int 14h)",
	"Printer services are supported (int 17h)",
	"CGA/mono video services are supported (int 10h)",
	"NEC PC-98", /* 31 */
}

var sCharateristicsExt1 = [...]string{
	"ACPI is supported", /* 0 */
	"USB legacy is supported",
	"AGP is supported",
	"I2O boot is supported",
	"LS-120 boot is supported",
	"ATAPI Zip drive boot is supported",
	"IEEE 1394 boot is supported",
	"Smart battery is supported", /* 7 */
}

var sCharateristicsExt2 = [...]string{
	"BIOS boot specification is supported", /* 0 */
	"Function key-initiated network boot is supported",
	"Targeted content distribution is supported",
	"UEFI is supported",
	"System is a virtual machine", /* 4 */
}

// BIOS Characteristics Extension Bytes
// Byte 1
const (
	BIOSCharacteristicsExt1ACPISupported = 1 << iota
	BIOSCharacteristicsExt1USBLegacySupported
	BIOSCharacteristicsExt1AGPSupported
	BIOSCharacteristicsExt1I2OBootSupported
	BIOSCharacteristicsExt1LS120SupperDiskBootSupported
	BIOSCharacteristicsExt1ATAPIZIPDriveBootSupported
	BIOSCharacteristicsExt11394BootSupported
	BIOSCharacteristicsExt1SmartBatterySupported
)

// Byte 2
const (
	BIOSCharacteristicsExt2BIOSBootSpecSupported = 1 << iota
	BIOSCharacteristicsExt2FuncKeyInitiatedNetworkBootSupported
	BIOSCharacteristicsExt2EnableTargetedContentDistribution
	BIOSCharacteristicsExt2UEFISpecSupported
	BIOSCharacteristicsExt2VirtualMachine
	// Bits 5:7 Reserved for future assignment
)

//BaseboardFeatureFlags
const (
	FeatureFlagsHostingBoard       = iota
	FeatureFlagsAtLeastOneDaughter = 1 << 1
	FeatureFlagsRemovable          = 1 << 2
	FeatureFlagsRepleaceable       = 1 << 3
	FeatureFlagsHotSwappable       = 1 << 4
	//FeatureFlagsReserved = 000b
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

func (ct ChassisType) String() string {
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
	return types[ct-1]
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

func (cc ChassisState) String() string {
	states := [...]string{
		"Other",
		"Unknown",
		"Safe",
		"Warning",
		"Critical",
		"NonRecoverable",
	}
	return states[cc-1]
}

type SecurityStatus byte

const (
	SecurityStatusOther SecurityStatus = 1 + iota
	SecurityStatusUnknown
	SecurityStatusNone
	SecurityStatusExternalInterfaceLockedOut
	SecurityStatusExternalInterfaceEnabled
)

func (ss SecurityStatus) String() string {
	status := [...]string{
		"Other",
		"Unknown",
		"None",
		"ExternalInterfaceLockedOut",
		"ExternalInterfaceEnabled",
	}
	return status[ss-1]
}

type ContainedElementType byte

type ContainedElements struct {
	Type    ContainedElementType
	Minimum byte
	Maximum byte
}

type Height byte

// type 3
type ChassisInformation struct {
	Type                         byte
	Length                       byte
	Handle                       byte
	Manufacturer                 string
	ChassisType                  ChassisType
	Version                      string
	AssetTag                     string
	SerialNumber                 string
	BootUpState                  ChassisState
	PowerSupplyState             ChassisState
	ThermalState                 ChassisState
	SecurityStatus               SecurityStatus
	OEMdefined                   uint16
	Height                       Height
	NumberOfPowerCords           byte
	ContainedElementCount        byte
	ContainedElementRecordLength byte
	ContainedElements            ContainedElements
	SKUNumber                    string
}

func (h DMIHeader) ChassisInformation() ChassisInformation {
	var ci ChassisInformation
	data := h.data
	ci.Manufacturer = h.FieldString(int(data[0x04]))
	ci.ChassisType = ChassisType(data[0x05])
	ci.Version = h.FieldString(int(data[0x06]))
	ci.SerialNumber = h.FieldString(int(data[0x07]))
	ci.AssetTag = h.FieldString(int(data[0x08]))
	ci.BootUpState = ChassisState(data[0x09])
	ci.PowerSupplyState = ChassisState(data[0xA])
	ci.ThermalState = ChassisState(data[0x0B])
	ci.SecurityStatus = SecurityStatus(data[0x0C])
	ci.OEMdefined = U16(data[0x0D : 0x0D+4])
	ci.Height = Height(data[0x11])
	ci.NumberOfPowerCords = data[0x12]
	ci.ContainedElementCount = data[0x13]
	ci.ContainedElementRecordLength = data[0x14]
	// TODO: 7.4.4
	//ci.ContainedElements =
	ci.SKUNumber = h.FieldString(int(data[0x15]))
	return ci
}

func (ci ChassisInformation) String() string {
	return fmt.Sprintf("Chassis Information:\n\tManufacturer: %s\n\tType: %s\n\tVersion: %s\n\tSerial Number: %s\n\tAsset Tag: %s\n\tBoot-up State: %s\n\tPower Supply State: %s\n\tThermal State: %s\n\tSecurity Status: %s\n\t", ci.Manufacturer, ci.ChassisType, ci.Version, ci.SerialNumber, ci.AssetTag, ci.BootUpState, ci.PowerSupplyState, ci.ThermalState, ci.SecurityStatus)
}

type ProcessorType byte

const (
	ProcessorTypeOther ProcessorType = 1 + iota
	ProcessorTypeUnknown
	ProcessorTypeCentralProcessor
	ProcessorTypeMathProcessor
	ProcessorTypeDSPProcessor
	ProcessorTypeVideoProcessor
)

func (pt ProcessorType) String() string {
	types := [...]string{
		"Other",
		"Unknown",
		"CentralProcessor",
		"MathProcessor",
		"DSPProcessor",
		"VideoProcessor",
	}
	return types[pt-1]
}

type ProcessorFamily uint16

const (
	_ ProcessorFamily = iota
	ProcessorOther
	ProcessorUnknown
	ProcessorProcessorFamily8086
	ProcessorProcessorFamily80286
	ProcessorIntel386TMprocessor
	ProcessorIntel486TMprocessor
	ProcessorProcessorFamily8087
	ProcessorProcessorFamily80287
	ProcessorProcessorFamily80387
	ProcessorProcessorFamily80487
	ProcessorIntelPentiumprocessor
	ProcessorPentiumProprocessor
	ProcessorPentiumIIprocessor
	ProcessorPentiumprocessorwithMMXTMtechnology
	ProcessorIntelCeleronprocessor
	ProcessorPentiumIIXeonTMprocessor
	ProcessorPentiumIIIprocessor
	ProcessorM1Family
	ProcessorM2Family
	ProcessorIntelCeleronMprocessor
	ProcessorIntelPentium4HTprocessor
	_
	_
	ProcessorAMDDuronTMProcessorFamily
	ProcessorK5Family
	ProcessorK6Family
	ProcessorK6_2
	ProcessorK6_3
	ProcessorAMDAthlonTMProcessorFamily
	ProcessorAMD29000Family
	ProcessorK6_2Plus
	ProcessorPowerPCFamily
	ProcessorPowerPC601
	ProcessorPowerPC603
	ProcessorPowerPC603Plus
	ProcessorPowerPC604
	ProcessorPowerPC620
	ProcessorPowerPCx704
	ProcessorPowerPC750
	ProcessorIntelCoreTMDuoprocessor
	ProcessorIntelCoreTMDuomobileprocessor
	ProcessorIntelCoreTMSolomobileprocessor
	ProcessorIntelAtomTMprocessor
	_
	_
	_
	_
	ProcessorAlphaFamily
	ProcessorAlpha21064
	ProcessorAlpha21066
	ProcessorAlpha21164
	ProcessorAlpha21164PC
	ProcessorAlpha21164a
	ProcessorAlpha21264
	ProcessorAlpha21364
	ProcessorAMDTurionTMIIUltraDual_CoreMobileMProcessorFamily
	ProcessorAMDTurionTMIIDual_CoreMobileMProcessorFamily
	ProcessorAMDAthlonTMIIDual_CoreMProcessorFamily
	ProcessorAMDOpteronTM6100SeriesProcessor
	ProcessorAMDOpteronTM4100SeriesProcessor
	ProcessorAMDOpteronTM6200SeriesProcessor
	ProcessorAMDOpteronTM4200SeriesProcessor
	ProcessorAMDFXTMSeriesProcessor
	ProcessorMIPSFamily
	ProcessorMIPSR4000
	ProcessorMIPSR4200
	ProcessorMIPSR4400
	ProcessorMIPSR4600
	ProcessorMIPSR10000
	ProcessorAMDC_SeriesProcessor
	ProcessorAMDE_SeriesProcessor
	ProcessorAMDA_SeriesProcessor
	ProcessorAMDG_SeriesProcessor
	ProcessorAMDZ_SeriesProcessor
	ProcessorAMDR_SeriesProcessor
	ProcessorAMDOpteronTM4300SeriesProcessor
	ProcessorAMDOpteronTM6300SeriesProcessor
	ProcessorAMDOpteronTM3300SeriesProcessor
	ProcessorAMDFireProTMSeriesProcessor
	ProcessorSPARCFamily
	ProcessorSuperSPARC
	ProcessormicroSPARCII
	ProcessormicroSPARCIIep
	ProcessorUltraSPARC
	ProcessorUltraSPARCII
	ProcessorUltraSPARCIii
	ProcessorUltraSPARCIII
	ProcessorUltraSPARCIIIi
	_
	_
	_
	_
	_
	_
	_
	Processor68040Family
	Processor68xxx
	ProcessorProcessorFamily68000
	ProcessorProcessorFamily68010
	ProcessorProcessorFamily68020
	ProcessorProcessorFamily68030
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	ProcessorHobbitFamily
	_
	_
	_
	_
	_
	_
	_
	ProcessorCrusoeTMTM5000Family
	ProcessorCrusoeTMTM3000Family
	ProcessorEfficeonTMTM8000Family
	_
	_
	_
	_
	_
	ProcessorWeitek
	_
	ProcessorItaniumTMprocessor
	ProcessorAMDAthlonTM64ProcessorFamily
	ProcessorAMDOpteronTMProcessorFamily
	ProcessorAMDSempronTMProcessorFamily
	ProcessorAMDTurionTM64MobileTechnology
	ProcessorDual_CoreAMDOpteronTMProcessorFamily
	ProcessorAMDAthlonTM64X2Dual_CoreProcessorFamily
	ProcessorAMDTurionTM64X2MobileTechnology
	ProcessorQuad_CoreAMDOpteronTMProcessorFamily
	ProcessorThird_GenerationAMDOpteronTMProcessorFamily
	ProcessorAMDPhenomTMFXQuad_CoreProcessorFamily
	ProcessorAMDPhenomTMX4Quad_CoreProcessorFamily
	ProcessorAMDPhenomTMX2Dual_CoreProcessorFamily
	ProcessorAMDAthlonTMX2Dual_CoreProcessorFamily
	ProcessorPA_RISCFamily
	ProcessorPA_RISC8500
	ProcessorPA_RISC8000
	ProcessorPA_RISC7300LC
	ProcessorPA_RISC7200
	ProcessorPA_RISC7100LC
	ProcessorPA_RISC7100
	_
	_
	_
	_
	_
	_
	_
	_
	_
	ProcessorV30Family
	ProcessorQuad_CoreIntelXeonprocessor3200Series
	ProcessorDual_CoreIntelXeonprocessor3000Series
	ProcessorQuad_CoreIntelXeonprocessor5300Series
	ProcessorDual_CoreIntelXeonprocessor5100Series
	ProcessorDual_CoreIntelXeonprocessor5000Series
	ProcessorDual_CoreIntelXeonprocessorLV
	ProcessorDual_CoreIntelXeonprocessorULV
	ProcessorDual_CoreIntelXeonprocessor7100Series
	ProcessorQuad_CoreIntelXeonprocessor5400Series
	ProcessorQuad_CoreIntelXeonprocessor
	ProcessorDual_CoreIntelXeonprocessor5200Series
	ProcessorDual_CoreIntelXeonprocessor7200Series
	ProcessorQuad_CoreIntelXeonprocessor7300Series
	ProcessorQuad_CoreIntelXeonprocessor7400Series
	ProcessorMulti_CoreIntelXeonprocessor7400Series
	ProcessorPentiumIIIXeonTMprocessor
	ProcessorPentiumIIIProcessorwithIntelSpeedStepTMTechnology
	ProcessorPentium4Processor
	ProcessorIntelXeonprocessor
	ProcessorAS400Family
	ProcessorIntelXeonTMprocessorMP
	ProcessorAMDAthlonTMXPProcessorFamily
	ProcessorAMDAthlonTMMPProcessorFamily
	ProcessorIntelItanium2processor
	ProcessorIntelPentiumMprocessor
	ProcessorIntelCeleronDprocessor
	ProcessorIntelPentiumDprocessor
	ProcessorIntelPentiumProcessorExtremeEdition
	ProcessorIntelCoreTMSoloProcessor
	ProcessorReserved
	ProcessorIntelCoreTM2DuoProcessor
	ProcessorIntelCoreTM2Soloprocessor
	ProcessorIntelCoreTM2Extremeprocessor
	ProcessorIntelCoreTM2Quadprocessor
	ProcessorIntelCoreTM2Extrememobileprocessor
	ProcessorIntelCoreTM2Duomobileprocessor
	ProcessorIntelCoreTM2Solomobileprocessor
	ProcessorIntelCoreTMi7processor
	ProcessorDual_CoreIntelCeleronprocessor
	ProcessorIBM390Family
	ProcessorG4
	ProcessorG5
	ProcessorESA390G6
	ProcessorzArchitecturebase
	ProcessorIntelCoreTMi5processor
	ProcessorIntelCoreTMi3processor
	_
	_
	_
	ProcessorVIAC7TM_MProcessorFamily
	ProcessorVIAC7TM_DProcessorFamily
	ProcessorVIAC7TMProcessorFamily
	ProcessorVIAEdenTMProcessorFamily
	ProcessorMulti_CoreIntelXeonprocessor
	ProcessorDual_CoreIntelXeonprocessor3xxxSeries
	ProcessorQuad_CoreIntelXeonprocessor3xxxSeries
	ProcessorVIANanoTMProcessorFamily
	ProcessorDual_CoreIntelXeonprocessor5xxxSeries
	ProcessorQuad_CoreIntelXeonprocessor5xxxSeries
	_
	ProcessorDual_CoreIntelXeonprocessor7xxxSeries
	ProcessorQuad_CoreIntelXeonprocessor7xxxSeries
	ProcessorMulti_CoreIntelXeonprocessor7xxxSeries
	ProcessorMulti_CoreIntelXeonprocessor3400Series
	_
	_
	_
	ProcessorAMDOpteronTM3000SeriesProcessor
	ProcessorAMDSempronTMIIProcessor
	ProcessorEmbeddedAMDOpteronTMQuad_CoreProcessorFamily
	ProcessorAMDPhenomTMTriple_CoreProcessorFamily
	ProcessorAMDTurionTMUltraDual_CoreMobileProcessorFamily
	ProcessorAMDTurionTMDual_CoreMobileProcessorFamily
	ProcessorAMDAthlonTMDual_CoreProcessorFamily
	ProcessorAMDSempronTMSIProcessorFamily
	ProcessorAMDPhenomTMIIProcessorFamily
	ProcessorAMDAthlonTMIIProcessorFamily
	ProcessorSix_CoreAMDOpteronTMProcessorFamily
	ProcessorAMDSempronTMMProcessorFamily
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_
	Processori860
	Processori960
	_
	_
	ProcessorIndicatortoobtaintheprocessorfamilyfromtheProcessorFamily2field
	_
	_
	_
	_
	_
	_
	ProcessorSH_3
	ProcessorSH_4
	ProcessorARM
	ProcessorStrongARM
	Processor6x86
	ProcessorMediaGX
	ProcessorMII
	ProcessorWinChip
	ProcessorDSP
	ProcessorVideoProcessor
	_
	_
)

func (pf ProcessorFamily) String() string {
	families := [...]string{
		"Other",
		"Unknown",
		"8086",
		"80286",
		"Intel386TM processor",
		"Intel486TM processor",
		"8087",
		"80287",
		"80387",
		"80487",
		"Intel® Pentium® processor",
		"Pentium® Pro processor",
		"Pentium® II processor",
		"Pentium® processor with MMXTM technology",
		"Intel® Celeron® processor",
		"Pentium® II XeonTM processor",
		"Pentium® III processor",
		"M1 Family",
		"M2 Family",
		"Intel® Celeron® M processor",
		"Intel® Pentium® 4 HT processor",
		"Available for assignment",
		"Available for assignment",
		"AMD DuronTM Processor Family",
		"K5 Family",
		"K6 Family",
		"K6-2",
		"K6-3",
		"AMD AthlonTM Processor Family",
		"AMD29000 Family",
		"K6-2+",
		"Power PC Family",
		"Power PC 601",
		"Power PC 603",
		"Power PC 603+",
		"Power PC 604",
		"Power PC 620",
		"Power PC x704",
		"Power PC 750",
		"Intel® CoreTM Duo processor",
		"Intel® CoreTM Duo mobile processor",
		"Intel® CoreTM Solo mobile processor",
		"Intel® AtomTM processor",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Alpha Family",
		"Alpha 21064",
		"Alpha 21066",
		"Alpha 21164",
		"Alpha 21164PC",
		"Alpha 21164a",
		"Alpha 21264",
		"Alpha 21364",
		"AMD TurionTM II Ultra Dual-Core Mobile M Processor Family",
		"AMD TurionTM II Dual-Core Mobile M Processor Family",
		"AMD AthlonTM II Dual-Core M Processor Family",
		"AMD OpteronTM 6100 Series Processor",
		"AMD OpteronTM 4100 Series Processor",
		"AMD OpteronTM 6200 Series Processor",
		"AMD OpteronTM 4200 Series Processor",
		"AMD FXTM Series Processor",
		"MIPS Family",
		"MIPS R4000",
		"MIPS R4200",
		"MIPS R4400",
		"MIPS R4600",
		"MIPS R10000",
		"AMD C-Series Processor",
		"AMD E-Series Processor",
		"AMD A-Series Processor",
		"AMD G-Series Processor",
		"AMD Z-Series Processor",
		"AMD R-Series Processor",
		"AMD OpteronTM 4300 Series Processor",
		"AMD OpteronTM 6300 Series Processor",
		"AMD OpteronTM 3300 Series Processor",
		"AMD FireProTM Series Processor",
		"SPARC Family",
		"SuperSPARC",
		"microSPARC II",
		"microSPARC IIep",
		"UltraSPARC",
		"UltraSPARC II",
		"UltraSPARC Iii",
		"UltraSPARC III",
		"UltraSPARC IIIi",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"68040 Family",
		"68xxx",
		"68000",
		"68010",
		"68020",
		"68030",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Hobbit Family",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"CrusoeTM TM5000 Family",
		"CrusoeTM TM3000 Family",
		"EfficeonTM TM8000 Family",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Weitek",
		"Available for assignment",
		"ItaniumTM processor",
		"AMD AthlonTM 64 Processor Family",
		"AMD OpteronTM Processor Family",
		"AMD SempronTM Processor Family",
		"AMD TurionTM 64 Mobile Technology",
		"Dual-Core AMD OpteronTM Processor Family",
		"AMD AthlonTM 64 X2 Dual-Core Processor Family",
		"AMD TurionTM 64 X2 Mobile Technology",
		"Quad-Core AMD OpteronTM Processor Family",
		"Third-Generation AMD OpteronTM Processor Family",
		"AMD PhenomTM FX Quad-Core Processor Family",
		"AMD PhenomTM X4 Quad-Core Processor Family",
		"AMD PhenomTM X2 Dual-Core Processor Family",
		"AMD AthlonTM X2 Dual-Core Processor Family",
		"PA-RISC Family",
		"PA-RISC 8500",
		"PA-RISC 8000",
		"PA-RISC 7300LC",
		"PA-RISC 7200",
		"PA-RISC 7100LC",
		"PA-RISC 7100",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"V30 Family",
		"Quad-Core Intel® Xeon® processor 3200 Series",
		"Dual-Core Intel® Xeon® processor 3000 Series",
		"Quad-Core Intel® Xeon® processor 5300 Series",
		"Dual-Core Intel® Xeon® processor 5100 Series",
		"Dual-Core Intel® Xeon® processor 5000 Series",
		"Dual-Core Intel® Xeon® processor LV",
		"Dual-Core Intel® Xeon® processor ULV",
		"Dual-Core Intel® Xeon® processor 7100 Series",
		"Quad-Core Intel® Xeon® processor 5400 Series",
		"Quad-Core Intel® Xeon® processor",
		"Dual-Core Intel® Xeon® processor 5200 Series",
		"Dual-Core Intel® Xeon® processor 7200 Series",
		"Quad-Core Intel® Xeon® processor 7300 Series",
		"Quad-Core Intel® Xeon® processor 7400 Series",
		"Multi-Core Intel® Xeon® processor 7400 Series",
		"Pentium® III XeonTM processor",
		"Pentium® III Processor with Intel® SpeedStepTM Technology",
		"Pentium® 4 Processor",
		"Intel® Xeon® processor",
		"AS400 Family",
		"Intel® XeonTM processor MP",
		"AMD AthlonTM XP Processor Family",
		"AMD AthlonTM MP Processor Family",
		"Intel® Itanium® 2 processor",
		"Intel® Pentium® M processor",
		"Intel® Celeron® D processor",
		"Intel® Pentium® D processor",
		"Intel® Pentium® Processor Extreme Edition",
		"Intel® CoreTM Solo Processor",
		"Reserved",
		"Intel® CoreTM 2 Duo Processor",
		"Intel® CoreTM 2 Solo processor",
		"Intel® CoreTM 2 Extreme processor",
		"Intel® CoreTM 2 Quad processor",
		"Intel® CoreTM 2 Extreme mobile processor",
		"Intel® CoreTM 2 Duo mobile processor",
		"Intel® CoreTM 2 Solo mobile processor",
		"Intel® CoreTM i7 processor",
		"Dual-Core Intel® Celeron® processor",
		"IBM390 Family",
		"G4",
		"G5",
		"ESA/390 G6",
		"z/Architecture base",
		"Intel® CoreTM i5 processor",
		"Intel® CoreTM i3 processor",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"VIA C7TM-M Processor Family",
		"VIA C7TM-D Processor Family",
		"VIA C7TM Processor Family",
		"VIA EdenTM Processor Family",
		"Multi-Core Intel® Xeon® processor",
		"Dual-Core Intel® Xeon® processor 3xxx Series",
		"Quad-Core Intel® Xeon® processor 3xxx Series",
		"VIA NanoTM Processor Family",
		"Dual-Core Intel® Xeon® processor 5xxx Series",
		"Quad-Core Intel® Xeon® processor 5xxx Series",
		"Available for assignment",
		"Dual-Core Intel® Xeon® processor 7xxx Series",
		"Quad-Core Intel® Xeon® processor 7xxx Series",
		"Multi-Core Intel® Xeon® processor 7xxx Series",
		"Multi-Core Intel® Xeon® processor 3400 Series",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"AMD OpteronTM 3000 Series Processor",
		"AMD SempronTM II Processor",
		"Embedded AMD OpteronTM Quad-Core Processor Family",
		"AMD PhenomTM Triple-Core Processor Family",
		"AMD TurionTM Ultra Dual-Core Mobile Processor Family",
		"AMD TurionTM Dual-Core Mobile Processor Family",
		"AMD AthlonTM Dual-Core Processor Family",
		"AMD SempronTM SI Processor Family",
		"AMD PhenomTM II Processor Family",
		"AMD AthlonTM II Processor Family",
		"Six-Core AMD OpteronTM Processor Family",
		"AMD SempronTM M Processor Family",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"i860",
		"i960",
		"Available for assignment",
		"Available for assignment",
		"Indicator to obtain the processor family from the Processor Family 2 field",
		"Reserved ￼",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"Available for assignment",
		"SH-3",
		"SH-4",
		"ARM",
		"StrongARM",
		"6x86",
		"MediaGX",
		"MII",
		"WinChip",
		"DSP",
		"Video Processor",
		"Available for assignment",
		"Reserved",
	}
	return families[pf]
}

type ProcessorID uint64

type ProcessorVoltage byte

const (
	ProcessorVoltage5V ProcessorVoltage = 1 << iota
	ProcessorVoltage3dot3V
	ProcessorVoltage2dot9V
	ProcessorVoltageReserved
	_
	_
	_
	ProcessorVoltageLegacy
)

func (pv ProcessorVoltage) String() string {
	voltages := [...]string{
		"5V",
		"3.3V",
		"2.9V",
	}
	if pv&ProcessorVoltageLegacy == 0 {
		return voltages[pv]
	}
	return fmt.Sprintf("%.1f", (pv-0x80)/10)
}

type ProcessorStatus byte

const (
	ProcessorStatusUnknown ProcessorStatus = 1 << iota
	ProcessorStatusEnabled
	ProcessorStatusDisabledByUser
	ProcessorStatusDisabledByBIOS
	ProcessorStatusIdle
	ProcessorStatusReserved
	ProcessorStatusOther
)

func (ps ProcessorStatus) String() string {
	status := [...]string{
		"Unknown",
		"CPU Enabled",
		"Disabled By User through BIOS Setup",
		"Disabled By BIOSa(POST Error)",
		"CPU is Idle, waiting to be enabled",
		"Reserved",
		"Other",
	}
	return status[ps]
}

type ProcessorUpgrade byte

const (
	_ ProcessorUpgrade = iota
	ProcessorUpgradeOther
	ProcessorUpgradeUnknown
	ProcessorUpgradeDaughterBoard
	ProcessorUpgradeZIFSocket
	ProcessorUpgradeReplaceablePiggyBack
	ProcessorUpgradeNone
	ProcessorUpgradeLIFSocket
	ProcessorUpgradeSlot1
	ProcessorUpgradeSlot2
	ProcessorUpgrade370_pinsocket
	ProcessorUpgradeSlotA
	ProcessorUpgradeSlotM
	ProcessorUpgradeSocket423
	ProcessorUpgradeSocketASocket462
	ProcessorUpgradeSocket478
	ProcessorUpgradeSocket754
	ProcessorUpgradeSocket940
	ProcessorUpgradeSocket939
	ProcessorUpgradeSocketmPGA604
	ProcessorUpgradeSocketLGA771
	ProcessorUpgradeSocketLGA775
	ProcessorUpgradeSocketS1
	ProcessorUpgradeSocketAM2
	ProcessorUpgradeSocketF1207
	ProcessorUpgradeSocketLGA1366
	ProcessorUpgradeSocketG34
	ProcessorUpgradeSocketAM3
	ProcessorUpgradeSocketC32
	ProcessorUpgradeSocketLGA1156
	ProcessorUpgradeSocketLGA1567
	ProcessorUpgradeSocketPGA988A
	ProcessorUpgradeSocketBGA1288
	ProcessorUpgradeSocketrPGA988B
	ProcessorUpgradeSocketBGA1023
	ProcessorUpgradeSocketBGA1224
	ProcessorUpgradeSocketLGA1155
	ProcessorUpgradeSocketLGA1356
	ProcessorUpgradeSocketLGA2011
	ProcessorUpgradeSocketFS1
	ProcessorUpgradeSocketFS2
	ProcessorUpgradeSocketFM1
	ProcessorUpgradeSocketFM2
	ProcessorUpgradeSocketLGA2011_3
	ProcessorUpgradeSocketLGA1356_3
)

func (pu ProcessorUpgrade) String() string {
	upgrades := [...]string{
		"Other",
		"Unknown",
		"Daughter Board",
		"ZIF Socket",
		"Replaceable Piggy Back",
		"None",
		"LIF Socket",
		"Slot 1",
		"Slot 2",
		"370-pin socket",
		"Slot A",
		"Slot M",
		"Socket 423",
		"Socket A (Socket 462)",
		"Socket 478",
		"Socket 754",
		"Socket 940",
		"Socket 939",
		"Socket mPGA604",
		"Socket LGA771",
		"Socket LGA775",
		"Socket S1",
		"Socket AM2",
		"Socket F (1207)",
		"Socket LGA1366",
		"Socket G34",
		"Socket AM3",
		"Socket C32",
		"Socket LGA1156",
		"Socket LGA1567",
		"Socket PGA988A",
		"Socket BGA1288",
		"Socket rPGA988B",
		"Socket BGA1023",
		"Socket BGA1224",
		"Socket LGA1155",
		"Socket LGA1356",
		"Socket LGA2011",
		"Socket FS1",
		"Socket FS2",
		"Socket FM1",
		"Socket FM2",
		"Socket LGA2011-3",
		"Socket LGA1356-3",
	}
	return upgrades[pu]
}

type ProcessorCharacteristics uint16

const (
	ProcessorCharacteristicsReserved ProcessorCharacteristics = 1 << iota
	ProcessorCharacteristicsUnknown
	ProcessorCharacteristics64_bitCapable
	ProcessorCharacteristicsMulti_Core
	ProcessorCharacteristicsHardwareThread
	ProcessorCharacteristicsExecuteProtection
	ProcessorCharacteristicsEnhancedVirtualization
	ProcessorCharacteristicsPowerPerformanceControl
)

func (pc ProcessorCharacteristics) String() string {
	chars := [...]string{
		"Reserved",
		"Unknown",
		"64-bit Capable",
		"Multi-Core",
		"Hardware Thread",
		"Execute Protection",
		"Enhanced Virtualization",
		"Power/Performance Control",
	}
	return chars[pc]
}

// type 4
type ProcessorInformation struct {
	Type              byte
	Length            byte
	Handle            uint16
	SocketDesignation string
	ProcessorType     ProcessorType
	Family            ProcessorFamily
	Manufacturer      string
	ID                ProcessorID
	Version           string
	Voltage           ProcessorVoltage
	ExternalClock     uint16
	MaxSpeed          uint16
	CurrentSpeed      uint16
	Status            ProcessorStatus
	Upgrade           ProcessorUpgrade
	L1CacheHandle     uint16
	L2CacheHandle     uint16
	L3CacheHandle     uint16
	SerialNumber      string
	AssetTag          string
	PartNumber        string
	CoreCount         byte
	CoreEnabled       byte
	ThreadCount       byte
	Characteristics   ProcessorCharacteristics
	Family2           ProcessorFamily
}

func (h DMIHeader) ProcessorInformation() ProcessorInformation {
	var pi ProcessorInformation
	data := h.data
	pi.SocketDesignation = h.FieldString(int(data[0x04]))
	pi.ProcessorType = ProcessorType(data[0x05])
	pi.Family = ProcessorFamily(data[0x06])
	pi.Manufacturer = h.FieldString(int(data[0x07]))
	// TODO:
	//pi.ProcessorID
	pi.Version = h.FieldString(int(data[0x10]))
	pi.Voltage = ProcessorVoltage(data[0x11])
	pi.ExternalClock = U16(data[0x12:0x14])
	pi.MaxSpeed = U16(data[0x14:0x16])
	pi.CurrentSpeed = U16(data[0x16:0x18])
	pi.Status = ProcessorStatus(data[0x18])
	pi.Upgrade = ProcessorUpgrade(data[0x19])
	pi.L1CacheHandle = U16(data[0x1A:0x1C])
	pi.L2CacheHandle = U16(data[0x1C:0x1E])
	pi.L3CacheHandle = U16(data[0x1E:0x20])
	pi.SerialNumber = h.FieldString(int(data[0x20]))
	pi.AssetTag = h.FieldString(int(data[0x21]))
	pi.PartNumber = h.FieldString(int(data[0x22]))
	pi.CoreCount = data[0x23]
	pi.CoreEnabled = data[0x24]
	pi.ThreadCount = data[0x25]
	pi.Characteristics = ProcessorCharacteristics(U16(data[0x26:0x28]))
	pi.Family2 = ProcessorFamily(data[0x28])
	return pi
}

type OperationalMode byte

const (
	OperationalModeWriteThrough OperationalMode = iota
	OperationalModeWriteBack
	OperationalModeVariesWithMemoryAddress
	OperationalModeUnknown
)

func (o OperationalMode) String() string {
	modes := [...]string{
		"Write Through",
		"Write Back",
		"Varies With Memory Address",
		"Unknown",
	}
	return modes[o]
}

type CacheLocation byte

const (
	CacheLocationInternal CacheLocation = iota
	CacheLocationExternal
	CacheLocationReserved
	CacheLocationUnknown
)

func (c CacheLocation) String() string {
	locations := [...]string{
		"Internal",
		"External",
		"Reserved",
		"Unknown",
	}
	return locations[c]
}

type CacheLevel byte

const (
	Level1 CacheLevel = iota
	Level2
	Level3
)

func (c CacheLevel) String() string {
	levels := [...]string{
		"Level1",
		"Level2",
		"Level3",
	}
	return levels[c]
}

type CacheConfiguration struct {
	Mode     OperationalMode
	Enabled  bool
	Location CacheLocation
	Socketed bool
	Level    CacheLevel
}

func NewCacheConfiguration(u uint16) CacheConfiguration {
	var c CacheConfiguration
	c.Level = CacheLevel(byte(u & 0x7))
	c.Socketed = (u&0x10 == 1)
	c.Location = CacheLocation((u >> 5) & 0x3)
	c.Enabled = (u&(0x1<<7) == 1)
	c.Mode = OperationalMode((u >> 8) & 0x7)
	return c
}

func (c CacheConfiguration) String() string {
	return fmt.Sprintf("Cache Configuration: \n\tLevel: %s\n\t\tSocketed: %v\n\t\tLocation: %s\n\t\tEnabled: %v\n\t\tMode:\n\t\t", c.Level, c.Socketed, c.Location, c.Enabled, c.Mode)
}

type CacheGranularity byte

const (
	CacheGranularity1K CacheGranularity = iota
	CacheGranularity64K
)

func (c CacheGranularity) String() string {
	grans := [...]string{
		"1K",
		"64K",
	}
	return grans[c]
}

type CacheSize struct {
	Granularity CacheGranularity
	Size        uint16
}

func NewCacheSize(u uint16) CacheSize {
	var c CacheSize
	c.Granularity = CacheGranularity(u >> 15)
	c.Size = u &^ (uint16(1) << 15)
	return c
}

func (c CacheSize) String() string {
	return fmt.Sprintf("%s * %s", c.Size, c.Granularity)
}

type SRAMType uint16

const (
	SRAMTypeOther SRAMType = 1 << iota
	SRAMTypeUnknown
	SRAMTypeNonBurst
	SRAMTypeBurst
	SRAMTypePipelineBurst
	SRAMTypeSynchronous
	SRAMTypeAsynchronous
	SRAMTypeReserved
)

func (st SRAMType) String() string {
	types := [...]string{
		"Other",
		"Unknown",
		"Non-Burst",
		"Burst",
		"Pipeline Burst",
		"Synchronous",
		"Asynchronous",
		"Reserved",
	}
	return types[st/2]
}

type CacheSpeed byte

type ErrorCorrectionType byte

const (
	ErrorCorrectionTypeOther ErrorCorrectionType = 1 + iota
	ErrorCorrectionTypeUnknown
	ErrorCorrectionTypeNone
	ErrorCorrectionTypeParity
	ErrorCorrectionTypeSinglebitECC
	ErrorCorrectionTypeMultibitECC
)

func (ec ErrorCorrectionType) String() string {
	types := [...]string{
		"Other",
		"Unknown",
		"None",
		"Parity",
		"Single-bit ECC",
		"Multi-bit ECC",
	}
	return types[ec-1]
}

type SystemCacheType byte

const (
	SystemCacheTypeOther SystemCacheType = 1 + iota
	SystemCacheTypeUnknown
	SystemCacheTypeInstruction
	SystemCacheTypeData
	SystemCacheTypeUnified
)

func (s SystemCacheType) String() string {
	types := [...]string{
		"Other",
		"Unknown",
		"Instruction",
		"Data",
		"Unified",
	}
	return types[s-1]
}

type CacheAssociativity byte

const (
	CacheAssociativityOther CacheAssociativity = 1 + iota
	CacheAssociativityUnknown
	CacheAssociativityDirectMapped
	CacheAssociativity2waySetAssociative
	CacheAssociativity4waySetAssociative
	CacheAssociativityFullyAssociative
	CacheAssociativity8waySetAssociative
	CacheAssociativity16waySetAssociative
	CacheAssociativity12waySetAssociative
	CacheAssociativity24waySetAssociative
	CacheAssociativity32waySetAssociative
	CacheAssociativity48waySetAssociative
	CacheAssociativity64waySetAssociative
	CacheAssociativity20waySetAssociative
)

func (a CacheAssociativity) String() string {
	caches := [...]string{
		"Other",
		"Unknown",
		"Direct Mapped",
		"2-way Set-Associative",
		"4-way Set-Associative",
		"Fully Associative",
		"8-way Set-Associative",
		"16-way Set-Associative",
		"12-way Set-Associative",
		"24-way Set-Associative",
		"32-way Set-Associative",
		"48-way Set-Associative",
		"64-way Set-Associative",
		"20-way Set-Associative",
	}
	return caches[a]
}

type CacheInformation struct {
	InfoCommon
	SocketDesignation   string
	Configuration       CacheConfiguration
	MaximumCacheSize    CacheSize
	InstalledSize       CacheSize
	SupportedSRAMType   SRAMType
	CurrentSRAMType     SRAMType
	CacheSpeed          CacheSpeed
	ErrorCorrectionType ErrorCorrectionType
	SystemCacheType     SystemCacheType
	Associativity       CacheAssociativity
}

func (h DMIHeader) CacheInformation() CacheInformation {
	var ci CacheInformation
	data := h.data
	ci.SocketDesignation = h.FieldString(int(data[0x04]))
	ci.Configuration = NewCacheConfiguration(U16(data[0x05:0x07]))
	ci.MaximumCacheSize = NewCacheSize(U16(data[0x07:0x09]))
	ci.InstalledSize = NewCacheSize(U16(data[0x09:0x0B]))
	ci.SupportedSRAMType = SRAMType(U16(data[0x0B:0x0D]))
	ci.CurrentSRAMType = SRAMType(U16(data[0x0D:0x0F]))
	ci.CacheSpeed = CacheSpeed(data[0x0F])
	ci.ErrorCorrectionType = ErrorCorrectionType(data[0x10])
	ci.SystemCacheType = SystemCacheType(data[0x11])
	ci.Associativity = CacheAssociativity(data[0x12])
	return ci
}

type PortConnectorType byte

const (
	PortConnectorTypeNone PortConnectorType = iota
	PortConnectorTypeCentronics
	PortConnectorTypeMiniCentronics
	PortConnectorTypeProprietary
	PortConnectorTypeDB_25pinmale
	PortConnectorTypeDB_25pinfemale
	PortConnectorTypeDB_15pinmale
	PortConnectorTypeDB_15pinfemale
	PortConnectorTypeDB_9pinmale
	PortConnectorTypeDB_9pinfemale
	PortConnectorTypeRJ_11
	PortConnectorTypeRJ_45
	PortConnectorType50_pinMiniSCSI
	PortConnectorTypeMini_DIN
	PortConnectorTypeMicro_DIN
	PortConnectorTypePS2
	PortConnectorTypeInfrared
	PortConnectorTypeHP_HIL
	PortConnectorTypeAccessBusUSB
	PortConnectorTypeSSASCSI
	PortConnectorTypeCircularDIN_8male
	PortConnectorTypeCircularDIN_8female
	PortConnectorTypeOnBoardIDE
	PortConnectorTypeOnBoardFloppy
	PortConnectorType9_pinDualInlinepin10cut
	PortConnectorType25_pinDualInlinepin26cut
	PortConnectorType50_pinDualInline
	PortConnectorType68_pinDualInline
	PortConnectorTypeOnBoardSoundInputfromCD_ROM
	PortConnectorTypeMini_CentronicsType_14
	PortConnectorTypeMini_CentronicsType_26
	PortConnectorTypeMini_jackheadphones
	PortConnectorTypeBNC
	PortConnectorType1394
	PortConnectorTypeSASSATAPlugReceptacle
	PortConnectorTypePC_98
	PortConnectorTypePC_98Hireso
	PortConnectorTypePC_H98
	PortConnectorTypePC_98Note
	PortConnectorTypePC_98Full
	PortConnectorTypeOther
)

func (p PortConnectorType) String() string {
	types := [...]string{
		"None",
		"Centronics",
		"Mini Centronics",
		"Proprietary",
		"DB-25 pin male",
		"DB-25 pin female",
		"DB-15 pin male",
		"DB-15 pin female",
		"DB-9 pin male",
		"DB-9 pin female",
		"RJ-11",
		"RJ-45",
		"50-pin MiniSCSI",
		"Mini-DIN",
		"Micro-DIN",
		"PS/2",
		"Infrared",
		"HP-HIL",
		"Access Bus (USB)",
		"SSA SCSI",
		"Circular DIN-8 male",
		"Circular DIN-8 female",
		"On Board IDE",
		"On Board Floppy",
		"9-pin Dual Inline (pin 10 cut)",
		"25-pin Dual Inline (pin 26 cut)",
		"50-pin Dual Inline",
		"68-pin Dual Inline",
		"On Board Sound Input from CD-ROM",
		"Mini-Centronics Type-14",
		"Mini-Centronics Type-26",
		"Mini-jack (headphones)",
		"BNC",
		"1394",
		"SAS/SATA Plug Receptacle",
		"PC-98",
		"PC-98Hireso",
		"PC-H98",
		"PC-98Note",
		"PC-98Full",
		"Other",
	}
	return types[p]
}

type PortType byte

const (
	PortTypeNone PortType = iota
	PortTypeParallelPortXTATCompatible
	PortTypeParallelPortPS2
	PortTypeParallelPortECP
	PortTypeParallelPortEPP
	PortTypeParallelPortECPEPP
	PortTypeSerialPortXTATCompatible
	PortTypeSerialPort16450Compatible
	PortTypeSerialPort16550Compatible
	PortTypeSerialPort16550ACompatible
	PortTypeSCSIPort
	PortTypeMIDIPort
	PortTypeJoyStickPort
	PortTypeKeyboardPort
	PortTypeMousePort
	PortTypeSSASCSI
	PortTypeUSB
	PortTypeFireWireIEEEP1394
	PortTypePCMCIATypeI2
	PortTypePCMCIATypeII
	PortTypePCMCIATypeIII
	PortTypeCardbus
	PortTypeAccessBusPort
	PortTypeSCSIII
	PortTypeSCSIWide
	PortTypePC_98
	PortTypePC_98_Hireso
	PortTypePC_H98
	PortTypeVideoPort
	PortTypeAudioPort
	PortTypeModemPort
	PortTypeNetworkPort
	PortTypeSATA
	PortTypeSAS
	PortType8251Compatible
	PortType8251FIFOCompatible
	PortTypeOther
)

func (p PortType) String() string {
	types := [...]string{
		"None",
		"Parallel Port XT/AT Compatible",
		"Parallel Port PS/2",
		"Parallel Port ECP",
		"Parallel Port EPP",
		"Parallel Port ECP/EPP",
		"Serial Port XT/AT Compatible",
		"Serial Port 16450 Compatible",
		"Serial Port 16550 Compatible",
		"Serial Port 16550A Compatible",
		"SCSI Port",
		"MIDI Port",
		"Joy Stick Port",
		"Keyboard Port",
		"Mouse Port",
		"SSA SCSI",
		"USB",
		"FireWire (IEEE P1394)",
		"PCMCIA Type I2",
		"PCMCIA Type II",
		"PCMCIA Type III",
		"Cardbus",
		"Access Bus Port",
		"SCSI II",
		"SCSI Wide",
		"PC-98",
		"PC-98-Hireso",
		"PC-H98",
		"Video Port",
		"Audio Port",
		"Modem Port",
		"Network Port",
		"SATA",
		"SAS",
		"8251 Compatible",
		"8251 FIFO Compatible",
		" Other",
	}
	return types[p]
}

type PortInformation struct {
	InfoCommon
	InternalReferenceDesignator string
	InternalConnectorType       PortConnectorType
	ExternalReferenceDesignator string
	ExternalConnectorType       PortConnectorType
	Type                        PortType
}

func (h DMIHeader) PortInformation() PortInformation {
	var pi PortInformation
	data := h.data
	pi.InternalReferenceDesignator = h.FieldString(int(data[0x04]))
	pi.InternalConnectorType = PortConnectorType(data[0x05])
	pi.ExternalReferenceDesignator = h.FieldString(int(data[0x06]))
	pi.ExternalConnectorType = PortConnectorType(data[0x07])
	pi.Type = PortType(data[0x08])
	return pi
}

type SystemSlotType byte

const (
	SystemSlotTypeOther SystemSlotType = 1 + iota
	SystemSlotTypeUnknown
	SystemSlotTypeISA
	SystemSlotTypeMCA
	SystemSlotTypeEISA
	SystemSlotTypePCI
	SystemSlotTypePCCardPCMCIA
	SystemSlotTypeVL_VESA
	SystemSlotTypeProprietary
	SystemSlotTypeProcessorCardSlot
	SystemSlotTypeProprietaryMemoryCardSlot
	SystemSlotTypeIORiserCardSlot
	SystemSlotTypeNuBus
	SystemSlotTypePCI_66MHzCapable
	SystemSlotTypeAGP
	SystemSlotTypeAGP2X
	SystemSlotTypeAGP4X
	SystemSlotTypePCI_X
	SystemSlotTypeAGP8X
	SystemSlotTypePC_98C20
	SystemSlotTypePC_98C24
	SystemSlotTypePC_98E
	SystemSlotTypePC_98LocalBus
	SystemSlotTypePC_98Card
	SystemSlotTypePCIExpress
	SystemSlotTypePCIExpressx1
	SystemSlotTypePCIExpressx2
	SystemSlotTypePCIExpressx4
	SystemSlotTypePCIExpressx8
	SystemSlotTypePCIExpressx16
	SystemSlotTypePCIExpressGen2
	SystemSlotTypePCIExpressGen2x1
	SystemSlotTypePCIExpressGen2x2
	SystemSlotTypePCIExpressGen2x4
	SystemSlotTypePCIExpressGen2x8
	SystemSlotTypePCIExpressGen2x16
	SystemSlotTypePCIExpressGen3
	SystemSlotTypePCIExpressGen3x1
	SystemSlotTypePCIExpressGen3x2
	SystemSlotTypePCIExpressGen3x4
	SystemSlotTypePCIExpressGen3x8
	SystemSlotTypePCIExpressGen3x16
)

func (s SystemSlotType) String() string {
	types := [...]string{
		"Other",
		"Unknown",
		"ISA",
		"MCA",
		"EISA",
		"PCI",
		"PC Card (PCMCIA)",
		"VL-VESA",
		"Proprietary",
		"Processor Card Slot",
		"Proprietary Memory Card Slot",
		"I/O Riser Card Slot",
		"NuBus",
		"PCI – 66MHz Capable",
		"AGP",
		"AGP 2X",
		"AGP 4X",
		"PCI-X",
		"AGP 8X",
		"PC-98/C20",
		"PC-98/C24",
		"PC-98/E",
		"PC-98/Local Bus",
		"PC-98/Card",
		"PCI Express",
		"PCI Express x1",
		"PCI Express x2",
		"PCI Express x4",
		"PCI Express x8",
		"PCI Express x16",
		"PCI Express Gen 2",
		"PCI Express Gen 2 x1",
		"PCI Express Gen 2 x2",
		"PCI Express Gen 2 x4",
		"PCI Express Gen 2 x8",
		"PCI Express Gen 2 x16",
		"PCI Express Gen 3",
		"PCI Express Gen 3 x1",
		"PCI Express Gen 3 x2",
		"PCI Express Gen 3 x4",
		"PCI Express Gen 3 x8",
		"PCI Express Gen 3 x16",
	}
	return types[s-1]
}

type SystemSlotDataBusWidth byte

const (
	SystemSlotDataBusWidthOther SystemSlotDataBusWidth = 1 + iota
	SystemSlotDataBusWidthUnknown
	SystemSlotDataBusWidth8bit
	SystemSlotDataBusWidth16bit
	SystemSlotDataBusWidth32bit
	SystemSlotDataBusWidth64bit
	SystemSlotDataBusWidth128bit
	SystemSlotDataBusWidth1xorx1
	SystemSlotDataBusWidth2xorx2
	SystemSlotDataBusWidth4xorx4
	SystemSlotDataBusWidth8xorx8
	SystemSlotDataBusWidth12xorx12
	SystemSlotDataBusWidth16xorx16
	SystemSlotDataBusWidth32xorx32
)

func (s SystemSlotDataBusWidth) String() string {
	widths := [...]string{
		"Other",
		"Unknown",
		"8 bit",
		"16 bit",
		"32 bit",
		"64 bit",
		"128 bit",
		"1x or x1",
		"2x or x2",
		"4x or x4",
		"8x or x8",
		"12x or x12",
		"16x or x16",
		"32x or x32",
	}
	return widths[s-1]
}

type SystemSlotUsage byte

const (
	SystemSlotUsageOther SystemSlotUsage = 1 + iota
	SystemSlotUsageUnknown
	SystemSlotUsageAvailable
	SystemSlotUsageInuse
)

func (s SystemSlotUsage) String() string {
	usages := [...]string{
		"Other",
		"Unknown",
		"Available",
		"In use",
	}
	return usages[s-1]
}

type SystemSlotLength byte

const (
	SystemSlotLengthOther SystemSlotLength = 1 + iota
	SystemSlotLengthUnknown
	SystemSlotLengthShortLength
	SystemSlotLengthLongLength
)

func (s SystemSlotLength) String() string {
	lengths := [...]string{
		"Other",
		"Unknown",
		"Short Length",
		"Long Length",
	}
	return lengths[s-1]
}

type SystemSlotID uint16

type SystemSlotCharacteristics1 byte

const (
	SystemSlotCharacteristicsunknown SystemSlotCharacteristics1 = 1 << iota
	SystemSlotCharacteristicsProvides5_0volts
	SystemSlotCharacteristicsProvides3_3volts
	SystemSlotCharacteristicsSlotsopeningissharedwithanotherslot
	SystemSlotCharacteristicsPCCardslotsupportsPCCard_16
	SystemSlotCharacteristicsPCCardslotsupportsCardBus
	SystemSlotCharacteristicsPCCardslotsupportsZoomVideo
	SystemSlotCharacteristicsPCCardslotsupportsModemRingResume
)

func (s SystemSlotCharacteristics1) String() string {
	chars := [...]string{
		"Characteristics unknown.",
		"Provides 5.0 volts.",
		"Provides 3.3 volts.",
		"Slot’s opening is shared with another slot (for example, PCI/EISA shared slot).",
		"PC Card slot supports PC Card-16.",
		"PC Card slot supports CardBus.",
		"PC Card slot supports Zoom Video.",
		"PC Card slot supports Modem Ring Resume.",
	}
	return chars[s>>1]
}

type SystemSlotCharacteristics2 byte

const (
	SystemSlotCharacteristics2PCIslotsupportsPowerManagementEventsignal SystemSlotCharacteristics2 = 1 << iota
	SystemSlotCharacteristics2Slotsupportshot_plugdevices
	SystemSlotCharacteristics2PCIslotsupportsSMBussignal
	SystemSlotCharacteristics2Reserved
)

func (s SystemSlotCharacteristics2) String() string {
	chars := [...]string{
		"PCI slot supports Power Management Event (PME#) signal.",
		"Slot supports hot-plug devices.",
		"PCI slot supports SMBus signal.",
		"Reserved",
	}
	return chars[s>>1]
}

type SystemSlotSegmengGroupNumber uint16

type SystemSlotNumber byte

type SystemSlot struct {
	InfoCommon
	Designation          string
	Type                 SystemSlotType
	DataBusWidth         SystemSlotDataBusWidth
	CurrentUsage         SystemSlotUsage
	Length               SystemSlotLength
	ID                   SystemSlotID
	Characteristics1     SystemSlotCharacteristics1
	Characteristics2     SystemSlotCharacteristics2
	SegmentGroupNumber   SystemSlotSegmengGroupNumber
	BusNumber            SystemSlotNumber
	DeviceFunctionNumber SystemSlotNumber
}

func (h DMIHeader) SystemSlot() SystemSlot {
	var ss SystemSlot
	data := h.data
	ss.Designation = h.FieldString(int(data[0x04]))
	ss.Type = SystemSlotType(data[0x05])
	ss.DataBusWidth = SystemSlotDataBusWidth(data[0x06])
	ss.CurrentUsage = SystemSlotUsage(data[0x07])
	ss.Length = SystemSlotLength(data[0x08])
	ss.ID = SystemSlotID(U16(data[0x09:0x0A]))
	ss.Characteristics1 = SystemSlotCharacteristics1(data[0x0B])
	ss.Characteristics2 = SystemSlotCharacteristics2(data[0x0C])
	ss.SegmentGroupNumber = SystemSlotSegmengGroupNumber(U16(data[0x0D:0x0F]))
	ss.BusNumber = SystemSlotNumber(data[0x0F])
	ss.DeviceFunctionNumber = SystemSlotNumber(data[0x10])
	return ss
}

func (s SystemSlot) String() string {
	return fmt.Sprintf("System Slot: %s\n\t\tSlot Designation: %s\n\t\tSlot Type: %s\n\t\tSlot Data Bus Width: %s\n\t\tCurrent Usage: %s\n\t\tSlot Length: %s\n\t\tSlot ID: %s\n\t\tSlot Characteristics1: %s\n\t\tSlot Characteristics2: %s\n\t\tSegment Group Number: %s\n\t\tBus Number: %s\n\t\tDevice/Function Number: %s\n", s.Designation, s.Type, s.DataBusWidth, s.CurrentUsage, s.Length, s.ID, s.Characteristics1, s.Characteristics2, s.SegmentGroupNumber, s.BusNumber, s.DeviceFunctionNumber)
}

type OnBoardDeviceTypeOfDevice byte

const (
	OnBoardDeviceOther OnBoardDeviceTypeOfDevice = 1 + iota
	OnBoardDeviceUnknown
	OnBoardDeviceVideo
	OnBoardDeviceSCSIController
	OnBoardDeviceEthernet
	OnBoardDeviceTokenRing
	OnBoardDeviceSound
	OnBoardDevicePATAController
	OnBoardDeviceSATAController
	OnBoardDeviceSASController
)

func (t OnBoardDeviceTypeOfDevice) String() string {
	types := [...]string{
		"Other",
		"Unknown",
		"Video",
		"SCSI Controller",
		"Ethernet",
		"Token Ring",
		"Sound",
		"PATA Controller",
		"SATA Controller",
		"SAS Controller",
	}
	return types[t-1]
}

type OnBoardDeviceType struct {
	status       bool
	typeOfDevice OnBoardDeviceTypeOfDevice
}

type OnBoardDeviceInformation struct {
	InfoCommon
	Type        []OnBoardDeviceType
	Description []string
}

func (h DMIHeader) OnBoardDeviceInformation() OnBoardDeviceInformation {
	var d OnBoardDeviceInformation
	data := h.data
	n := (data[0x01] - 4) / 2
	for i := byte(1); i <= n; i++ {
		var t OnBoardDeviceType
		index := 4 + 2*(i-1)
		sindex := 5 + 2*(i-1)
		t.status = data[index]&0x80 != 0
		t.typeOfDevice = OnBoardDeviceTypeOfDevice(data[index] & 0x7F)
		d.Type = append(d.Type, t)
		desc := h.FieldString(int(data[sindex]))
		d.Description = append(d.Description, desc)
	}
	return d
}

func (d OnBoardDeviceInformation) String() string {
	var info string
	title := "On Board Devices Information"
	for i, v := range d.Type {
		s := fmt.Sprintf("Device %d: Enabled: %v: Description: %s", i, v.status, v.typeOfDevice, d.Description[i])
		info += "\n\t\t" + s
	}
	return title + "\n\t\t" + info
}

func U16(data []byte) uint16 {
	var u16 uint16
	binary.Read(bytes.NewBuffer(data[0:2]), binary.LittleEndian, &u16)
	return u16
}

func U64(data []byte) uint64 {
	var u64 uint64
	binary.Read(bytes.NewBuffer(data[0:8]), binary.LittleEndian, &u64)
	return u64
}

func NewDMIHeader(data []byte) DMIHeader {
	var h uint16
	binary.Read(bytes.NewBuffer(data[2:4]), binary.LittleEndian, &h)
	hd := DMIHeader{Type: data[0], Length: data[1], Handle: h, data: data}
	return hd
}

func NewSMBIOS_EPS() (SMBIOS_EPS, error) {
	var eps SMBIOS_EPS
	var u16 uint16
	var u32 uint32

	mem, err := getMem(0xF0000, 0x10000)
	if err != nil {
		return eps, err
	}
	data := anchor(mem)
	eps.Anchor = data[:0x04]
	eps.Checksum = data[0x04]
	eps.Length = data[0x05]
	eps.MajorVersion = data[0x06]
	eps.MinorVersion = data[0x07]
	binary.Read(bytes.NewBuffer(data[0x08:0x0A]), binary.LittleEndian, &u16)
	eps.MaxSize = u16
	eps.FormattedArea = data[0x0B:0x0F]
	eps.InterAnchor = data[0x10:0x15]
	eps.InterChecksum = data[0x15]
	binary.Read(bytes.NewBuffer(data[0x16:0x18]), binary.LittleEndian, &u16)
	eps.TableLength = u16
	binary.Read(bytes.NewBuffer(data[0x18:0x1C]), binary.LittleEndian, &u32)
	eps.TableAddress = u32
	binary.Read(bytes.NewBuffer(data[0x1C:0x1E]), binary.LittleEndian, &u16)
	eps.NumberOfSM = u16
	eps.BCDRevision = data[0x1E]
	return eps, nil
}

func (e SMBIOS_EPS) StructureTableMem() ([]byte, error) {
	return getMem(e.TableAddress, uint32(e.TableLength))
}

func (h DMIHeader) Next() DMIHeader {
	de := []byte{0, 0}
	next := h.data[h.Length:]
	index := bytes.Index(next, de)
	hd := NewDMIHeader(next[index+2:])
	return hd
}

func (h DMIHeader) Decode() {
	switch h.Type {
	case 0:
		bi := h.GetBIOSInformation()
		fmt.Println(bi)
	case 1:
		si := h.GetSystemInformation()
		fmt.Println(si)
	case 2:
		bi := h.GetBaseboardInformation()
		fmt.Println(bi)
	case 3:
		ci := h.ChassisInformation()
		fmt.Println(ci)
	case 4:
		pi := h.ProcessorInformation()
		fmt.Println(pi)
	case 7:
		ci := h.CacheInformation()
		fmt.Println(ci)
	case 8:
		pi := h.PortInformation()
		fmt.Println(pi)
	case 9:
		ss := h.SystemSlot()
		fmt.Println(ss)
	case 10:
		di := h.OnBoardDeviceInformation()
		fmt.Println(di)
	default:
		fmt.Println("Unknown")
	}
}

func (h DMIHeader) FieldString(offset int) string {
	d := h.data
	index := int(h.Length)
	if offset == 0 {
		return "Not Specified"
	}
	for i := offset; i > 1 && d[index] != 0; i-- {
		ib := bytes.IndexByte(d[index:], 0)
		if ib != -1 {
			index += ib
			index++
		}
	}
	ib := bytes.IndexByte(d[index:], 0)
	return string(d[index : index+ib])
}

func (h DMIHeader) GetBIOSInformation() BIOSInformation {
	var bi BIOSInformation
	data := h.data
	if h.Type != 0 {
		panic("h.Type is not 0")
	}
	bi.Vendor = h.FieldString(int(data[0x04]))
	bi.BIOSVersion = h.FieldString(int(data[0x05]))
	bi.StartingAddressSegment = U16(data[0x06:0x08])
	bi.ReleaseDate = h.FieldString(int(data[0x08]))
	bi.Characteristics = Characteristics(U64(data[0x0A:0x12]))
	bi.CharacteristicsExt1 = CharacteristicsExt1(data[0x12])
	bi.CharacteristicsExt2 = CharacteristicsExt2(data[0x12])
	return bi
}

func (c Characteristics) String() string {
	var s string
	for i := uint32(4); i < 32; i++ {
		//fmt.Printf("char\n%064b\n%064b\n", char, 1<<i)
		if c&(1<<i) != 0 {
			s += "\n\t\t" + sBIOSCharacteristics[i-3]
		}
	}
	return s
}

func (c CharacteristicsExt1) String() string {
	var s string
	for i := uint32(0); i < 7; i++ {
		if c&(1<<i) != 0 {
			s += "\n\t\t" + sCharateristicsExt1[i]
		}
	}
	return s
}

func (c CharacteristicsExt2) String() string {
	var s string
	for i := uint32(0); i < 5; i++ {
		if c&(1<<i) != 0 {
			s += "\n\t\t" + sCharateristicsExt2[i]
		}
	}
	return s
}

func (bi BIOSInformation) String() string {
	return fmt.Sprintf("BIOS Information\n\tVendor: %s\n\tVersion: %s\n\tAddress: %4X0\n\tCharacteristics: %s\n\tExt1:%s\n\tExt2: %s", bi.Vendor, bi.BIOSVersion, bi.StartingAddressSegment, bi.Characteristics, bi.CharacteristicsExt1, bi.CharacteristicsExt2)
}

func uuid(data []byte, ver string) string {
	if bytes.Index(data, []byte{0x00}) != -1 {
		return "Not present"
	}

	if bytes.Index(data, []byte{0xFF}) != -1 {
		return "Not settable"
	}

	if ver > "2.6" {
		return fmt.Sprintf("%02X%02X%02X%02X-%02X%02X-%02X%02X-%02X%02X-%02X%02X%02X%02X%02X%02X",
			data[3], data[2], data[1], data[0], data[5], data[4], data[7], data[6],
			data[8], data[9], data[10], data[11], data[12], data[13], data[14], data[15])
	}
	return fmt.Sprintf("%02X%02X%02X%02X-%02X%02X-%02X%02X-%02X%02X-%02X%02X%02X%02X%02X%02X",
		data[0], data[1], data[2], data[3], data[4], data[5], data[6], data[7],
		data[8], data[9], data[10], data[11], data[12], data[13], data[14], data[15])
}

func (h DMIHeader) GetSystemInformation() SystemInformation {
	var si SystemInformation
	data := h.data
	if h.Type != 1 {
		panic("Type is not 1")
	}
	si.Manufacturer = h.FieldString(int(data[0x04]))
	si.ProductName = h.FieldString(int(data[0x05]))
	si.Version = h.FieldString(int(data[0x06]))
	si.SerialNumber = h.FieldString(int(data[0x07]))
	si.UUID = uuid(data[0x08:0x18], si.Version)
	si.WakeUpType = WakeUpType(data[0x18])
	si.SKUNumber = h.FieldString(int(data[0x19]))
	si.Family = h.FieldString(int(data[0x1A]))
	return si
}

func (h DMIHeader) GetBaseboardInformation() BaseboardInformation {
	var bi BaseboardInformation
	data := h.data
	if h.Type != 2 {
		panic("Type is not 2")
	}
	bi.Manufacturer = h.FieldString(int(data[0x04]))
	bi.Product = h.FieldString(int(data[0x05]))
	bi.Version = h.FieldString(int(data[0x06]))
	bi.SerailNumber = h.FieldString(int(data[0x07]))
	bi.AssetTag = h.FieldString(int(data[0x08]))
	bi.FeatureFlags = FeatureFlags(data[0x09])
	bi.LocationInChassis = h.FieldString(int(data[0x0A]))
	bi.BoardType = BoardType(data[0x0D])
	return bi
}

func (e SMBIOS_EPS) StructureTable() {
	tmem, err := e.StructureTableMem()
	if err != nil {
		return
	}
	//for i := 0, hd := NewDMIHeader(tmem); i < e.NumberOfSM ; i++, hd = hd.Next() {
	hd := NewDMIHeader(tmem)
	for i := 0; i < 4; i++ {
		hd.Decode()
		hd = hd.Next()
	}
}

func getMem(base uint32, length uint32) ([]byte, error) {
	file, err := os.Open("/dev/mem")
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()
	fd := file.Fd()
	mmoffset := base % uint32(os.Getpagesize())
	mm, err := syscall.Mmap(int(fd), int64(base-mmoffset), int(mmoffset+length), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return []byte{}, err
	}
	mem := make([]byte, len(mm))
	copy(mem, mm)
	err = syscall.Munmap(mm)
	if err != nil {
		return []byte{}, err
	}
	return mem, nil
}

func readMem() ([]byte, error) {
	base := 0xF0000
	file, err := os.Open("/dev/mem")
	if err != nil {
		return []byte{}, err
	}
	fd := file.Fd()
	mmoffset := base % os.Getpagesize()
	mm, err := syscall.Mmap(int(fd), int64(base-mmoffset), mmoffset+0x10000, syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return []byte{}, err
	}
	mem := make([]byte, len(mm))
	copy(mem, mm)
	err = syscall.Munmap(mm)
	if err != nil {
		return []byte{}, err
	}
	return mem, nil
}

func anchor(mem []byte) []byte {
	anchor := []byte{'_', 'S', 'M', '_'}
	i := bytes.Index(mem, anchor)
	return mem[i:]
}

func version(mem []byte) string {
	ver := strconv.Itoa(int(mem[0x06])) + "." + strconv.Itoa(int(mem[0x07]))
	return ver
}

func main() {
	eps, err := NewSMBIOS_EPS()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
	eps.StructureTable()
	//fmt.Printf("%2X", m)
}
