/*
* godmi.go
* DMI SMBIOS information
*
* Chapman Ou <ochapman.cn@gmail.com>
*
 */
package godmi

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"strconv"
	"syscall"
)

var gdmi map[SMBIOSStructureType]interface{}

type SMBIOSStructureType byte

const (
	SMBIOSStructureTypeBIOS SMBIOSStructureType = iota
	SMBIOSStructureTypeSystem
	SMBIOSStructureTypeBaseBoard
	SMBIOSStructureTypeChassis
	SMBIOSStructureTypeProcessor
	SMBIOSStructureTypeMemoryController
	SMBIOSStructureTypeMemoryModule
	SMBIOSStructureTypeCache
	SMBIOSStructureTypePortConnector
	SMBIOSStructureTypeSystemSlots
	SMBIOSStructureTypeOnBoardDevices
	SMBIOSStructureTypeOEMStrings
	SMBIOSStructureTypeSystemConfigurationOptions
	SMBIOSStructureTypeBIOSLanguage
	SMBIOSStructureTypeGroupAssociations
	SMBIOSStructureTypeSystemEventLog
	SMBIOSStructureTypePhysicalMemoryArray
	SMBIOSStructureTypeMemoryDevice
	SMBIOSStructureType32_bitMemoryError
	SMBIOSStructureTypeMemoryArrayMappedAddress
	SMBIOSStructureTypeMemoryDeviceMappedAddress
	SMBIOSStructureTypeBuilt_inPointingDevice
	SMBIOSStructureTypePortableBattery
	SMBIOSStructureTypeSystemReset
	SMBIOSStructureTypeHardwareSecurity
	SMBIOSStructureTypeSystemPowerControls
	SMBIOSStructureTypeVoltageProbe
	SMBIOSStructureTypeCoolingDevice
	SMBIOSStructureTypeTemperatureProbe
	SMBIOSStructureTypeElectricalCurrentProbe
	SMBIOSStructureTypeOut_of_bandRemoteAccess
	SMBIOSStructureTypeBootIntegrityServices
	SMBIOSStructureTypeSystemBoot
	SMBIOSStructureType64_bitMemoryError
	SMBIOSStructureTypeManagementDevice
	SMBIOSStructureTypeManagementDeviceComponent
	SMBIOSStructureTypeManagementDeviceThresholdData
	SMBIOSStructureTypeMemoryChannel
	SMBIOSStructureTypeIPMIDevice
	SMBIOSStructureTypePowerSupply
	SMBIOSStructureTypeAdditionalInformation
	SMBIOSStructureTypeOnBoardDevicesExtendedInformation
	SMBIOSStructureTypeManagementControllerHostInterface                     /*42*/
	SMBIOSStructureTypeInactive                          SMBIOSStructureType = 126
	SMBIOSStructureTypeEndOfTable                        SMBIOSStructureType = 127
)

func (b SMBIOSStructureType) String() string {
	types := [...]string{
		"BIOS", /* 0 */
		"System",
		"Base Board",
		"Chassis",
		"Processor",
		"Memory Controller",
		"Memory Module",
		"Cache",
		"Port Connector",
		"System Slots",
		"On Board Devices",
		"OEM Strings",
		"System Configuration Options",
		"BIOS Language",
		"Group Associations",
		"System Event Log",
		"Physical Memory Array",
		"Memory Device",
		"32-bit Memory Error",
		"Memory Array Mapped Address",
		"Memory Device Mapped Address",
		"Built-in Pointing Device",
		"Portable Battery",
		"System Reset",
		"Hardware Security",
		"System Power Controls",
		"Voltage Probe",
		"Cooling Device",
		"Temperature Probe",
		"Electrical Current Probe",
		"Out-of-band Remote Access",
		"Boot Integrity Services",
		"System Boot",
		"64-bit Memory Error",
		"Management Device",
		"Management Device Component",
		"Management Device Threshold Data",
		"Memory Channel",
		"IPMI Device",
		"Power Supply",
		"Additional Information",
		"Onboard Device",
		"Management Controller Host Interface", /* 42 */
	}
	return types[b]
}

type SMBIOSStructureHandle uint16

type InfoCommon struct {
	Type   SMBIOSStructureType
	Length byte
	Handle SMBIOSStructureHandle
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
	InfoCommon
	data []byte
}

type SMBIOS_Structure struct {
}

type Characteristics uint64
type CharacteristicsExt1 byte
type CharacteristicsExt2 byte

type bIOSInformation struct {
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

type systemInformation struct {
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

func (si systemInformation) String() string {
	return fmt.Sprintf("SystemInformation:"+
		"\n\tManufacturer: %s"+
		"\n\tProduct Name: %s"+
		"\n\tVersion: %s"+
		"\n\tSerial Number: %s"+
		"\n\tUUID: %s"+
		"\n\tWake-up Type: %s"+
		"\n\tSKU Number: %s"+
		"\n\tFamily: %s\n\t",
		si.Manufacturer,
		si.ProductName,
		si.Version,
		si.SerialNumber,
		si.UUID,
		si.WakeUpType,
		si.SKUNumber,
		si.Family)
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

type baseboardInformation struct {
	Type                           byte
	Length                         byte
	Handle                         uint16
	Manufacturer                   string
	Product                        string
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

func (bi baseboardInformation) String() string {
	return fmt.Sprintf("BaseboardInformation:"+
		"\n\tManufacturer: %s"+
		"\n\tProduct: %s"+
		"\n\tVersion: %s"+
		"\n\tSerial Number: %s"+
		"\n\tAsset Tag: %s"+
		"\n\tFeature Flags: %s"+
		"\n\tLocation In Chassis: %s"+
		"\n\tBoard Type: %s\n\t",
		bi.Manufacturer,
		bi.Product,
		bi.Version,
		bi.SerialNumber,
		bi.AssetTag,
		bi.FeatureFlags,
		bi.LocationInChassis,
		bi.BoardType)
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

func (h DMIHeader) ChassisInformation() *ChassisInformation {
	data := h.data
	return &ChassisInformation{
		Manufacturer:                 h.FieldString(int(data[0x04])),
		ChassisType:                  ChassisType(data[0x05]),
		Version:                      h.FieldString(int(data[0x06])),
		SerialNumber:                 h.FieldString(int(data[0x07])),
		AssetTag:                     h.FieldString(int(data[0x08])),
		BootUpState:                  ChassisState(data[0x09]),
		PowerSupplyState:             ChassisState(data[0xA]),
		ThermalState:                 ChassisState(data[0x0B]),
		SecurityStatus:               SecurityStatus(data[0x0C]),
		OEMdefined:                   U16(data[0x0D : 0x0D+4]),
		Height:                       Height(data[0x11]),
		NumberOfPowerCords:           data[0x12],
		ContainedElementCount:        data[0x13],
		ContainedElementRecordLength: data[0x14],
		// TODO: 7.4.4
		//ci.ContainedElements:
		SKUNumber: h.FieldString(int(data[0x15])),
	}
}

func (ci ChassisInformation) String() string {
	return fmt.Sprintf("Chassis Information:\n\t"+
		"Manufacturer: %s"+
		"\n\tType: %s"+
		"\n\tVersion: %s"+
		"\n\tSerial Number: %s"+
		"\n\tAsset Tag: %s"+
		"\n\tBoot-up State: %s"+
		"\n\tPower Supply State: %s"+
		"\n\tThermal State: %s"+
		"\n\tSecurity Status: %s\n\t",
		ci.Manufacturer,
		ci.ChassisType,
		ci.Version,
		ci.SerialNumber,
		ci.AssetTag,
		ci.BootUpState,
		ci.PowerSupplyState,
		ci.ThermalState,
		ci.SecurityStatus)
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
	return fmt.Sprintf("Cache Configuration:"+
		"\n\tLevel: %s"+
		"\n\t\tSocketed: %v"+
		"\n\t\tLocation: %s"+
		"\n\t\tEnabled: %v"+
		"\n\t\tMode:\n\t\t",
		c.Level,
		c.Socketed,
		c.Location,
		c.Enabled,
		c.Mode)
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
	return fmt.Sprintf("System Slot: %s"+
		"\n\t\tSlot Designation: %s"+
		"\n\t\tSlot Type: %s"+
		"\n\t\tSlot Data Bus Width: %s"+
		"\n\t\tCurrent Usage: %s"+
		"\n\t\tSlot Length: %s"+
		"\n\t\tSlot ID: %s"+
		"\n\t\tSlot Characteristics1: %s"+
		"\n\t\tSlot Characteristics2: %s"+
		"\n\t\tSegment Group Number: %s"+
		"\n\t\tBus Number: %s"+
		"\n\t\tDevice/Function Number: %s\n",
		s.Designation,
		s.Type,
		s.DataBusWidth,
		s.CurrentUsage,
		s.Length,
		s.ID,
		s.Characteristics1,
		s.Characteristics2,
		s.SegmentGroupNumber,
		s.BusNumber,
		s.DeviceFunctionNumber)
}

type BIOSLanguageInformationFlag byte

const (
	BIOSLanguageInformationFlagLongFormat BIOSLanguageInformationFlag = iota
	BIOSLanguageInformationFlagAbbreviatedFormat
)

func NewBIOSLanguageInformationFlag(f byte) BIOSLanguageInformationFlag {
	return BIOSLanguageInformationFlag(f & 0xFE)
}

type BIOSLanguageInformation struct {
	InfoCommon
	InstallableLanguage []string
	Flags               BIOSLanguageInformationFlag
	CurrentLanguage     string
}

func (h DMIHeader) BIOSLanguageInformation() BIOSLanguageInformation {
	var bl BIOSLanguageInformation
	data := h.data
	cnt := data[0x04]
	for i := byte(1); i <= cnt; i++ {
		bl.InstallableLanguage = append(bl.InstallableLanguage, h.FieldString(int(data[i])))
	}
	bl.Flags = NewBIOSLanguageInformationFlag(data[0x05])
	bl.CurrentLanguage = bl.InstallableLanguage[data[0x15]]
	return bl
}

func (b BIOSLanguageInformation) String() string {
	return fmt.Sprintf("BIOS Language Information:"+
		"\n\t\tInstallable Languages %s"+
		"\n\t\tFlags: %s"+
		"\n\t\tCurrent Language: %s\n",
		b.InstallableLanguage,
		b.Flags,
		b.CurrentLanguage)
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

type SystemConfigurationOptions struct {
	InfoCommon
	Count   byte
	strings string
}

//Type 11
type OEMStrings struct {
	InfoCommon
	Count   byte
	strings string
}

func (h DMIHeader) SystemConfigurationOptions() SystemConfigurationOptions {
	var sc SystemConfigurationOptions
	data := h.data
	sc.Count = data[0x04]
	for i := byte(1); i <= sc.Count; i++ {
		sc.strings += fmt.Sprintf("string %d: %s\n\t\t", i, h.FieldString(int(data[0x04+i])))
	}
	return sc
}

func (s SystemConfigurationOptions) String() string {
	return fmt.Sprintf("System Configuration Option\n\t\t%s", s.strings)
}

func (h DMIHeader) OEMStrings() OEMStrings {
	var o OEMStrings
	data := h.data
	o.Count = data[0x04]
	for i := byte(0); i < o.Count; i++ {
		o.strings += fmt.Sprintf("strings: %d %s\n\t\t", i, h.FieldString(int(data[i])))
	}
	return o
}

func (o OEMStrings) String() string {
	return fmt.Sprintf("OEM strings: %s", o.strings)
}

type GroupAssociationsItem struct {
	Type   SMBIOSStructureType
	Handle SMBIOSStructureHandle
}

// Type 14
type GroupAssociations struct {
	InfoCommon
	GroupName string
	Item      []GroupAssociationsItem
}

func (h DMIHeader) GroupAssociations() GroupAssociations {
	var ga GroupAssociations
	data := h.data
	ga.GroupName = h.FieldString(int(data[0x04]))
	cnt := (h.Length - 5) / 3
	items := data[5:]
	var i byte
	for i = 0; i < cnt; i++ {
		var gai GroupAssociationsItem
		gai.Type = SMBIOSStructureType(items[i*3])
		gai.Handle = SMBIOSStructureHandle(U16(items[i*3+1:]))
		ga.Item = append(ga.Item, gai)
	}
	return ga
}

func (g GroupAssociations) String() string {
	return fmt.Sprintf("Group Associations:"+
		"\n\t\tGroup Name: %s"+
		"\n\t\tItem: %#v\n",
		g.GroupName,
		g.Item)
}

type PhysicalMemoryArrayLocation byte

const (
	PhysicalMemoryArrayLocationOther PhysicalMemoryArrayLocation = 1 + iota
	PhysicalMemoryArrayLocationUnknown
	PhysicalMemoryArrayLocationSystemboardormotherboard
	PhysicalMemoryArrayLocationISAadd_oncard
	PhysicalMemoryArrayLocationEISAadd_oncard
	PhysicalMemoryArrayLocationPCIadd_oncard
	PhysicalMemoryArrayLocationMCAadd_oncard
	PhysicalMemoryArrayLocationPCMCIAadd_oncard
	PhysicalMemoryArrayLocationProprietaryadd_oncard
	PhysicalMemoryArrayLocationNuBus
	PhysicalMemoryArrayLocationPC_98C20add_oncard
	PhysicalMemoryArrayLocationPC_98C24add_oncard
	PhysicalMemoryArrayLocationPC_98Eadd_oncard
	PhysicalMemoryArrayLocationPC_98Localbusadd_oncard
)

func (p PhysicalMemoryArrayLocation) String() string {
	locations := [...]string{
		"Other",
		"Unknown",
		"System board or motherboard",
		"ISA add-on card",
		"EISA add-on card",
		"PCI add-on card",
		"MCA add-on card",
		"PCMCIA add-on card",
		"Proprietary add-on card",
		"NuBus",
		"PC-98/C20 add-on card",
		"PC-98/C24 add-on card",
		"PC-98/E add-on card",
		"PC-98/Local bus add-on card",
	}
	return locations[p-1]
}

type PhysicalMemoryArrayUse byte

const (
	PhysicalMemoryArrayUseOther PhysicalMemoryArrayUse = 1 + iota
	PhysicalMemoryArrayUseUnknown
	PhysicalMemoryArrayUseSystemmemory
	PhysicalMemoryArrayUseVideomemory
	PhysicalMemoryArrayUseFlashmemory
	PhysicalMemoryArrayUseNon_volatileRAM
	PhysicalMemoryArrayUseCachememory
)

func (p PhysicalMemoryArrayUse) String() string {
	uses := [...]string{
		"Other",
		"Unknown",
		"System memory",
		"Video memory",
		"Flash memory",
		"Non-volatile RAM",
		"Cache memory",
	}
	return uses[p-1]
}

type PhysicalMemoryArrayErrorCorrection byte

const (
	PhysicalMemoryArrayErrorCorrectionOther PhysicalMemoryArrayErrorCorrection = 1 + iota
	PhysicalMemoryArrayErrorCorrectionUnknown
	PhysicalMemoryArrayErrorCorrectionNone
	PhysicalMemoryArrayErrorCorrectionParity
	PhysicalMemoryArrayErrorCorrectionSingle_bitECC
	PhysicalMemoryArrayErrorCorrectionMulti_bitECC
	PhysicalMemoryArrayErrorCorrectionCRC
)

func (p PhysicalMemoryArrayErrorCorrection) String() string {
	types := [...]string{
		"Other",
		"Unknown",
		"None",
		"Parity",
		"Single-bit ECC",
		"Multi-bit ECC",
		"CRC",
	}
	return types[p-1]
}

type PhysicalMemoryArray struct {
	InfoCommon
	Location                PhysicalMemoryArrayLocation
	Use                     PhysicalMemoryArrayUse
	ErrorCorrection         PhysicalMemoryArrayErrorCorrection
	MaximumCapacity         uint32
	ErrorInformationHandle  uint16
	NumberOfMemoryDevices   uint16
	ExtendedMaximumCapacity uint64
}

func (h DMIHeader) PhysicalMemoryArray() PhysicalMemoryArray {
	var pma PhysicalMemoryArray
	data := h.data
	pma.Location = PhysicalMemoryArrayLocation(data[0x04])
	pma.Use = PhysicalMemoryArrayUse(data[0x05])
	pma.ErrorCorrection = PhysicalMemoryArrayErrorCorrection(data[0x06])
	pma.MaximumCapacity = U32(data[0x07:0x0B])
	pma.ErrorInformationHandle = U16(data[0x0B:0x0D])
	pma.NumberOfMemoryDevices = U16(data[0x0D:0x0F])
	pma.ExtendedMaximumCapacity = U64(data[0x0F:])
	return pma
}

func (p PhysicalMemoryArray) String() string {
	return fmt.Sprintf("Physcial Memory Array:\n\t\t"+
		"Location: %s\n\t\t"+
		"Use: %s\n\t\t"+
		"Memory Error Correction: %s\n\t\t"+
		"Maximum Capacity: %d\n\t\t"+
		"Memory Error Information Handle: %d\n\t\t"+
		"Number of Memory Devices: %d\n\t\t"+
		"Extended Maximum Capacity: %d\n",
		p.Location,
		p.Use,
		p.ErrorCorrection,
		p.MaximumCapacity,
		p.ErrorInformationHandle,
		p.NumberOfMemoryDevices,
		p.ExtendedMaximumCapacity)
}

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
	InfoCommon
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

func (h DMIHeader) MemoryDevice() MemoryDevice {
	var md MemoryDevice
	data := h.data
	md.PhysicalMemoryArrayHandle = U16(data[0x04:0x06])
	md.ErrorInformationHandle = U16(data[0x06:0x08])
	md.TotalWidth = U16(data[0x08:0x0A])
	md.DataWidth = U16(data[0x0A:0x0C])
	md.Size = U16(data[0x0C:0x0e])
	md.FormFactor = MemoryDeviceFormFactor(data[0x0E])
	md.DeviceSet = data[0x0F]
	md.DeviceLocator = h.FieldString(int(data[0x10]))
	md.BankLocator = h.FieldString(int(data[0x11]))
	md.Type = MemoryDeviceType(data[0x12])
	md.TypeDetail = MemoryDeviceTypeDetail(U16(data[0x13:0x15]))
	md.Speed = U16(data[0x15:0x17])
	md.Manufacturer = h.FieldString(int(data[0x17]))
	md.SerialNumber = h.FieldString(int(data[0x18]))
	md.PartNumber = h.FieldString(int(data[0x1A]))
	md.Attributes = data[0x1B]
	md.ExtendedSize = U32(data[0x1C:0x20])
	md.ConfiguredVoltage = U16(data[0x20:0x22])
	md.MinimumVoltage = U16(data[0x22:0x24])
	md.MaximumVoltage = U16(data[0x24:0x26])
	md.ConfiguredVoltage = U16(data[0x26:0x28])
	return md
}

func (m MemoryDevice) String() string {
	return fmt.Sprintf("Memory Device:\n\t\t"+
		"Physical Memory Array Handle: %d\n\t\t"+
		"Memory Error Information Handle: %d\n\t\t"+
		"Total Width: %d\n\t\t"+
		"Data Width: %d\n\t\t"+
		"Size: %d\n\t\t"+
		"Form Factor: %s\n\t\t"+
		"Device Set: %d\n\t\t"+
		"Device Locator: %s\n\t\t"+
		"Bank Locator: %s\n\t\t"+
		"Memory Type: %s\n\t\t"+
		"Type Detail: %s\n\t\t"+
		"Speed: %d\n\t\t"+
		"Manufacturer: %s\n\t\t"+
		"Serial Number: %s\n\t\t"+
		"Asset Tag: %s\n\t\t"+
		"Part Number: %s\n\t\t"+
		"Attributes: %s\n\t\t"+
		"Extended Size: %s\n\t\t"+
		"Configured Memory Clock Speed: %d\n\t\t"+
		"Minimum voltage: %d\n\t\t"+
		"Maximum voltage: %d\n\t\t"+
		"Configured voltage: %d\n\t\t",
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
	InfoCommon
	Type              MemoryErrorInformationType
	Granularity       MemoryErrorInformationGranularity
	Operation         MemoryErrorInformationOperation
	VendorSyndrome    uint32
	ArrayErrorAddress uint32
	ErrorAddress      uint32
	Resolution        uint32
}

func (h DMIHeader) _32BitMemoryErrorInformation() _32BitMemoryErrorInformation {
	var mei _32BitMemoryErrorInformation
	data := h.data
	mei.Type = MemoryErrorInformationType(data[0x04])
	mei.Granularity = MemoryErrorInformationGranularity(data[0x05])
	mei.Operation = MemoryErrorInformationOperation(data[0x06])
	mei.VendorSyndrome = U32(data[0x07:0x0B])
	mei.ArrayErrorAddress = U32(data[0x0B:0x0F])
	mei.ErrorAddress = U32(data[0x0F:0x13])
	mei.Resolution = U32(data[0x13:0x22])
	return mei
}

func (m _32BitMemoryErrorInformation) String() string {
	return fmt.Sprintf("32 Bit Memory Error Information:\n\t\t"+
		"Error Type: %s\n\t\t"+
		"Error Granularity: %s\n\t\t"+
		"Error Operation: %s\n\t\t"+
		"Vendor Syndrome: %d\n\t\t"+
		"Memory Array Error Address: %d\n\t\t"+
		"Device Error Address: %d\n\t\t"+
		"Error Resoluton: %d\n\t\t",
		m.Type,
		m.Granularity,
		m.Operation,
		m.VendorSyndrome,
		m.ArrayErrorAddress,
		m.ErrorAddress,
		m.Resolution,
	)
}

type BuiltinPointingDeviceType byte

const (
	BuiltinPointingDeviceTypeOther BuiltinPointingDeviceType = 1 + iota
	BuiltinPointingDeviceTypeUnknown
	BuiltinPointingDeviceTypeMouse
	BuiltinPointingDeviceTypeTrackBall
	BuiltinPointingDeviceTypeTrackPoint
	BuiltinPointingDeviceTypeGlidePoint
	BuiltinPointingDeviceTypeTouchPad
	BuiltinPointingDeviceTypeTouchScreen
	BuiltinPointingDeviceTypeOpticalSensor
)

func (b BuiltinPointingDeviceType) String() string {
	types := [...]string{
		"Other",
		"Unknown",
		"Mouse",
		"Track Ball",
		"Track Point",
		"Glide Point",
		"Touch Pad",
		"Touch Screen",
		"Optical Sensor",
	}
	return types[b-1]
}

type BuiltinPointingDeviceInterface byte

const (
	BuiltinPointingDeviceInterfaceOther BuiltinPointingDeviceInterface = 1 + iota
	BuiltinPointingDeviceInterfaceUnknown
	BuiltinPointingDeviceInterfaceSerial
	BuiltinPointingDeviceInterfacePS2
	BuiltinPointingDeviceInterfaceInfrared
	BuiltinPointingDeviceInterfaceHP_HIL
	BuiltinPointingDeviceInterfaceBusmouse
	BuiltinPointingDeviceInterfaceADB
	BuiltinPointingDeviceInterfaceBusmouseDB_9
	BuiltinPointingDeviceInterfaceBusmousemicro_DIN
	BuiltinPointingDeviceInterfaceUSB
)

func (b BuiltinPointingDeviceInterface) String() string {
	interfaces := [...]string{
		"Other",
		"Unknown",
		"Serial",
		"PS/2",
		"Infrared",
		"HP-HIL",
		"Bus mouse",
		"ADB (Apple Desktop Bus)",
		"Bus mouse DB-9",
		"Bus mouse micro-DIN",
		"USB",
	}
	return interfaces[b-1]
}

type BuiltinPointingDevice struct {
	InfoCommon
	Type            BuiltinPointingDeviceType
	Interface       BuiltinPointingDeviceInterface
	NumberOfButtons byte
}

func (h DMIHeader) BuiltinPointingDevice() BuiltinPointingDevice {
	var b BuiltinPointingDevice
	data := h.data
	b.Type = BuiltinPointingDeviceType(data[0x04])
	b.Interface = BuiltinPointingDeviceInterface(data[0x05])
	b.NumberOfButtons = data[0x06]
	return b
}

func (b BuiltinPointingDevice) String() string {
	return fmt.Sprintf("Built-in Pointing Device:\n\t\t"+
		"Type: %s\n\t\t"+
		"Interface: %s\n\t\t"+
		"Number of Buttons: %d\n",
		b.Type,
		b.Interface,
		b.NumberOfButtons,
	)
}

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
	InfoCommon
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

func (h DMIHeader) PortableBattery() PortableBattery {
	var p PortableBattery
	data := h.data
	p.Location = h.FieldString(int(data[0x04]))
	p.Manufacturer = h.FieldString(int(data[0x05]))
	p.ManufacturerDate = h.FieldString(int(data[0x06]))
	p.SerialNumber = h.FieldString(int(data[0x07]))
	p.DeviceName = h.FieldString(int(data[0x08]))
	p.DeviceChemistry = PortableBatteryDeviceChemistry(data[0x09])
	p.DesignCapacity = U16(data[0x0A:0x0C])
	p.DesignVoltage = U16(data[0x0C:0x0E])
	p.SBDSVersionNumber = h.FieldString(int(data[0x0E]))
	p.MaximumErrorInBatteryData = data[0x0F]
	p.SBDSSerialNumber = U16(data[0x10:0x12])
	p.SBDSManufactureDate = U16(data[0x12:0x14])
	p.SBDSDeviceChemistry = h.FieldString(int(data[0x14]))
	p.DesignCapacityMultiplier = data[0x15]
	p.OEMSepecific = U32(data[0x16:0x1A])
	return p
}

func (p PortableBattery) String() string {
	return fmt.Sprintf("Portable Battery\n\t\t"+
		"Location: %s\n\t\t"+
		"Manufacturer: %s\n\t\t"+
		"Manufacturer Date: %s\n\t\t"+
		"Serial Number: %s\n\t\t"+
		"Device Name: %s\n\t\t"+
		"Device Chemistry: %s\n\t\t"+
		"Design Capacity: %d\n\t\t"+
		"Design Voltage: %d\n\t\t"+
		"SBDS Version Number: %s\n\t\t"+
		"Maximum Error in Battery Data: %d\n\t\t"+
		"SBDS Serial Numberd: %d\n\t\t"+
		"SBDS Manufacturer Date: %d\n\t\t"+
		"SBDS Device Chemistry: %s\n\t\t"+
		"Design Capacity Multiplier: %d\n\t\t"+
		"OEM-specific: %d\n",
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

type SystemResetBootOption byte

const (
	SystemResetBootOptionReserved SystemResetBootOption = iota
	SystemResetBootOptionOperatingSystem
	SystemResetBootOptionSystemUtilities
	SystemResetBootOptionDoNotReboot
)

func (s SystemResetBootOption) String() string {
	options := [...]string{
		"Reserved",
		"Operating System",
		"System Utilities",
		"Do Not Reboot",
	}
	return options[s]
}

type SystemResetCapabilities struct {
	Status            bool
	BootOptionOnLimit SystemResetBootOption
	BootOption        SystemResetBootOption
	WatchdogTimer     bool
}

func NewSystemResetCapablities(data byte) SystemResetCapabilities {
	var s SystemResetCapabilities
	s.Status = (data&0x01 != 0)
	s.BootOption = SystemResetBootOption(data & 0x06)
	s.BootOptionOnLimit = SystemResetBootOption(data & 0x18)
	s.WatchdogTimer = data&0x20 != 0
	return s
}

func (s SystemResetCapabilities) String() string {
	return fmt.Sprintf("Capablities:\n\t\t"+
		"Status: %t\n\t\t"+
		"Boot Option: %s\n\t\t"+
		"Boot Option On Limit: %s\n\t\t"+
		"Watchdog Timer: %t\n",
		s.Status,
		s.BootOption,
		s.BootOptionOnLimit,
		s.WatchdogTimer)
}

type SystemReset struct {
	InfoCommon
	Capabilities  byte
	ResetCount    uint16
	ResetLimit    uint16
	TimerInterval uint16
	Timeout       uint16
}

func (s SystemReset) String() string {
	return fmt.Sprintf("System Reset:\n\t\t"+
		"Capabilities: %s\n\t\t"+
		"Reset Count: %d\n\t\t"+
		"Reset Limit: %d\n\t\t"+
		"Timer Interval: %d\n\t\t"+
		"Timeout: %d\n",
		s.Capabilities,
		s.ResetCount,
		s.ResetLimit,
		s.TimerInterval,
		s.Timeout)
}

func (h DMIHeader) SystemReset() SystemReset {
	var s SystemReset
	data := h.data
	s.Capabilities = data[0x04]
	s.ResetCount = U16(data[0x05:0x07])
	s.ResetLimit = U16(data[0x07:0x09])
	s.TimerInterval = U16(data[0x09:0x0B])
	s.Timeout = U16(data[0x0B:0x0D])
	return s
}

type HardwareSecurityStatus byte

const (
	HardwareSecurityStatusDisabled HardwareSecurityStatus = iota
	HardwareSecurityStatusEnabled
	HardwareSecurityStatusNotImplemented
	HardwareSecurityStatusUnknown
)

func (h HardwareSecurityStatus) String() string {
	status := [...]string{
		"Disabled",
		"Enabled",
		"Not Implemented",
		"Unknown",
	}
	return status[h]
}

type HardwareSecuritySettings struct {
	PowerOnPassword       HardwareSecurityStatus
	KeyboardPassword      HardwareSecurityStatus
	AdministratorPassword HardwareSecurityStatus
	FrontPanelReset       HardwareSecurityStatus
}

func NewHardwareSecurity(data byte) HardwareSecuritySettings {
	var h HardwareSecuritySettings
	h.PowerOnPassword = HardwareSecurityStatus(data & 0xC0)
	h.KeyboardPassword = HardwareSecurityStatus(data & 0x30)
	h.AdministratorPassword = HardwareSecurityStatus(data & 0x0C)
	h.FrontPanelReset = HardwareSecurityStatus(data & 0x03)
	return h
}

func (h HardwareSecuritySettings) String() string {
	return fmt.Sprintf("Power-on Password Status: %s\n"+
		"Keyboard Password Status: %s\n"+
		"Administrator Password Status: %s\n"+
		"Front Panel Reset Status: %s\n",
		h.PowerOnPassword,
		h.KeyboardPassword,
		h.AdministratorPassword,
		h.FrontPanelReset)
}

type HardwareSecurity struct {
	InfoCommon
	Setting HardwareSecuritySettings
}

func (h DMIHeader) HardwareSecurity() HardwareSecurity {
	var hw HardwareSecurity
	data := h.data
	hw.Setting = NewHardwareSecurity(data[0x04])
	return hw
}

func (h HardwareSecurity) String() string {
	return fmt.Sprintf("Hardware Security\n\t\t"+
		"Setting: %s\n\t\t",
		h.Setting)
}

type SystemPowerControlsMonth byte
type SystemPowerControlsDayOfMonth byte
type SystemPowerControlsHour byte
type SystemPowerControlsMinute byte
type SystemPowerControlsSecond byte

type SystemPowerControls struct {
	InfoCommon
	NextScheduledPowerOnMonth      SystemPowerControlsMonth
	NextScheduledPowerOnDayOfMonth SystemPowerControlsDayOfMonth
	NextScheduledPowerOnHour       SystemPowerControlsHour
	NextScheduledPowerMinute       SystemPowerControlsMinute
	NextScheduledPowerSecond       SystemPowerControlsSecond
}

func (h DMIHeader) SystemPowerControls() *SystemPowerControls {
	data := h.data
	return &SystemPowerControls{
		NextScheduledPowerOnMonth:      SystemPowerControlsMonth(bcd(data[0x04:0x05])),
		NextScheduledPowerOnDayOfMonth: SystemPowerControlsDayOfMonth(bcd(data[0x05:0x06])),
		NextScheduledPowerOnHour:       SystemPowerControlsHour(bcd(data[0x06:0x07])),
		NextScheduledPowerMinute:       SystemPowerControlsMinute(bcd(data[0x07:0x08])),
		NextScheduledPowerSecond:       SystemPowerControlsSecond(bcd(data[0x08:0x09])),
	}
}

func (s SystemPowerControls) String() string {
	return fmt.Sprintf("System Power Controls:\n\t\t"+
		"Next Scheduled Power-on Month: %d"+
		"Next Scheduled Power-on Day-of-month: %d"+
		"Next Scheduled Power-on Hour: %d"+
		"Next Scheduled Power-on Minute: %d"+
		"Next Scheduled Power-on Second: %d",
		s.NextScheduledPowerOnMonth,
		s.NextScheduledPowerOnDayOfMonth,
		s.NextScheduledPowerOnHour,
		s.NextScheduledPowerMinute,
		s.NextScheduledPowerSecond)
}

type VoltageProbeStatus byte

const (
	VoltageProbeStatusOther VoltageProbeStatus = 0x20 + iota
	VoltageProbeStatusUnknown
	VoltageProbeStatusOK
	VoltageProbeStatusNon_critical
	VoltageProbeStatusCritical
	VoltageProbeStatusNon_recoverable
)

func (v VoltageProbeStatus) String() string {
	status := [...]string{
		"Other",
		"Unknown",
		"OK",
		"Non-critical",
		"Critical",
		"Non-recoverable",
	}
	return status[v-6]
}

type VoltageProbeLocation byte

const (
	VoltageProbeLocationOther VoltageProbeLocation = 1 + iota
	VoltageProbeLocationUnknown
	VoltageProbeLocationOK
	VoltageProbeLocationNon_critical
	VoltageProbeLocationCritical
	VoltageProbeLocationNon_recoverable
	VoltageProbeLocationMotherboard
	VoltageProbeLocationMemoryModule
	VoltageProbeLocationProcessorModule
	VoltageProbeLocationPowerUnit
	VoltageProbeLocationAdd_inCard
)

func (v VoltageProbeLocation) String() string {
	locations := [...]string{
		"Other",
		"Unknown",
		"OK",
		"Non-critical",
		"Critical",
		"Non-recoverable",
		"Motherboard",
		"Memory Module",
		"Processor Module",
		"Power Unit",
		"Add-in Card",
	}
	return locations[v-1]
}

type VoltageProbeLocationAndStatus struct {
	Status   VoltageProbeStatus
	Location VoltageProbeLocation
}

func NewVoltageProbeLocationAndStatus(data byte) VoltageProbeLocationAndStatus {
	return VoltageProbeLocationAndStatus{
		Status:   VoltageProbeStatus(data & 0x1F),
		Location: VoltageProbeLocation(data & 0xE0),
	}
}

func (v VoltageProbeLocationAndStatus) String() string {
	return fmt.Sprintf("\n\t\t\t\tStatus: %s\n\t\t\t\tLocation: %s",
		v.Status, v.Location)
}

type VoltageProbe struct {
	InfoCommon
	Description       string
	LocationAndStatus VoltageProbeLocationAndStatus
	MaximumValue      uint16
	MinimumValude     uint16
	Resolution        uint16
	Tolerance         uint16
	Accuracy          uint16
	OEMdefined        uint16
	NominalValue      uint16
}

func (v VoltageProbe) String() string {
	return fmt.Sprintf("Voltage Probe:\n\t\t"+
		"Description: %s\n\t\t"+
		"Location And Status: %s\n\t\t"+
		"Maximum Value: %d\n\t\t"+
		"Minimum Valude: %d\n\t\t"+
		"Resolution: %d\n\t\t"+
		"Tolerance: %d\n\t\t"+
		"Accuracy: %d\n\t\t"+
		"OE Mdefined: %d\n\t\t"+
		"Nominal Value: %d\n",
		v.Description,
		v.LocationAndStatus,
		v.MaximumValue,
		v.MinimumValude,
		v.Resolution,
		v.Tolerance,
		v.Accuracy,
		v.OEMdefined,
		v.NominalValue)
}

func (h DMIHeader) VoltageProbe() *VoltageProbe {
	data := h.data
	return &VoltageProbe{
		Description:       h.FieldString(int(data[0x04])),
		LocationAndStatus: NewVoltageProbeLocationAndStatus(data[0x05]),
		MaximumValue:      U16(data[0x06:0x08]),
		MinimumValude:     U16(data[0x08:0x0A]),
		Resolution:        U16(data[0x0A:0x0C]),
		Tolerance:         U16(data[0x0C:0x0E]),
		Accuracy:          U16(data[0x0E:0x10]),
		OEMdefined:        U16(data[0x10:0x12]),
		NominalValue:      U16(data[0x12:0x14]),
	}
}

type CoolingDeviceStatus byte

const (
	CoolingDeviceStatusOther CoolingDeviceStatus = 0x20 + iota
	CoolingDeviceStatusUnknown
	CoolingDeviceStatusOK
	CoolingDeviceStatusNon_critical
	CoolingDeviceStatusCritical
	CoolingDeviceStatusNon_recoverable
)

func (c CoolingDeviceStatus) String() string {
	status := [...]string{
		"Other",
		"Unknown",
		"OK",
		"Non-critical",
		"Critical",
		"Non-recoverable",
	}
	return status[c-0x20]
}

type CoolingDeviceType byte

const (
	CoolingDeviceTypeOther CoolingDeviceType = 1 + iota
	CoolingDeviceTypeUnknown
	CoolingDeviceTypeFan
	CoolingDeviceTypeCentrifugalBlower
	CoolingDeviceTypeChipFan
	CoolingDeviceTypeCabinetFan
	CoolingDeviceTypePowerSupplyFan
	CoolingDeviceTypeHeatPipe
	CoolingDeviceTypeIntegratedRefrigeration
	CoolingDeviceTypeActiveCooling
	CoolingDeviceTypePassiveCooling
)

func (c CoolingDeviceType) String() string {
	types := [...]string{
		"Other",
		"Unknown",
		"Fan",
		"Centrifugal Blower",
		"Chip Fan",
		"Cabinet Fan",
		"Power Supply Fan",
		"Heat Pipe",
		"Integrated Refrigeration",
		"Active Cooling",
		"Passive Cooling",
	}
	return types[c-1]
}

type CoolingDeviceTypeAndStatus struct {
	Status CoolingDeviceStatus
	Type   CoolingDeviceType
}

func NewCoolingDeviceTypeAndStatus(data byte) CoolingDeviceTypeAndStatus {
	return CoolingDeviceTypeAndStatus{
		Status: CoolingDeviceStatus(data & 0xE0),
		Type:   CoolingDeviceType(data & 0x1F),
	}
}

type CoolingDevice struct {
	InfoCommon
	TemperatureProbeHandle uint16
	DeviceTypeAndStatus    CoolingDeviceTypeAndStatus
	CoolingUintGroup       byte
	OEMdefined             uint32
	NominalSpeed           uint16
	Description            string
}

func (c CoolingDevice) String() string {
	s := fmt.Sprintf("Cooling Device:\n\t\t"+
		"Temperature Probe Handle: %d\n\t\t"+
		"Device Type And Status: %s\n\t\t"+
		"Cooling Uint Group: %d\n\t\t"+
		"OE Mdefined: %d\n\t\t",
		c.TemperatureProbeHandle,
		c.DeviceTypeAndStatus,
		c.CoolingUintGroup,
		c.OEMdefined,
	)
	if c.Length > 0x0C {
		s += fmt.Sprintf("Nominal Speed: %d\n\t\t", c.NominalSpeed)
	}
	if c.Length > 0x0F {
		s += fmt.Sprintf("Description: %s\n", c.Description)
	}
	return s
}

func (h DMIHeader) CoolingDevice() *CoolingDevice {
	data := h.data
	cd := &CoolingDevice{
		TemperatureProbeHandle: U16(data[0x04:0x06]),
		DeviceTypeAndStatus:    NewCoolingDeviceTypeAndStatus(data[0x06]),
		CoolingUintGroup:       data[0x07],
		OEMdefined:             U32(data[0x08:0x0C]),
	}
	if h.Length > 0x0C {
		cd.NominalSpeed = U16(data[0x0C:0x0E])
	}
	if h.Length > 0x0F {
		cd.Description = h.FieldString(int(data[0x0E]))
	}
	return cd
}

type TemperatureProbeStatus byte

const (
	TemperatureProbeStatusOther TemperatureProbeStatus = 0x20 + iota
	TemperatureProbeStatusUnknown
	TemperatureProbeStatusOK
	TemperatureProbeStatusNon_critical
	TemperatureProbeStatusCritical
	TemperatureProbeStatusNon_recoverable
)

func (t TemperatureProbeStatus) String() string {
	status := [...]string{
		"Other",
		"Unknown",
		"OK",
		"Non-critical",
		"Critical",
		"Non-recoverable",
	}
	return status[t-0x20]
}

type TemperatureProbeLocation byte

const (
	TemperatureProbeLocationOther TemperatureProbeStatus = 1 + iota
	TemperatureProbeLocationUnknown
	TemperatureProbeLocationProcessor
	TemperatureProbeLocationDisk
	TemperatureProbeLocationPeripheralBay
	TemperatureProbeLocationSystemManagementModule
	TemperatureProbeLocationMotherboard
	TemperatureProbeLocationMemoryModule
	TemperatureProbeLocationProcessorModule
	TemperatureProbeLocationPowerUnit
	TemperatureProbeLocationAdd_inCard
	TemperatureProbeLocationFrontPanelBoard
	TemperatureProbeLocationBackPanelBoard
	TemperatureProbeLocationPowerSystemBoard
	TemperatureProbeLocationDriveBackPlane
)

func (t TemperatureProbeLocation) String() string {
	locations := [...]string{
		"Other",
		"Unknown",
		"Processor",
		"Disk",
		"Peripheral Bay",
		"System Management Module",
		"Motherboard",
		"Memory Module",
		"Processor Module",
		"Power Unit",
		"Add-in Card",
		"Front Panel Board",
		"Back Panel Board",
		"Power System Board",
		"Drive Back Plane",
	}
	return locations[t-1]
}

type TemperatureProbeLocationAndStatus struct {
	Status   TemperatureProbeStatus
	Location TemperatureProbeLocation
}

func (t TemperatureProbeLocationAndStatus) String() string {
	return fmt.Sprintf("\n\t\t\t\tStatus: %s\n\t\t\t\tLocation: %s",
		t.Status, t.Location)
}

func NewTemperatureProbeLocationAndStatus(data byte) TemperatureProbeLocationAndStatus {
	return TemperatureProbeLocationAndStatus{
		Status:   TemperatureProbeStatus(data & 0xE0),
		Location: TemperatureProbeLocation(data & 0x1F),
	}
}

type TemperatureProbe struct {
	InfoCommon
	Description       string
	LocationAndStatus TemperatureProbeLocationAndStatus
	MaximumValue      uint16
	MinimumValue      uint16
	Resolution        uint16
	Tolerance         uint16
	Accuracy          uint16
	OEMdefined        uint32
	NominalValue      uint16
}

func (t TemperatureProbe) String() string {
	return fmt.Sprintf("Temperature Probe:\n\t\t"+
		"Description: %s\n\t\t"+
		"Location And Status: %s\n\t\t"+
		"Maximum Value: %d\n\t\t"+
		"Minimum Value: %d\n\t\t"+
		"Resolution: %d\n\t\t"+
		"Tolerance: %d\n\t\t"+
		"Accuracy: %d\n\t\t"+
		"OE Mdefined: %d\n\t\t"+
		"Nominal Value: %d\n",
		t.Description,
		t.LocationAndStatus,
		t.MaximumValue,
		t.MinimumValue,
		t.Resolution,
		t.Tolerance,
		t.Accuracy,
		t.OEMdefined,
		t.NominalValue)
}

func (h DMIHeader) TemperatureProbe() *TemperatureProbe {
	data := h.data
	return &TemperatureProbe{
		Description:       h.FieldString(int(data[0x04])),
		LocationAndStatus: NewTemperatureProbeLocationAndStatus(data[0x05]),
		MaximumValue:      U16(data[0x06:0x08]),
		MinimumValue:      U16(data[0x08:0x0A]),
		Resolution:        U16(data[0x0A:0x0C]),
		Tolerance:         U16(data[0x0C:0x0E]),
		Accuracy:          U16(data[0x0E:0x10]),
		OEMdefined:        U32(data[0x10:0x14]),
		NominalValue:      U16(data[0x14:0x16]),
	}
}

type ElectricalCurrentProbeStatus byte

const (
	ElectricalCurrentProbeStatusOther ElectricalCurrentProbeStatus = 0x20 + iota
	ElectricalCurrentProbeStatusUnknown
	ElectricalCurrentProbeStatusOK
	ElectricalCurrentProbeStatusNon_critical
	ElectricalCurrentProbeStatusCritical
	ElectricalCurrentProbeStatusNon_recoverable
)

func (e ElectricalCurrentProbeStatus) String() string {
	status := [...]string{
		"Other",
		"Unknown",
		"OK",
		"Non-critical",
		"Critical",
		"Non-recoverable",
	}
	return status[e-0x20]
}

type ElectricalCurrentProbeLocation byte

const (
	ElectricalCurrentProbeLocationOther ElectricalCurrentProbeLocation = 1 + iota
	ElectricalCurrentProbeLocationUnknown
	ElectricalCurrentProbeLocationProcessor
	ElectricalCurrentProbeLocationDisk
	ElectricalCurrentProbeLocationPeripheralBay
	ElectricalCurrentProbeLocationSystemManagementModule
	ElectricalCurrentProbeLocationMotherboard
	ElectricalCurrentProbeLocationMemoryModule
	ElectricalCurrentProbeLocationProcessorModule
	ElectricalCurrentProbeLocationPowerUnit
	ElectricalCurrentProbeLocationAdd_inCard
)

func (e ElectricalCurrentProbeLocation) String() string {
	locations := [...]string{
		"Other",
		"Unknown",
		"Processor",
		"Disk",
		"Peripheral Bay",
		"System Management Module",
		"Motherboard",
		"Memory Module",
		"Processor Module",
		"Power Unit",
		"Add-in Card",
	}
	return locations[e-1]
}

type ElectricalCurrentProbeLocationAndStatus struct {
	Status   ElectricalCurrentProbeStatus
	Location ElectricalCurrentProbeLocation
}

func (e ElectricalCurrentProbeLocationAndStatus) String() string {
	return fmt.Sprintf("\n\t\t\t\tStatus: %s\n\t\t\t\tLocation: %s",
		e.Status, e.Location)

}

func NewElectricalCurrentProbeLocationAndStatus(data byte) ElectricalCurrentProbeLocationAndStatus {
	return ElectricalCurrentProbeLocationAndStatus{
		Status:   ElectricalCurrentProbeStatus(data & 0xE0),
		Location: ElectricalCurrentProbeLocation(data & 0x1F),
	}
}

type ElectricalCurrentProbe struct {
	InfoCommon
	Description       string
	LocationAndStatus ElectricalCurrentProbeLocationAndStatus
	MaximumValue      uint16
	MinimumValue      uint16
	Resolution        uint16
	Tolerance         uint16
	Accuracy          uint16
	OEMdefined        uint32
	NomimalValue      uint16
}

func (e ElectricalCurrentProbe) String() string {
	return fmt.Sprintf("Electrical Current Probe:\n\t\t"+
		"Description: %s\n\t\t"+
		"Location And Status: %s\n\t\t"+
		"Maximum Value: %d\n\t\t"+
		"Minimum Value: %d\n\t\t"+
		"Resolution: %d\n\t\t"+
		"Tolerance: %d\n\t\t"+
		"Accuracy: %d\n\t\t"+
		"OE Mdefined: %d\n\t\t"+
		"Nomimal Value: %d\n",
		e.Description,
		e.LocationAndStatus,
		e.MaximumValue,
		e.MinimumValue,
		e.Resolution,
		e.Tolerance,
		e.Accuracy,
		e.OEMdefined,
		e.NomimalValue)
}

func (h DMIHeader) ElectricalCurrentProbe() *ElectricalCurrentProbe {
	data := h.data
	return &ElectricalCurrentProbe{
		Description:       h.FieldString(int(data[0x04])),
		LocationAndStatus: NewElectricalCurrentProbeLocationAndStatus(data[0x05]),
		MaximumValue:      U16(data[0x06:0x08]),
		MinimumValue:      U16(data[0x08:0x0A]),
		Resolution:        U16(data[0x0A:0x0C]),
		Tolerance:         U16(data[0x0C:0x0E]),
		Accuracy:          U16(data[0x0E:0x10]),
		OEMdefined:        U32(data[0x10:0x14]),
		NomimalValue:      U16(data[0x14:0x16]),
	}
}

type OutOfBandRemoteAccessConnections struct {
	OutBoundEnabled bool
	InBoundEnabled  bool
}

func NewOutOfBandRemoteAccessConnections(data byte) OutOfBandRemoteAccessConnections {
	return OutOfBandRemoteAccessConnections{
		OutBoundEnabled: (data&0x02 != 0),
		InBoundEnabled:  (data&0x01 != 0),
	}
}

func (o OutOfBandRemoteAccessConnections) String() string {
	return fmt.Sprintf("\n\t\t\t\tOutbound Enabled: %t\n\t\t\t\tInbound Enabled: %t",
		o.OutBoundEnabled, o.InBoundEnabled)
}

type OutOfBandRemoteAccess struct {
	InfoCommon
	ManufacturerName string
	Connections      OutOfBandRemoteAccessConnections
}

func (o OutOfBandRemoteAccess) String() string {
	return fmt.Sprintf("Out Of Band Remote Access:\n\t\t"+
		"Manufacturer Name: %s\n\t\t"+
		"Connections: %s\n",
		o.ManufacturerName,
		o.Connections)
}

func (h DMIHeader) OutOfBandRemoteAccess() *OutOfBandRemoteAccess {
	data := h.data
	return &OutOfBandRemoteAccess{
		ManufacturerName: h.FieldString(int(data[0x04])),
		Connections:      NewOutOfBandRemoteAccessConnections(data[0x05]),
	}
}

type SystemBootInformationStatus byte

func (s SystemBootInformationStatus) String() string {
	status := [...]string{
		"No errors detected", /* 0 */
		"No bootable media",
		"Operating system failed to load",
		"Firmware-detected hardware failure",
		"Operating system-detected hardware failure",
		"User-requested boot",
		"System security violation",
		"Previously-requested image",
		"System watchdog timer expired",
	}
	if s <= 8 {
		return status[s]
	} else if s >= 128 && s <= 191 {
		return "OEM-specific"
	} else if s > 192 && s <= 255 {
		return "Product-specific"
	}
	return "Error"
}

type SystemBootInformation struct {
	InfoCommon
	BootStatus SystemBootInformationStatus
}

func (s SystemBootInformation) String() string {
	return fmt.Sprintf("System Boot Information:\n\t\t"+
		"Boot Status: %s\n",
		s.BootStatus)
}

func (h DMIHeader) SystemBootInformation() *SystemBootInformation {
	data := h.data
	return &SystemBootInformation{
		BootStatus: SystemBootInformationStatus(data[0x0A]),
	}
}

type _64BitMemoryErrorInformation struct {
	InfoCommon
	Type              MemoryErrorInformationType
	Granularity       MemoryErrorInformationGranularity
	Operation         MemoryErrorInformationOperation
	VendorSyndrome    uint32
	ArrayErrorAddress uint32
	ErrorAddress      uint32
	Reslution         uint32
}

func (m _64BitMemoryErrorInformation) String() string {
	return fmt.Sprintf("_64 Bit Memory Error Information:\n\t\t"+
		"Type: %s\n\t\t"+
		"Granularity: %s\n\t\t"+
		"Operation: %s\n\t\t"+
		"Vendor Syndrome: %d\n\t\t"+
		"Array Error Address: %d\n\t\t"+
		"Error Address: %d\n\t\t"+
		"Reslution: %d\n",
		m.Type,
		m.Granularity,
		m.Operation,
		m.VendorSyndrome,
		m.ArrayErrorAddress,
		m.ErrorAddress,
		m.Reslution)
}

func (h DMIHeader) _64BitMemoryErrorInformation() *_64BitMemoryErrorInformation {
	data := h.data
	return &_64BitMemoryErrorInformation{
		Type:              MemoryErrorInformationType(data[0x04]),
		Granularity:       MemoryErrorInformationGranularity(data[0x05]),
		Operation:         MemoryErrorInformationOperation(data[0x06]),
		VendorSyndrome:    U32(data[0x07:0x0B]),
		ArrayErrorAddress: U32(data[0x0B:0x0F]),
		ErrorAddress:      U32(data[0x0F:0x13]),
		Reslution:         U32(data[0x13:0x17]),
	}
}

type ManagementDeviceType byte

const (
	ManagementDeviceTypeOther ManagementDeviceType = 1 + iota
	ManagementDeviceTypeUnknown
	ManagementDeviceTypeNationalSemiconductorLM75
	ManagementDeviceTypeNationalSemiconductorLM78
	ManagementDeviceTypeNationalSemiconductorLM79
	ManagementDeviceTypeNationalSemiconductorLM80
	ManagementDeviceTypeNationalSemiconductorLM81
	ManagementDeviceTypeAnalogDevicesADM9240
	ManagementDeviceTypeDallasSemiconductorDS1780
	ManagementDeviceTypeMaxim1617
	ManagementDeviceTypeGenesysGL518SM
	ManagementDeviceTypeWinbondW83781D
	ManagementDeviceTypeHoltekHT82H791
)

func (m ManagementDeviceType) String() string {
	types := [...]string{
		"Other",
		"Unknown",
		"National Semiconductor LM75",
		"National Semiconductor LM78",
		"National Semiconductor LM79",
		"National Semiconductor LM80",
		"National Semiconductor LM81",
		"Analog Devices ADM9240",
		"Dallas Semiconductor DS1780",
		"Maxim 1617",
		"Genesys GL518SM",
		"Winbond W83781D",
		"Holtek HT82H791",
	}
	return types[m-1]
}

type ManagementDeviceAddressType byte

const (
	ManagementDeviceAddressTypeOther ManagementDeviceAddressType = 1 + iota
	ManagementDeviceAddressTypeUnknown
	ManagementDeviceAddressTypeIOPort
	ManagementDeviceAddressTypeMemory
	ManagementDeviceAddressTypeSMBus
)

func (m ManagementDeviceAddressType) String() string {
	types := [...]string{
		"Other",
		"Unknown",
		"I/O Port",
		"Memory",
		"SM Bus",
	}
	return types[m-1]
}

type ManagementDevice struct {
	InfoCommon
	Description string
	Type        ManagementDeviceType
	Address     uint32
	AddressType ManagementDeviceAddressType
}

func (m ManagementDevice) String() string {
	return fmt.Sprintf("Management Device:\n\t\t"+
		"Description: %s\n\t\t"+
		"Type: %s\n\t\t"+
		"Address: %d\n\t\t"+
		"Address Type: %s\n",
		m.Description,
		m.Type,
		m.Address,
		m.AddressType)
}

func (h DMIHeader) ManagementDevice() *ManagementDevice {
	data := h.data
	return &ManagementDevice{
		Description: h.FieldString(int(data[0x04])),
		Type:        ManagementDeviceType(data[0x05]),
		Address:     U32(data[0x06:0x0A]),
		AddressType: ManagementDeviceAddressType(data[0x0A]),
	}
}

type ManagementDeviceComponent struct {
	InfoCommon
	Description            string
	ManagementDeviceHandle uint16
	ComponentHandle        uint16
	ThresholdHandle        uint16
}

func (m ManagementDeviceComponent) String() string {
	return fmt.Sprintf("Management Device Component:\n\t\t"+
		"Description: %s\n\t\t"+
		"Management Device Handle: %d\n\t\t"+
		"Component Handle: %d\n\t\t"+
		"Threshold Handle: %d\n",
		m.Description,
		m.ManagementDeviceHandle,
		m.ComponentHandle,
		m.ThresholdHandle)
}

func (h DMIHeader) ManagementDeviceComponent() *ManagementDeviceComponent {
	data := h.data
	return &ManagementDeviceComponent{
		Description:            h.FieldString(int(data[0x04])),
		ManagementDeviceHandle: U16(data[0x05:0x07]),
		ComponentHandle:        U16(data[0x07:0x09]),
		ThresholdHandle:        U16(data[0x09:0x0B]),
	}
}

type ManagementDeviceThresholdData struct {
	InfoCommon
	LowerThresholdNonCritical    uint16
	UpperThresholdNonCritical    uint16
	LowerThresholdCritical       uint16
	UpperThresholdCritical       uint16
	LowerThresholdNonRecoverable uint16
	UpperThresholdNonRecoverable uint16
}

func (m ManagementDeviceThresholdData) String() string {
	return fmt.Sprintf("Management Device Threshold Data:\n\t\t"+
		"Lower Threshold Non Critical: %d\n\t\t"+
		"Upper Threshold Non Critical: %d\n\t\t"+
		"Lower Threshold Critical: %d\n\t\t"+
		"Upper Threshold Critical: %d\n\t\t"+
		"Lower Threshold Non Recoverable: %d\n\t\t"+
		"Upper Threshold Non Recoverable: %d\n",
		m.LowerThresholdNonCritical,
		m.UpperThresholdNonCritical,
		m.LowerThresholdCritical,
		m.UpperThresholdCritical,
		m.LowerThresholdNonRecoverable,
		m.UpperThresholdNonRecoverable)
}

func (h DMIHeader) ManagementDeviceThresholdData() *ManagementDeviceThresholdData {
	data := h.data
	return &ManagementDeviceThresholdData{
		LowerThresholdNonCritical:    U16(data[0x04:0x06]),
		UpperThresholdNonCritical:    U16(data[0x06:0x08]),
		LowerThresholdCritical:       U16(data[0x08:0x0A]),
		UpperThresholdCritical:       U16(data[0x0A:0x0C]),
		LowerThresholdNonRecoverable: U16(data[0x0C:0x0E]),
		UpperThresholdNonRecoverable: U16(data[0x0E:0x10]),
	}
}

type MemoryChannelType byte

const (
	MemoryChannelTypeOther MemoryChannelType = 1 + iota
	MemoryChannelTypeUnknown
	MemoryChannelTypeRamBus
	MemoryChannelTypeSyncLink
)

func (m MemoryChannelType) String() string {
	types := [...]string{
		"Other",
		"Unknown",
		"RamBus",
		"SyncLink",
	}
	return types[m-1]
}

type MemoryDeviceLoadHandle struct {
	Load   byte
	Handle uint16
}

type MemoryDeviceLoadHandles []MemoryDeviceLoadHandle

func newMemoryDeviceLoadHandles(data []byte, count byte, length byte) MemoryDeviceLoadHandles {
	md := make([]MemoryDeviceLoadHandle, 0)
	if length < 0x07+count {
		return md
	}
	for i := byte(1); i <= count; i++ {
		var mem MemoryDeviceLoadHandle
		offset := 3 * (i - 1)
		mem.Load = data[0x07+offset]
		mem.Handle = U16(data[0x08+offset : 0x0A+offset])
		md = append(md, mem)
	}
	return md
}

func (m MemoryDeviceLoadHandles) String() string {
	var s string
	for _, md := range m {
		s += fmt.Sprintf("\n\t\tDevice: %d\tHandle %d", md.Load, md.Handle)
	}
	return s
}

type MemoryChannel struct {
	InfoCommon
	ChannelType        MemoryChannelType
	MaximumChannelLoad byte
	MemoryDeviceCount  byte
	LoadHandle         MemoryDeviceLoadHandles
}

func (m MemoryChannel) String() string {
	return fmt.Sprintf("Memory Channel:\n\t\t"+
		"Channel Type: %s\n\t\t"+
		"Maximum Channel Load: %d\n\t\t"+
		"%s",
		m.ChannelType,
		m.MaximumChannelLoad,
		m.LoadHandle)
}

func (h DMIHeader) MemoryChannel() *MemoryChannel {
	data := h.data
	mc := &MemoryChannel{
		ChannelType:        MemoryChannelType(data[0x04]),
		MaximumChannelLoad: data[0x05],
		MemoryDeviceCount:  data[0x06],
	}
	mc.LoadHandle = newMemoryDeviceLoadHandles(data, data[0x06], h.Length)
	return mc
}

type IPMIDeviceInformationInterfaceType byte

const (
	IPMIDeviceInformationInterfaceTypeUnknown IPMIDeviceInformationInterfaceType = 1 + iota
	IPMIDeviceInformationInterfaceTypeKCSKeyboardControllerStyle
	IPMIDeviceInformationInterfaceTypeSMICServerManagementInterfaceChip
	IPMIDeviceInformationInterfaceTypeBTBlockTransfer
	IPMIDeviceInformationInterfaceTypeReservedforfutureassignmentbythisspecification
)

func (i IPMIDeviceInformationInterfaceType) String() string {
	types := [...]string{
		"Unknown",
		"KCS: Keyboard Controller Style",
		"SMIC: Server Management Interface Chip",
		"BT: Block Transfer",
		"Reserved for future assignment by this specification",
	}
	if i <= 3 {
		return types[i]
	}
	return types[4]
}

type IPMIDeviceInformationInfo byte

const (
	IPMIDeviceInformationInfoNotSpecified IPMIDeviceInformationInfo = iota
	IPMIDeviceInformationInfoSpecified
)

func (i IPMIDeviceInformationInfo) String() string {
	info := [...]string{
		"not specified",
		"specified",
	}
	return info[i]
}

type IPMIDeviceInformationPolarity byte

const (
	IPMIDeviceInformationPolarityActiveLow IPMIDeviceInformationPolarity = iota
	IPMIDeviceInformationPolarityActiveHigh
)

func (i IPMIDeviceInformationPolarity) String() string {
	polarities := [...]string{
		"active low",
		"active high",
	}
	return polarities[i]
}

type IPMIDeviceInformationTriggerMode byte

const (
	IPMIDeviceInformationTriggerModeEdge IPMIDeviceInformationTriggerMode = iota
	IPMIDeviceInformationTriggerModeLevel
)

func (i IPMIDeviceInformationTriggerMode) String() string {
	modes := [...]string{
		"edge",
		"level",
	}
	return modes[i]
}

type IPMIDeviceInformationInterruptInfo struct {
	Info        IPMIDeviceInformationInfo
	Polarity    IPMIDeviceInformationPolarity
	TriggerMode IPMIDeviceInformationTriggerMode
}

type IPMIDeviceInformationRegisterSpacing byte

const (
	IPMIDeviceInformationRegisterSpacingSuccessiveByteBoundaries IPMIDeviceInformationRegisterSpacing = iota
	IPMIDeviceInformationRegisterSpacing32BitBoundaries
	IPMIDeviceInformationRegisterSpacing16ByteBoundaries
	IPMIDeviceInformationRegisterSpacingReserved
)

func (i IPMIDeviceInformationRegisterSpacing) String() string {
	space := [...]string{
		"Interface registers are on successive byte boundaries",
		"Interface registers are on 32-bit boundaries",
		"Interface registers are on 16-byte boundaries",
		"Reserved",
	}
	return space[i]
}

type IPMIDeviceInformationLSbit byte

type IPMIDeviceInformationBaseModifier struct {
	RegisterSpacing IPMIDeviceInformationRegisterSpacing
	LSbit           IPMIDeviceInformationLSbit
}

type IPMIDeviceInformationAddressModiferInterruptInfo struct {
	BaseAddressModifier IPMIDeviceInformationBaseModifier
	InterruptInfo       IPMIDeviceInformationInterruptInfo
}

func (i IPMIDeviceInformationAddressModiferInterruptInfo) String() string {
	return fmt.Sprintf("Base Address Modifier:"+
		"\n\t\t\t\tRegister spacing: %s"+
		"\n\t\t\t\tLs-bit for addresses: %d"+
		"\n\t\tInterrupt Info:"+
		"\n\t\t\t\tInfo: %s"+
		"\n\t\t\t\tPolarity: %s"+
		"\n\t\t\t\tTrigger Mode: %s",
		i.BaseAddressModifier.RegisterSpacing,
		i.BaseAddressModifier.LSbit,
		i.InterruptInfo.Info,
		i.InterruptInfo.Polarity,
		i.InterruptInfo.TriggerMode)
}

func newIPMIDeviceInformationAddressModiferInterruptInfo(base byte) IPMIDeviceInformationAddressModiferInterruptInfo {
	var ipmi IPMIDeviceInformationAddressModiferInterruptInfo
	ipmi.BaseAddressModifier.RegisterSpacing = IPMIDeviceInformationRegisterSpacing((base & 0xC0) >> 6)
	ipmi.BaseAddressModifier.LSbit = IPMIDeviceInformationLSbit((base & 0x10) >> 4)
	ipmi.InterruptInfo.Info = IPMIDeviceInformationInfo((base & 0x08) >> 3)
	ipmi.InterruptInfo.Polarity = IPMIDeviceInformationPolarity((base & 0x02) >> 1)
	ipmi.InterruptInfo.TriggerMode = IPMIDeviceInformationTriggerMode(base & 0x01)
	return ipmi
}

type IPMIDeviceInformation struct {
	InfoCommon
	InterfaceType                  IPMIDeviceInformationInterfaceType
	Revision                       byte
	I2CSlaveAddress                byte
	NVStorageAddress               byte
	BaseAddress                    uint64
	BaseAddressModiferInterrutInfo IPMIDeviceInformationAddressModiferInterruptInfo
	InterruptNumbe                 byte
}

func (i IPMIDeviceInformation) String() string {
	return fmt.Sprintf("IPMI Device Information:\n\t\t"+
		"Interface Type: %s\n\t\t"+
		"Revision: %d\n\t\t"+
		"I2C Slave Address: %d\n\t\t"+
		"NV Storage Address: %d\n\t\t"+
		"Base Address: %d\n\t\t"+
		"Base Address Modifer Interrut Info: %s\n\t\t"+
		"Interrupt Numbe: %d\n",
		i.InterfaceType,
		i.Revision,
		i.I2CSlaveAddress,
		i.NVStorageAddress,
		i.BaseAddress,
		i.BaseAddressModiferInterrutInfo,
		i.InterruptNumbe)
}

func (h DMIHeader) IPMIDeviceInformation() *IPMIDeviceInformation {
	data := h.data
	return &IPMIDeviceInformation{
		InterfaceType:                  IPMIDeviceInformationInterfaceType(data[0x04]),
		Revision:                       data[0x05],
		I2CSlaveAddress:                data[0x06],
		NVStorageAddress:               data[0x07],
		BaseAddress:                    U64(data[0x08:0x10]),
		BaseAddressModiferInterrutInfo: newIPMIDeviceInformationAddressModiferInterruptInfo(data[0x10]),
		InterruptNumbe:                 data[0x11],
	}
}

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
	return fmt.Sprintf("System Power Supply Characteristics:"+
		"\n\t\t\t\tDMTF Power Supply Type: %s"+
		"\n\t\t\t\tStatus: %s"+
		"\n\t\t\t\tDMTF Input Voltage Switching: %s"+
		"\n\t\t\t\tIs Unplugged From Wall: %t"+
		"\n\t\t\t\tIs Present: %t"+
		"\n\t\t\t\tIs Hot Repleaceable: %t\n",
		s.DMTFPowerSupplyType,
		s.Status,
		s.DMTFInputVoltageSwitching,
		s.IsUnpluggedFromWall,
		s.IsPresent,
		s.IsHotRepleaceable)
}

type SystemPowerSupply struct {
	InfoCommon
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
	return fmt.Sprintf("System Power Supply:\n\t\t"+
		"Power Unit Group: %d\n\t\t"+
		"Location: %s\n\t\t"+
		"Device Name: %s\n\t\t"+
		"Manufacturer: %s\n\t\t"+
		"Serial Number: %s\n\t\t"+
		"Asset Tag Number: %s\n\t\t"+
		"Model Part Number: %s\n\t\t"+
		"Revision Level: %s\n\t\t"+
		"Max Power Capacity: %d\n\t\t"+
		"Power Supply Characteristics: %s\n\t\t"+
		"Input Voltage Probe Handle: %d\n\t\t"+
		"Cooling Device Handle: %d\n\t\t"+
		"Input Current Probe Handle: %d\n",
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

func (h DMIHeader) SystemPowerSupply() *SystemPowerSupply {
	data := h.data
	return &SystemPowerSupply{
		PowerUnitGroup:             data[0x04],
		Location:                   h.FieldString(int(data[0x05])),
		DeviceName:                 h.FieldString(int(data[0x06])),
		Manufacturer:               h.FieldString(int(data[0x07])),
		SerialNumber:               h.FieldString(int(data[0x08])),
		AssetTagNumber:             h.FieldString(int(data[0x09])),
		ModelPartNumber:            h.FieldString(int(data[0x0A])),
		RevisionLevel:              h.FieldString(int(data[0x0B])),
		MaxPowerCapacity:           U16(data[0x0C:0x0E]),
		PowerSupplyCharacteristics: newSystemPowerSupplyCharacteristics(U16(data[0x0E : 0x0E+2])),
		InputVoltageProbeHandle:    U16(data[0x0F:0x11]),
		CoolingDeviceHandle:        U16(data[0x11:0x13]),
		InputCurrentProbeHandle:    U16(data[0x13:0x15]),
	}
}

type AdditionalInformationEntries struct {
	Length           byte
	ReferencedHandle uint16
	ReferencedOffset byte
	String           string
	Value            []byte
}

type AdditionalInformationEntriess []AdditionalInformationEntries

func (a AdditionalInformationEntriess) String() string {
	var str string
	for _, s := range a {
		str += fmt.Sprintf("\n\t\t\t\tReferenced Handle: %d"+
			"\n\t\t\t\tReferenced Offset: %d"+
			"\n\t\t\t\tString: %s"+
			"\n\t\t\t\tValue: %v",
			s.ReferencedHandle,
			s.ReferencedOffset,
			s.String,
			s.Value)
	}
	return str
}

type AdditionalInformation struct {
	InfoCommon
	NumberOfEntries byte
	Entries         []AdditionalInformationEntries
}

func (a AdditionalInformation) String() string {
	return fmt.Sprintf("Additional Information:\n\t\t"+
		"Number Of Entries: %d\n\t\t"+
		"Entries: %s\n",
		a.NumberOfEntries,
		AdditionalInformationEntriess(a.Entries))
}

func (h DMIHeader) AdditionalInformation() *AdditionalInformation {
	data := h.data
	ai := new(AdditionalInformation)
	ai.NumberOfEntries = data[0x04]
	en := make([]AdditionalInformationEntries, 0)
	d := data[0x05:]
	for i := byte(0); i < ai.NumberOfEntries; i++ {
		var e AdditionalInformationEntries
		e.Length = d[0x0]
		e.ReferencedHandle = U16(d[0x01:0x03])
		e.ReferencedOffset = d[0x03]
		e.String = h.FieldString(int(d[0x04]))
		e.Value = data[0x05:e.Length]
		en = append(en, e)
		d = data[0x05+e.Length:]
	}
	ai.Entries = en
	return ai
}

type OnBoardDevicesExtendedInformationType byte

const (
	OnBoardDevicesExtendedInformationTypeOther OnBoardDevicesExtendedInformationType = 1 + iota
	OnBoardDevicesExtendedInformationTypeUnknown
	OnBoardDevicesExtendedInformationTypeVideo
	OnBoardDevicesExtendedInformationTypeSCSIController
	OnBoardDevicesExtendedInformationTypeEthernet
	OnBoardDevicesExtendedInformationTypeTokenRing
	OnBoardDevicesExtendedInformationTypeSound
	OnBoardDevicesExtendedInformationTypePATAController
	OnBoardDevicesExtendedInformationTypeSATAController
	OnBoardDevicesExtendedInformationTypeSASController
)

func (o OnBoardDevicesExtendedInformationType) String() string {
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
	return types[o-1]
}

type OnBoardDevicesExtendedInformation struct {
	InfoCommon
	ReferenceDesignation string
	DeviceType           OnBoardDevicesExtendedInformationType
	DeviceTypeInstance   byte
	SegmentGroupNumber   uint16
	BusNumber            byte
	DeviceFunctionNumber byte
}

func (o OnBoardDevicesExtendedInformation) SlotSegment() string {
	if o.SegmentGroupNumber == 0xFFFF || o.BusNumber == 0xFF || o.DeviceFunctionNumber == 0xFF {
		return "Not of types PCI/AGP/PCI-X/PCI-Express"
	}
	return fmt.Sprintf("Bus Address: %04x:%02x:%02x.%x",
		o.SegmentGroupNumber,
		o.BusNumber,
		o.DeviceFunctionNumber>>3,
		o.DeviceFunctionNumber&0x7)
}

func (o OnBoardDevicesExtendedInformation) String() string {
	return fmt.Sprintf("On Board Devices Extended Information:\n\t\t"+
		"Reference Designation: %s\n\t\t"+
		"Device Type: %s\n\t\t"+
		"Device Type Instance: %d\n\t\t"+
		"%s\n",
		o.ReferenceDesignation,
		o.DeviceType,
		o.DeviceTypeInstance,
		o.SlotSegment())
}

func (h DMIHeader) OnBoardDevicesExtendedInformation() *OnBoardDevicesExtendedInformation {
	data := h.data
	return &OnBoardDevicesExtendedInformation{
		ReferenceDesignation: h.FieldString(int(data[0x04])),
		DeviceType:           OnBoardDevicesExtendedInformationType(data[0x05]),
		DeviceTypeInstance:   data[0x06],
		SegmentGroupNumber:   U16(data[0x07:0x09]),
		BusNumber:            data[0x09],
		DeviceFunctionNumber: data[0x0A],
	}
}

type ManagementControllerHostInterfaceType byte

const (
	ManagementControllerHostInterfaceTypeKCSKeyboardControllerStyle ManagementControllerHostInterfaceType = 0x02 + iota
	ManagementControllerHostInterfaceType8250UARTRegisterCompatible
	ManagementControllerHostInterfaceType16450UARTRegisterCompatible
	ManagementControllerHostInterfaceType16550_16550AUARTRegisterCompatible
	ManagementControllerHostInterfaceType16650_16650AUARTRegisterCompatible
	ManagementControllerHostInterfaceType16750_16750AUARTRegisterCompatible
	ManagementControllerHostInterfaceType16850_16850AUARTRegisterCompatible
)

func (m ManagementControllerHostInterfaceType) String() string {
	types := [...]string{
		"KCS: Keyboard Controller Style",
		"8250 UART Register Compatible",
		"16450 UART Register Compatible",
		"16550/16550A UART Register Compatible",
		"16650/16650A UART Register Compatible",
		"16750/16750A UART Register Compatible",
		"16850/16850A UART Register Compatible",
	}
	if m >= 0x02 && m <= 0x08 {
		return types[m-0x02]
	}
	if m == 0xf0 {
		return "OEM"
	}
	return "<OUT OF SPEC>"
}

type ManagementControllerHostInterfaceData []byte

type ManagementControllerHostInterface struct {
	InfoCommon
	Type ManagementControllerHostInterfaceType
	Data ManagementControllerHostInterfaceData
}

func (m ManagementControllerHostInterface) String() string {
	return fmt.Sprintf("Management Controller Host Interface:\n\t\t"+
		"Type: %s\n\t\t"+
		"MC Host Interface Data: %s\n",
		m.Type,
		m.MCHostInterfaceData)
}

func (m ManagementControllerHostInterface) MCHostInterfaceData() string {
	if m.Type == 0xF0 {
		return fmt.Sprintf("Vendor ID:0x%02X%02X%02X%02X",
			m.Data[0x01], m.Data[0x02], m.Data[0x03], m.Data[0x04])
	}
	return ""
}

func (h DMIHeader) ManagementControllerHostInterface() *ManagementControllerHostInterface {
	data := h.data
	mc := &ManagementControllerHostInterface{
		Type: ManagementControllerHostInterfaceType(data[0x04]),
	}
	if mc.Type == 0xF0 {
		mc.Data = data[0x05 : 0x05+4]
	}
	return mc
}

type Inactive struct {
	InfoCommon
}

func (i Inactive) String() string {
	return "Inactive"
}

func (h DMIHeader) Inactive() *Inactive {
	return &Inactive{}
}

type EndOfTable struct {
	InfoCommon
}

func (e EndOfTable) String() string {
	return "End-of-Table"
}

func (h DMIHeader) EndOfTable() *EndOfTable {
	return &EndOfTable{}
}

func bcd(data []byte) int64 {
	var b int64
	l := len(data)
	if l > 8 {
		panic("bcd: Out of range")
	}
	// Number of 4-bits
	nb := int64(l * 2)
	for i := int64(0); i < nb; i++ {
		var shift uint64
		if i%2 == 0 {
			shift = 0
		} else {
			shift = 4
		}
		b += int64((data[i/2]>>shift)&0x0F) * int64(math.Pow10(int(i)))
	}
	return b
}

func U16(data []byte) uint16 {
	var u16 uint16
	binary.Read(bytes.NewBuffer(data[0:2]), binary.LittleEndian, &u16)
	return u16
}

func U32(data []byte) uint32 {
	var u32 uint32
	binary.Read(bytes.NewBuffer(data[0:4]), binary.LittleEndian, &u32)
	return u32
}

func U64(data []byte) uint64 {
	var u64 uint64
	binary.Read(bytes.NewBuffer(data[0:8]), binary.LittleEndian, &u64)
	return u64
}

func NewDMIHeader(data []byte) *DMIHeader {
	if len(data) < 0x04 {
		return nil
	}
	return &DMIHeader{
		InfoCommon: InfoCommon{
			Type:   SMBIOSStructureType(data[0x00]),
			Length: data[1],
			Handle: SMBIOSStructureHandle(U16(data[0x02:0x04])),
		},
		data: data}
}

func newSMBIOS_EPS() (eps *SMBIOS_EPS, err error) {
	eps = new(SMBIOS_EPS)

	mem, err := getMem(0xF0000, 0x10000)
	if err != nil {
		return
	}
	data := anchor(mem)
	eps.Anchor = data[:0x04]
	eps.Checksum = data[0x04]
	eps.Length = data[0x05]
	eps.MajorVersion = data[0x06]
	eps.MinorVersion = data[0x07]
	eps.MaxSize = U16(data[0x08:0x0A])
	eps.FormattedArea = data[0x0B:0x0F]
	eps.InterAnchor = data[0x10:0x15]
	eps.InterChecksum = data[0x15]
	eps.TableLength = U16(data[0x16:0x18])
	eps.TableAddress = U32(data[0x18:0x1C])
	eps.NumberOfSM = U16(data[0x1C:0x1E])
	eps.BCDRevision = data[0x1E]
	return
}

func (e SMBIOS_EPS) StructureTableMem() ([]byte, error) {
	return getMem(e.TableAddress, uint32(e.TableLength))
}

func (h DMIHeader) Next() *DMIHeader {
	de := []byte{0, 0}
	next := h.data[h.Length:]
	index := bytes.Index(next, de)
	if index == -1 {
		return nil
	}
	return NewDMIHeader(next[index+2:])
}

func (h DMIHeader) Decode() interface{} {
	switch h.Type {
	case SMBIOSStructureTypeBIOS:
		return h.BIOSInformation()
	case SMBIOSStructureTypeSystem:
		return h.SystemInformation()
	case SMBIOSStructureTypeBaseBoard:
		return h.BaseboardInformation()
	case SMBIOSStructureTypeChassis:
		return h.ChassisInformation()
	case SMBIOSStructureTypeProcessor:
		return h.ProcessorInformation()
	case SMBIOSStructureTypeCache:
		return h.CacheInformation()
	case SMBIOSStructureTypePortConnector:
		return h.PortInformation()
	case SMBIOSStructureTypeSystemSlots:
		return h.SystemSlot()
	case SMBIOSStructureTypeOnBoardDevices:
		return h.OnBoardDeviceInformation()
	case SMBIOSStructureTypeBIOSLanguage:
		return h.BIOSLanguageInformation()
	case SMBIOSStructureTypeSystemConfigurationOptions:
		return h.SystemConfigurationOptions()
	case SMBIOSStructureTypeOEMStrings:
		return h.OEMStrings()
	case SMBIOSStructureTypeGroupAssociations:
		return h.GroupAssociations()
	case SMBIOSStructureTypePhysicalMemoryArray:
		return h.PhysicalMemoryArray()
	case SMBIOSStructureTypeMemoryDevice:
		return h.MemoryDevice()
	case SMBIOSStructureType32_bitMemoryError:
		return h._32BitMemoryErrorInformation()
	case SMBIOSStructureTypeBuilt_inPointingDevice:
		return h.BuiltinPointingDevice()
	case SMBIOSStructureTypePortableBattery:
		return h.PortableBattery()
	case SMBIOSStructureTypeSystemReset:
		return h.SystemReset()
	case SMBIOSStructureTypeHardwareSecurity:
		return h.HardwareSecurity()
	case SMBIOSStructureTypeSystemPowerControls:
		return h.SystemPowerControls()
	case SMBIOSStructureTypeVoltageProbe:
		return h.VoltageProbe()
	case SMBIOSStructureTypeCoolingDevice:
		return h.CoolingDevice()
	case SMBIOSStructureTypeTemperatureProbe:
		return h.TemperatureProbe()
	case SMBIOSStructureTypeElectricalCurrentProbe:
		return h.ElectricalCurrentProbe()
	case SMBIOSStructureTypeOut_of_bandRemoteAccess:
		return h.OutOfBandRemoteAccess()
	case SMBIOSStructureTypeSystemBoot:
		return h.SystemBootInformation()
	case SMBIOSStructureType64_bitMemoryError:
		return h._64BitMemoryErrorInformation()
	case SMBIOSStructureTypeManagementDevice:
		return h.ManagementDevice()
	case SMBIOSStructureTypeManagementDeviceComponent:
		return h.ManagementDeviceComponent()
	case SMBIOSStructureTypeMemoryChannel:
		return h.MemoryChannel()
	case SMBIOSStructureTypeIPMIDevice:
		return h.IPMIDeviceInformation()
	case SMBIOSStructureTypePowerSupply:
		return h.SystemPowerSupply()
	case SMBIOSStructureTypeAdditionalInformation:
		return h.AdditionalInformation()
	case SMBIOSStructureTypeOnBoardDevicesExtendedInformation:
		return h.OnBoardDevicesExtendedInformation()
	case SMBIOSStructureTypeManagementControllerHostInterface:
		return h.ManagementControllerHostInterface()
	case SMBIOSStructureTypeInactive:
		return h.Inactive()
	case SMBIOSStructureTypeEndOfTable:
		return h.EndOfTable()
	}
	return nil
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

func (h DMIHeader) BIOSInformation() bIOSInformation {
	var bi bIOSInformation
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

func (bi bIOSInformation) String() string {
	return fmt.Sprintf("BIOS Information:"+
		"\n\tVendor: %s"+
		"\n\tVersion: %s"+
		"\n\tAddress: %4"+
		"X0\n\tCharacteristics: %s"+
		"\n\tExt1:%s"+
		"\n\tExt2: %s",
		bi.Vendor,
		bi.BIOSVersion,
		bi.StartingAddressSegment,
		bi.Characteristics,
		bi.CharacteristicsExt1,
		bi.CharacteristicsExt2)
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

func (h DMIHeader) SystemInformation() systemInformation {
	var si systemInformation
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

func (h DMIHeader) BaseboardInformation() baseboardInformation {
	var bi baseboardInformation
	data := h.data
	if h.Type != 2 {
		panic("Type is not 2")
	}
	bi.Manufacturer = h.FieldString(int(data[0x04]))
	bi.Product = h.FieldString(int(data[0x05]))
	bi.Version = h.FieldString(int(data[0x06]))
	bi.SerialNumber = h.FieldString(int(data[0x07]))
	bi.AssetTag = h.FieldString(int(data[0x08]))
	bi.FeatureFlags = FeatureFlags(data[0x09])
	bi.LocationInChassis = h.FieldString(int(data[0x0A]))
	bi.BoardType = BoardType(data[0x0D])
	return bi
}

func (e SMBIOS_EPS) StructureTable() map[SMBIOSStructureType]interface{} {
	tmem, err := e.StructureTableMem()
	if err != nil {
		return nil
	}
	m := make(map[SMBIOSStructureType]interface{}, 0)
	for hd := NewDMIHeader(tmem); hd != nil; hd = hd.Next() {
		m[hd.Type] = hd.Decode()
	}
	return m
}

func init() {
	eps, err := newSMBIOS_EPS()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		panic(err)
	}
	gdmi = eps.StructureTable()
}

func SystemInformation() systemInformation {
	return gdmi[SMBIOSStructureTypeSystem].(systemInformation)
}

func BIOSInformation() bIOSInformation {
	return gdmi[SMBIOSStructureTypeBIOS].(bIOSInformation)
}

func BaseboardInformation() baseboardInformation {
	return gdmi[SMBIOSStructureTypeBaseBoard].(baseboardInformation)
}

func Chassis() *ChassisInformation {
	return gdmi[SMBIOSStructureTypeChassis].(*ChassisInformation)
}

func Processor() ProcessorInformation {
	return gdmi[SMBIOSStructureTypeProcessor].(ProcessorInformation)
}

func getMem(base uint32, length uint32) (mem []byte, err error) {
	file, err := os.Open("/dev/mem")
	if err != nil {
		return
	}
	defer file.Close()
	fd := file.Fd()
	mmoffset := base % uint32(os.Getpagesize())
	mm, err := syscall.Mmap(int(fd), int64(base-mmoffset), int(mmoffset+length), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return
	}
	mem = make([]byte, len(mm))
	copy(mem, mm)
	err = syscall.Munmap(mm)
	if err != nil {
		return
	}
	return
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
