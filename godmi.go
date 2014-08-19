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

const OUT_OF_SPEC = "<OUT OF SPEC>"

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

type infoCommon struct {
	SMType SMBIOSStructureType
	Length byte
	Handle SMBIOSStructureHandle
}

type entryPoint struct {
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

type dmiHeader struct {
	infoCommon
	data []byte
}

// BIOS Characteristics
const (
	BIOSCharacteristicsReserved0 BIOSCharacteristics = 1 << iota
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

// BIOS Characteristics Extension Bytes(Ext1)
// Byte 1
const (
	BIOSCharacteristicsExt1ACPISupported BIOSCharacteristicsExt1 = 1 << iota
	BIOSCharacteristicsExt1USBLegacySupported
	BIOSCharacteristicsExt1AGPSupported
	BIOSCharacteristicsExt1I2OBootSupported
	BIOSCharacteristicsExt1LS120SupperDiskBootSupported
	BIOSCharacteristicsExt1ATAPIZIPDriveBootSupported
	BIOSCharacteristicsExt11394BootSupported
	BIOSCharacteristicsExt1SmartBatterySupported
)

// BIOS Characteristics Extension Bytes(Ext2)
// Byte 2
const (
	BIOSCharacteristicsExt2BIOSBootSpecSupported BIOSCharacteristicsExt2 = 1 << iota
	BIOSCharacteristicsExt2FuncKeyInitiatedNetworkBootSupported
	BIOSCharacteristicsExt2EnableTargetedContentDistribution
	BIOSCharacteristicsExt2UEFISpecSupported
	BIOSCharacteristicsExt2VirtualMachine
	// Bits 5:7 Reserved for future assignment
)

func (h dmiHeader) ChassisInformation() *ChassisInformation {
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

func (h dmiHeader) ProcessorInformation() *ProcessorInformation {
	data := h.data
	return &ProcessorInformation{
		SocketDesignation: h.FieldString(int(data[0x04])),
		ProcessorType:     ProcessorType(data[0x05]),
		Family:            ProcessorFamily(data[0x06]),
		Manufacturer:      h.FieldString(int(data[0x07])),
		// TODO:
		//pi.ProcessorID
		Version:         h.FieldString(int(data[0x10])),
		Voltage:         ProcessorVoltage(data[0x11]),
		ExternalClock:   u16(data[0x12:0x14]),
		MaxSpeed:        u16(data[0x14:0x16]),
		CurrentSpeed:    u16(data[0x16:0x18]),
		Status:          ProcessorStatus(data[0x18]),
		Upgrade:         ProcessorUpgrade(data[0x19]),
		L1CacheHandle:   u16(data[0x1A:0x1C]),
		L2CacheHandle:   u16(data[0x1C:0x1E]),
		L3CacheHandle:   u16(data[0x1E:0x20]),
		SerialNumber:    h.FieldString(int(data[0x20])),
		AssetTag:        h.FieldString(int(data[0x21])),
		PartNumber:      h.FieldString(int(data[0x22])),
		CoreCount:       data[0x23],
		CoreEnabled:     data[0x24],
		ThreadCount:     data[0x25],
		Characteristics: ProcessorCharacteristics(u16(data[0x26:0x28])),
		Family2:         ProcessorFamily(data[0x28]),
	}
}

func (h dmiHeader) CacheInformation() *CacheInformation {
	data := h.data
	return &CacheInformation{
		SocketDesignation:   h.FieldString(int(data[0x04])),
		Configuration:       NewCacheConfiguration(u16(data[0x05:0x07])),
		MaximumCacheSize:    NewCacheSize(u16(data[0x07:0x09])),
		InstalledSize:       NewCacheSize(u16(data[0x09:0x0B])),
		SupportedSRAMType:   CacheSRAMType(u16(data[0x0B:0x0D])),
		CurrentSRAMType:     CacheSRAMType(u16(data[0x0D:0x0F])),
		CacheSpeed:          CacheSpeed(data[0x0F]),
		ErrorCorrectionType: CacheErrorCorrectionType(data[0x10]),
		SystemCacheType:     CacheSystemCacheType(data[0x11]),
		Associativity:       CacheAssociativity(data[0x12]),
	}
}

func (h dmiHeader) PortInformation() *PortInformation {
	data := h.data
	return &PortInformation{
		InternalReferenceDesignator: h.FieldString(int(data[0x04])),
		InternalConnectorType:       PortConnectorType(data[0x05]),
		ExternalReferenceDesignator: h.FieldString(int(data[0x06])),
		ExternalConnectorType:       PortConnectorType(data[0x07]),
		Type: PortType(data[0x08]),
	}
}

func (h dmiHeader) SystemSlot() *SystemSlot {
	data := h.data
	return &SystemSlot{
		Designation:          h.FieldString(int(data[0x04])),
		Type:                 SystemSlotType(data[0x05]),
		DataBusWidth:         SystemSlotDataBusWidth(data[0x06]),
		CurrentUsage:         SystemSlotUsage(data[0x07]),
		Length:               SystemSlotLength(data[0x08]),
		ID:                   SystemSlotID(u16(data[0x09:0x0A])),
		Characteristics1:     SystemSlotCharacteristics1(data[0x0B]),
		Characteristics2:     SystemSlotCharacteristics2(data[0x0C]),
		SegmentGroupNumber:   SystemSlotSegmengGroupNumber(u16(data[0x0D:0x0F])),
		BusNumber:            SystemSlotNumber(data[0x0F]),
		DeviceFunctionNumber: SystemSlotNumber(data[0x10]),
	}
}

func (h dmiHeader) BIOSLanguageInformation() *BIOSLanguageInformation {
	var bl BIOSLanguageInformation
	data := h.data
	cnt := data[0x04]
	for i := byte(1); i <= cnt; i++ {
		bl.InstallableLanguage = append(bl.InstallableLanguage, h.FieldString(int(data[i])))
	}
	bl.Flags = NewBIOSLanguageInformationFlag(data[0x05])
	bl.CurrentLanguage = bl.InstallableLanguage[data[0x15]]
	return &bl
}

func (h dmiHeader) OnBoardDeviceInformation() *OnBoardDeviceInformation {
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
	return &d
}

func (h dmiHeader) SystemConfigurationOptions() *SystemConfigurationOptions {
	var sc SystemConfigurationOptions
	data := h.data
	sc.Count = data[0x04]
	for i := byte(1); i <= sc.Count; i++ {
		sc.strings += fmt.Sprintf("string %d: %s\n\t\t", i, h.FieldString(int(data[0x04+i])))
	}
	return &sc
}

func (h dmiHeader) OEMStrings() *OEMStrings {
	var o OEMStrings
	data := h.data
	o.Count = data[0x04]
	for i := byte(0); i < o.Count; i++ {
		o.strings += fmt.Sprintf("strings: %d %s\n\t\t", i, h.FieldString(int(data[i])))
	}
	return &o
}

func (h dmiHeader) GroupAssociations() *GroupAssociations {
	var ga GroupAssociations
	data := h.data
	ga.GroupName = h.FieldString(int(data[0x04]))
	cnt := (h.Length - 5) / 3
	items := data[5:]
	var i byte
	for i = 0; i < cnt; i++ {
		var gai GroupAssociationsItem
		gai.Type = SMBIOSStructureType(items[i*3])
		gai.Handle = SMBIOSStructureHandle(u16(items[i*3+1:]))
		ga.Item = append(ga.Item, gai)
	}
	return &ga
}

func (h dmiHeader) PhysicalMemoryArray() *PhysicalMemoryArray {
	data := h.data
	return &PhysicalMemoryArray{
		Location:                PhysicalMemoryArrayLocation(data[0x04]),
		Use:                     PhysicalMemoryArrayUse(data[0x05]),
		ErrorCorrection:         PhysicalMemoryArrayErrorCorrection(data[0x06]),
		MaximumCapacity:         u32(data[0x07:0x0B]),
		ErrorInformationHandle:  u16(data[0x0B:0x0D]),
		NumberOfMemoryDevices:   u16(data[0x0D:0x0F]),
		ExtendedMaximumCapacity: u64(data[0x0F:]),
	}
}

func (h dmiHeader) MemoryDevice() *MemoryDevice {
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

func (h dmiHeader) _32BitMemoryErrorInformation() *_32BitMemoryErrorInformation {
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

func (h dmiHeader) BuiltinPointingDevice() *BuiltinPointingDevice {
	data := h.data
	return &BuiltinPointingDevice{
		Type:            BuiltinPointingDeviceType(data[0x04]),
		Interface:       BuiltinPointingDeviceInterface(data[0x05]),
		NumberOfButtons: data[0x06],
	}
}

func (h dmiHeader) PortableBattery() *PortableBattery {
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

func (h dmiHeader) SystemReset() *SystemReset {
	data := h.data
	return &SystemReset{
		Capabilities:  data[0x04],
		ResetCount:    u16(data[0x05:0x07]),
		ResetLimit:    u16(data[0x07:0x09]),
		TimerInterval: u16(data[0x09:0x0B]),
		Timeout:       u16(data[0x0B:0x0D]),
	}
}

func (h dmiHeader) HardwareSecurity() *HardwareSecurity {
	var hw HardwareSecurity
	data := h.data
	hw.Setting = NewHardwareSecurity(data[0x04])
	return &hw
}

func (h dmiHeader) SystemPowerControls() *SystemPowerControls {
	data := h.data
	return &SystemPowerControls{
		NextScheduledPowerOnMonth:      SystemPowerControlsMonth(bcd(data[0x04:0x05])),
		NextScheduledPowerOnDayOfMonth: SystemPowerControlsDayOfMonth(bcd(data[0x05:0x06])),
		NextScheduledPowerOnHour:       SystemPowerControlsHour(bcd(data[0x06:0x07])),
		NextScheduledPowerMinute:       SystemPowerControlsMinute(bcd(data[0x07:0x08])),
		NextScheduledPowerSecond:       SystemPowerControlsSecond(bcd(data[0x08:0x09])),
	}
}

func (h dmiHeader) VoltageProbe() *VoltageProbe {
	data := h.data
	return &VoltageProbe{
		Description:       h.FieldString(int(data[0x04])),
		LocationAndStatus: NewVoltageProbeLocationAndStatus(data[0x05]),
		MaximumValue:      u16(data[0x06:0x08]),
		MinimumValude:     u16(data[0x08:0x0A]),
		Resolution:        u16(data[0x0A:0x0C]),
		Tolerance:         u16(data[0x0C:0x0E]),
		Accuracy:          u16(data[0x0E:0x10]),
		OEMdefined:        u16(data[0x10:0x12]),
		NominalValue:      u16(data[0x12:0x14]),
	}
}

func (h dmiHeader) CoolingDevice() *CoolingDevice {
	data := h.data
	cd := &CoolingDevice{
		TemperatureProbeHandle: u16(data[0x04:0x06]),
		DeviceTypeAndStatus:    NewCoolingDeviceTypeAndStatus(data[0x06]),
		CoolingUintGroup:       data[0x07],
		OEMdefined:             u32(data[0x08:0x0C]),
	}
	if h.Length > 0x0C {
		cd.NominalSpeed = u16(data[0x0C:0x0E])
	}
	if h.Length > 0x0F {
		cd.Description = h.FieldString(int(data[0x0E]))
	}
	return cd
}

func (h dmiHeader) TemperatureProbe() *TemperatureProbe {
	data := h.data
	return &TemperatureProbe{
		Description:       h.FieldString(int(data[0x04])),
		LocationAndStatus: NewTemperatureProbeLocationAndStatus(data[0x05]),
		MaximumValue:      u16(data[0x06:0x08]),
		MinimumValue:      u16(data[0x08:0x0A]),
		Resolution:        u16(data[0x0A:0x0C]),
		Tolerance:         u16(data[0x0C:0x0E]),
		Accuracy:          u16(data[0x0E:0x10]),
		OEMdefined:        u32(data[0x10:0x14]),
		NominalValue:      u16(data[0x14:0x16]),
	}
}

func (h dmiHeader) ElectricalCurrentProbe() *ElectricalCurrentProbe {
	data := h.data
	return &ElectricalCurrentProbe{
		Description:       h.FieldString(int(data[0x04])),
		LocationAndStatus: NewElectricalCurrentProbeLocationAndStatus(data[0x05]),
		MaximumValue:      u16(data[0x06:0x08]),
		MinimumValue:      u16(data[0x08:0x0A]),
		Resolution:        u16(data[0x0A:0x0C]),
		Tolerance:         u16(data[0x0C:0x0E]),
		Accuracy:          u16(data[0x0E:0x10]),
		OEMdefined:        u32(data[0x10:0x14]),
		NomimalValue:      u16(data[0x14:0x16]),
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
	infoCommon
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

func (h dmiHeader) OutOfBandRemoteAccess() *OutOfBandRemoteAccess {
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
	infoCommon
	BootStatus SystemBootInformationStatus
}

func (s SystemBootInformation) String() string {
	return fmt.Sprintf("System Boot Information\n"+
		"\tBoot Status: %s",
		s.BootStatus)
}

func (h dmiHeader) SystemBootInformation() *SystemBootInformation {
	data := h.data
	return &SystemBootInformation{
		BootStatus: SystemBootInformationStatus(data[0x0A]),
	}
}

type _64BitMemoryErrorInformation struct {
	infoCommon
	Type              MemoryErrorInformationType
	Granularity       MemoryErrorInformationGranularity
	Operation         MemoryErrorInformationOperation
	VendorSyndrome    uint32
	ArrayErrorAddress uint32
	ErrorAddress      uint32
	Reslution         uint32
}

func (m _64BitMemoryErrorInformation) String() string {
	return fmt.Sprintf("_64 Bit Memory Error Information\n"+
		"\tType: %s\n"+
		"\tGranularity: %s\n"+
		"\tOperation: %s\n"+
		"\tVendor Syndrome: %d\n"+
		"\tArray Error Address: %d\n"+
		"\tError Address: %d\n"+
		"\tReslution: %d",
		m.Type,
		m.Granularity,
		m.Operation,
		m.VendorSyndrome,
		m.ArrayErrorAddress,
		m.ErrorAddress,
		m.Reslution)
}

func (h dmiHeader) _64BitMemoryErrorInformation() *_64BitMemoryErrorInformation {
	data := h.data
	return &_64BitMemoryErrorInformation{
		Type:              MemoryErrorInformationType(data[0x04]),
		Granularity:       MemoryErrorInformationGranularity(data[0x05]),
		Operation:         MemoryErrorInformationOperation(data[0x06]),
		VendorSyndrome:    u32(data[0x07:0x0B]),
		ArrayErrorAddress: u32(data[0x0B:0x0F]),
		ErrorAddress:      u32(data[0x0F:0x13]),
		Reslution:         u32(data[0x13:0x17]),
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
	infoCommon
	Description string
	Type        ManagementDeviceType
	Address     uint32
	AddressType ManagementDeviceAddressType
}

func (m ManagementDevice) String() string {
	return fmt.Sprintf("Management Device\n"+
		"\tDescription: %s\n"+
		"\tType: %s\n"+
		"\tAddress: %d\n"+
		"\tAddress Type: %s",
		m.Description,
		m.Type,
		m.Address,
		m.AddressType)
}

func (h dmiHeader) ManagementDevice() *ManagementDevice {
	data := h.data
	return &ManagementDevice{
		Description: h.FieldString(int(data[0x04])),
		Type:        ManagementDeviceType(data[0x05]),
		Address:     u32(data[0x06:0x0A]),
		AddressType: ManagementDeviceAddressType(data[0x0A]),
	}
}

type ManagementDeviceComponent struct {
	infoCommon
	Description            string
	ManagementDeviceHandle uint16
	ComponentHandle        uint16
	ThresholdHandle        uint16
}

func (m ManagementDeviceComponent) String() string {
	return fmt.Sprintf("Management Device Component\n"+
		"\tDescription: %s\n"+
		"\tManagement Device Handle: %d\n"+
		"\tComponent Handle: %d\n"+
		"\tThreshold Handle: %d",
		m.Description,
		m.ManagementDeviceHandle,
		m.ComponentHandle,
		m.ThresholdHandle)
}

func (h dmiHeader) ManagementDeviceComponent() *ManagementDeviceComponent {
	data := h.data
	return &ManagementDeviceComponent{
		Description:            h.FieldString(int(data[0x04])),
		ManagementDeviceHandle: u16(data[0x05:0x07]),
		ComponentHandle:        u16(data[0x07:0x09]),
		ThresholdHandle:        u16(data[0x09:0x0B]),
	}
}

type ManagementDeviceThresholdData struct {
	infoCommon
	LowerThresholdNonCritical    uint16
	UpperThresholdNonCritical    uint16
	LowerThresholdCritical       uint16
	UpperThresholdCritical       uint16
	LowerThresholdNonRecoverable uint16
	UpperThresholdNonRecoverable uint16
}

func (m ManagementDeviceThresholdData) String() string {
	return fmt.Sprintf("Management Device Threshold Data\n"+
		"\tLower Threshold Non Critical: %d\n"+
		"\tUpper Threshold Non Critical: %d\n"+
		"\tLower Threshold Critical: %d\n"+
		"\tUpper Threshold Critical: %d\n"+
		"\tLower Threshold Non Recoverable: %d\n"+
		"\tUpper Threshold Non Recoverable: %d",
		m.LowerThresholdNonCritical,
		m.UpperThresholdNonCritical,
		m.LowerThresholdCritical,
		m.UpperThresholdCritical,
		m.LowerThresholdNonRecoverable,
		m.UpperThresholdNonRecoverable)
}

func (h dmiHeader) ManagementDeviceThresholdData() *ManagementDeviceThresholdData {
	data := h.data
	return &ManagementDeviceThresholdData{
		LowerThresholdNonCritical:    u16(data[0x04:0x06]),
		UpperThresholdNonCritical:    u16(data[0x06:0x08]),
		LowerThresholdCritical:       u16(data[0x08:0x0A]),
		UpperThresholdCritical:       u16(data[0x0A:0x0C]),
		LowerThresholdNonRecoverable: u16(data[0x0C:0x0E]),
		UpperThresholdNonRecoverable: u16(data[0x0E:0x10]),
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
		mem.Handle = u16(data[0x08+offset : 0x0A+offset])
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
	infoCommon
	ChannelType        MemoryChannelType
	MaximumChannelLoad byte
	MemoryDeviceCount  byte
	LoadHandle         MemoryDeviceLoadHandles
}

func (m MemoryChannel) String() string {
	return fmt.Sprintf("Memory Channel\n"+
		"\tChannel Type: %s\n"+
		"\tMaximum Channel Load: %d\n"+
		"%s",
		m.ChannelType,
		m.MaximumChannelLoad,
		m.LoadHandle)
}

func (h dmiHeader) MemoryChannel() *MemoryChannel {
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
	return fmt.Sprintf("\tBase Address Modifier:\n"+
		"\t\tRegister spacing: %s\n"+
		"\t\tLs-bit for addresses: %d\n"+
		"\tInterrupt Info:\n"+
		"\t\tInfo: %s\n"+
		"\t\tPolarity: %s\n"+
		"\t\tTrigger Mode: %s",
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
	infoCommon
	InterfaceType                  IPMIDeviceInformationInterfaceType
	Revision                       byte
	I2CSlaveAddress                byte
	NVStorageAddress               byte
	BaseAddress                    uint64
	BaseAddressModiferInterrutInfo IPMIDeviceInformationAddressModiferInterruptInfo
	InterruptNumbe                 byte
}

func (i IPMIDeviceInformation) String() string {
	return fmt.Sprintf("IPMI Device Information\n"+
		"\tInterface Type: %s\n"+
		"\tRevision: %d\n"+
		"\tI2C Slave Address: %d\n"+
		"\tNV Storage Address: %d\n"+
		"\tBase Address: %d\n"+
		"\tBase Address Modifer Interrut Info: %s\n"+
		"\tInterrupt Numbe: %d",
		i.InterfaceType,
		i.Revision,
		i.I2CSlaveAddress,
		i.NVStorageAddress,
		i.BaseAddress,
		i.BaseAddressModiferInterrutInfo,
		i.InterruptNumbe)
}

func (h dmiHeader) IPMIDeviceInformation() *IPMIDeviceInformation {
	data := h.data
	return &IPMIDeviceInformation{
		InterfaceType:                  IPMIDeviceInformationInterfaceType(data[0x04]),
		Revision:                       data[0x05],
		I2CSlaveAddress:                data[0x06],
		NVStorageAddress:               data[0x07],
		BaseAddress:                    u64(data[0x08:0x10]),
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

func (h dmiHeader) SystemPowerSupply() *SystemPowerSupply {
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
		MaxPowerCapacity:           u16(data[0x0C:0x0E]),
		PowerSupplyCharacteristics: newSystemPowerSupplyCharacteristics(u16(data[0x0E : 0x0E+2])),
		InputVoltageProbeHandle:    u16(data[0x0F:0x11]),
		CoolingDeviceHandle:        u16(data[0x11:0x13]),
		InputCurrentProbeHandle:    u16(data[0x13:0x15]),
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
		str += fmt.Sprintf("\n\t\t\t\tReferenced Handle: %d\n"+
			"\t\t\t\tReferenced Offset: %d\n"+
			"\t\t\t\tString: %s\n"+
			"\t\t\t\tValue: %v",
			s.ReferencedHandle,
			s.ReferencedOffset,
			s.String,
			s.Value)
	}
	return str
}

type AdditionalInformation struct {
	infoCommon
	NumberOfEntries byte
	Entries         []AdditionalInformationEntries
}

func (a AdditionalInformation) String() string {
	return fmt.Sprintf("Additional Information\n"+
		"\tNumber Of Entries: %d\n"+
		"\tEntries: %s",
		a.NumberOfEntries,
		AdditionalInformationEntriess(a.Entries))
}

func (h dmiHeader) AdditionalInformation() *AdditionalInformation {
	data := h.data
	ai := new(AdditionalInformation)
	ai.NumberOfEntries = data[0x04]
	en := make([]AdditionalInformationEntries, 0)
	d := data[0x05:]
	for i := byte(0); i < ai.NumberOfEntries; i++ {
		var e AdditionalInformationEntries
		e.Length = d[0x0]
		e.ReferencedHandle = u16(d[0x01:0x03])
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
	infoCommon
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
	return fmt.Sprintf("On Board Devices Extended Information\n"+
		"\tReference Designation: %s\n"+
		"\tDevice Type: %s\n"+
		"\tDevice Type Instance: %d\n"+
		"%s\n",
		o.ReferenceDesignation,
		o.DeviceType,
		o.DeviceTypeInstance,
		o.SlotSegment())
}

func (h dmiHeader) OnBoardDevicesExtendedInformation() *OnBoardDevicesExtendedInformation {
	data := h.data
	return &OnBoardDevicesExtendedInformation{
		ReferenceDesignation: h.FieldString(int(data[0x04])),
		DeviceType:           OnBoardDevicesExtendedInformationType(data[0x05]),
		DeviceTypeInstance:   data[0x06],
		SegmentGroupNumber:   u16(data[0x07:0x09]),
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
	infoCommon
	Type ManagementControllerHostInterfaceType
	Data ManagementControllerHostInterfaceData
}

func (m ManagementControllerHostInterface) String() string {
	return fmt.Sprintf("Management Controller Host Interface\n"+
		"\tType: %s\n"+
		"\tMC Host Interface Data: %s\n",
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

func (h dmiHeader) ManagementControllerHostInterface() *ManagementControllerHostInterface {
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
	infoCommon
}

func (i Inactive) String() string {
	return "Inactive"
}

func (h dmiHeader) Inactive() *Inactive {
	return &Inactive{}
}

type EndOfTable struct {
	infoCommon
}

func (e EndOfTable) String() string {
	return "End-of-Table"
}

func (h dmiHeader) EndOfTable() *EndOfTable {
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

func u16(data []byte) uint16 {
	var u uint16
	binary.Read(bytes.NewBuffer(data[0:2]), binary.LittleEndian, &u)
	return u
}

func u32(data []byte) uint32 {
	var u uint32
	binary.Read(bytes.NewBuffer(data[0:4]), binary.LittleEndian, &u)
	return u
}

func u64(data []byte) uint64 {
	var u uint64
	binary.Read(bytes.NewBuffer(data[0:8]), binary.LittleEndian, &u)
	return u
}

func newdmiHeader(data []byte) *dmiHeader {
	if len(data) < 0x04 {
		return nil
	}
	return &dmiHeader{
		infoCommon: infoCommon{
			SMType: SMBIOSStructureType(data[0x00]),
			Length: data[1],
			Handle: SMBIOSStructureHandle(u16(data[0x02:0x04])),
		},
		data: data}
}

func newEntryPoint() (eps *entryPoint, err error) {
	eps = new(entryPoint)

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
	eps.MaxSize = u16(data[0x08:0x0A])
	eps.FormattedArea = data[0x0B:0x0F]
	eps.InterAnchor = data[0x10:0x15]
	eps.InterChecksum = data[0x15]
	eps.TableLength = u16(data[0x16:0x18])
	eps.TableAddress = u32(data[0x18:0x1C])
	eps.NumberOfSM = u16(data[0x1C:0x1E])
	eps.BCDRevision = data[0x1E]
	return
}

func (e entryPoint) StructureTableMem() ([]byte, error) {
	return getMem(e.TableAddress, uint32(e.TableLength))
}

func (h dmiHeader) Next() *dmiHeader {
	de := []byte{0, 0}
	next := h.data[h.Length:]
	index := bytes.Index(next, de)
	if index == -1 {
		return nil
	}
	return newdmiHeader(next[index+2:])
}

func (h dmiHeader) Decode() interface{} {
	switch h.SMType {
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

func (h dmiHeader) FieldString(offset int) string {
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

func (h dmiHeader) BIOSInformation() *BIOSInformation {
	data := h.data
	sas := u16(data[0x06:0x08])
	bi := &BIOSInformation{
		Vendor:                 h.FieldString(int(data[0x04])),
		BIOSVersion:            h.FieldString(int(data[0x05])),
		StartingAddressSegment: sas,
		ReleaseDate:            h.FieldString(int(data[0x08])),
		RomSize:                BIOSRomSize(64 * (data[0x09] + 1)),
		RuntimeSize:            BIOSRuntimeSize((uint(0x10000) - uint(sas)) << 4),
		Characteristics:        BIOSCharacteristics(u64(data[0x0A:0x12])),
	}
	if h.Length >= 0x13 {
		bi.CharacteristicsExt1 = BIOSCharacteristicsExt1(data[0x12])
	}
	if h.Length >= 0x14 {
		bi.CharacteristicsExt2 = BIOSCharacteristicsExt2(data[0x13])
	}
	return bi
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

func (h dmiHeader) SystemInformation() *SystemInformation {
	data := h.data
	version := h.FieldString(int(data[0x06]))
	return &SystemInformation{
		Manufacturer: h.FieldString(int(data[0x04])),
		ProductName:  h.FieldString(int(data[0x05])),
		Version:      version,
		SerialNumber: h.FieldString(int(data[0x07])),
		UUID:         uuid(data[0x08:0x18], version),
		WakeUpType:   SystemInformationWakeUpType(data[0x18]),
		SKUNumber:    h.FieldString(int(data[0x19])),
		Family:       h.FieldString(int(data[0x1A])),
	}
}

func (h dmiHeader) BaseboardInformation() *BaseboardInformation {
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

func (e entryPoint) StructureTable() map[SMBIOSStructureType]interface{} {
	tmem, err := e.StructureTableMem()
	if err != nil {
		return nil
	}
	m := make(map[SMBIOSStructureType]interface{})
	for hd := newdmiHeader(tmem); hd != nil; hd = hd.Next() {
		m[hd.SMType] = hd.Decode()
	}
	return m
}

func init() {
	eps, err := newEntryPoint()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		panic(err)
	}
	gdmi = eps.StructureTable()
}

func GetPortInformation() *PortInformation {
	if d, ok := gdmi[SMBIOSStructureTypePortConnector]; ok {
		return d.(*PortInformation)
	}
	return nil
}

func GetSystemSlot() *SystemSlot {
	if d, ok := gdmi[SMBIOSStructureTypeSystemSlots]; ok {
		return d.(*SystemSlot)
	}
	return nil
}

func GetOnBoardDeviceInformation() *OnBoardDeviceInformation {
	if d, ok := gdmi[SMBIOSStructureTypeOnBoardDevices]; ok {
		return d.(*OnBoardDeviceInformation)
	}
	return nil
}

func GetBIOSLanguageInformation() *BIOSLanguageInformation {
	if d, ok := gdmi[SMBIOSStructureTypeBIOSLanguage]; ok {
		return d.(*BIOSLanguageInformation)
	}
	return nil
}

func GetChassisInformation() *ChassisInformation {
	if d, ok := gdmi[SMBIOSStructureTypeChassis]; ok {
		return d.(*ChassisInformation)
	}
	return nil
}

func GetProcessorInformation() *ProcessorInformation {
	if d, ok := gdmi[SMBIOSStructureTypeProcessor]; ok {
		return d.(*ProcessorInformation)
	}
	return nil
}

func GetCacheInformation() *CacheInformation {
	if d, ok := gdmi[SMBIOSStructureTypeCache]; ok {
		return d.(*CacheInformation)
	}
	return nil
}

func GetSystemConfigurationOptions() *SystemConfigurationOptions {
	if d, ok := gdmi[SMBIOSStructureTypeSystemConfigurationOptions]; ok {
		return d.(*SystemConfigurationOptions)
	}
	return nil
}

func GetOEMStrings() *OEMStrings {
	if d, ok := gdmi[SMBIOSStructureTypeOEMStrings]; ok {
		return d.(*OEMStrings)
	}
	return nil
}

func GetGroupAssociations() *GroupAssociations {
	if d, ok := gdmi[SMBIOSStructureTypeGroupAssociations]; ok {
		return d.(*GroupAssociations)
	}
	return nil
}

func GetPhysicalMemoryArray() *PhysicalMemoryArray {
	if d, ok := gdmi[SMBIOSStructureTypePhysicalMemoryArray]; ok {
		return d.(*PhysicalMemoryArray)
	}
	return nil
}

func GetMemoryDevice() *MemoryDevice {
	if d, ok := gdmi[SMBIOSStructureTypeMemoryDevice]; ok {
		return d.(*MemoryDevice)
	}
	return nil
}

func Get_32BitMemoryErrorInformation() *_32BitMemoryErrorInformation {
	if d, ok := gdmi[SMBIOSStructureType32_bitMemoryError]; ok {
		return d.(*_32BitMemoryErrorInformation)
	}
	return nil
}

func GetBuiltinPointingDevice() *BuiltinPointingDevice {
	if d, ok := gdmi[SMBIOSStructureTypeBuilt_inPointingDevice]; ok {
		return d.(*BuiltinPointingDevice)
	}
	return nil
}

func GetPortableBattery() *PortableBattery {
	if d, ok := gdmi[SMBIOSStructureTypePortableBattery]; ok {
		return d.(*PortableBattery)
	}
	return nil
}

func GetSystemReset() *SystemReset {
	if d, ok := gdmi[SMBIOSStructureTypeSystemReset]; ok {
		return d.(*SystemReset)
	}
	return nil
}

func GetHardwareSecurity() *HardwareSecurity {
	if d, ok := gdmi[SMBIOSStructureTypeHardwareSecurity]; ok {
		return d.(*HardwareSecurity)
	}
	return nil
}

func GetSystemPowerControls() *SystemPowerControls {
	if d, ok := gdmi[SMBIOSStructureTypeSystemPowerControls]; ok {
		return d.(*SystemPowerControls)
	}
	return nil
}

func GetVoltageProbe() *VoltageProbe {
	if d, ok := gdmi[SMBIOSStructureTypeVoltageProbe]; ok {
		return d.(*VoltageProbe)
	}
	return nil
}

func GetCoolingDevice() *CoolingDevice {
	if d, ok := gdmi[SMBIOSStructureTypeCoolingDevice]; ok {
		return d.(*CoolingDevice)
	}
	return nil
}

func GetTemperatureProbe() *TemperatureProbe {
	if d, ok := gdmi[SMBIOSStructureTypeTemperatureProbe]; ok {
		return d.(*TemperatureProbe)
	}
	return nil
}

func GetElectricalCurrentProbe() *ElectricalCurrentProbe {
	if d, ok := gdmi[SMBIOSStructureTypeElectricalCurrentProbe]; ok {
		return d.(*ElectricalCurrentProbe)
	}
	return nil
}

func GetOutOfBandRemoteAccess() *OutOfBandRemoteAccess {
	if d, ok := gdmi[SMBIOSStructureTypeOut_of_bandRemoteAccess]; ok {
		return d.(*OutOfBandRemoteAccess)
	}
	return nil
}

func GetSystemBootInformation() *SystemBootInformation {
	if d, ok := gdmi[SMBIOSStructureTypeSystemBoot]; ok {
		return d.(*SystemBootInformation)
	}
	return nil
}

func Get_64BitMemoryErrorInformation() *_64BitMemoryErrorInformation {
	if d, ok := gdmi[SMBIOSStructureType64_bitMemoryError]; ok {
		return d.(*_64BitMemoryErrorInformation)
	}
	return nil
}

func GetManagementDevice() *ManagementDevice {
	if d, ok := gdmi[SMBIOSStructureTypeManagementDevice]; ok {
		return d.(*ManagementDevice)
	}
	return nil
}

func GetManagementDeviceComponent() *ManagementDeviceComponent {
	if d, ok := gdmi[SMBIOSStructureTypeManagementDeviceComponent]; ok {
		return d.(*ManagementDeviceComponent)
	}
	return nil
}

func GetManagementDeviceThresholdData() *ManagementDeviceThresholdData {
	if d, ok := gdmi[SMBIOSStructureTypeManagementDeviceThresholdData]; ok {
		return d.(*ManagementDeviceThresholdData)
	}
	return nil
}

func GetMemoryChannel() *MemoryChannel {
	if d, ok := gdmi[SMBIOSStructureTypeMemoryChannel]; ok {
		return d.(*MemoryChannel)
	}
	return nil
}

func GetIPMIDeviceInformation() *IPMIDeviceInformation {
	if d, ok := gdmi[SMBIOSStructureTypeIPMIDevice]; ok {
		return d.(*IPMIDeviceInformation)
	}
	return nil
}

func GetSystemPowerSupply() *SystemPowerSupply {
	if d, ok := gdmi[SMBIOSStructureTypePowerSupply]; ok {
		return d.(*SystemPowerSupply)
	}
	return nil
}

func GetAdditionalInformation() *AdditionalInformation {
	if d, ok := gdmi[SMBIOSStructureTypeAdditionalInformation]; ok {
		return d.(*AdditionalInformation)
	}
	return nil
}

func GetOnBoardDevicesExtendedInformation() *OnBoardDevicesExtendedInformation {
	if d, ok := gdmi[SMBIOSStructureTypeOnBoardDevicesExtendedInformation]; ok {
		return d.(*OnBoardDevicesExtendedInformation)
	}
	return nil
}

func GetManagementControllerHostInterface() *ManagementControllerHostInterface {
	if d, ok := gdmi[SMBIOSStructureTypeManagementControllerHostInterface]; ok {
		return d.(*ManagementControllerHostInterface)
	}
	return nil
}

func GetBIOSInformation() *BIOSInformation {
	if d, ok := gdmi[SMBIOSStructureTypeBIOS]; ok {
		return d.(*BIOSInformation)
	}
	return nil
}

func GetSystemInformation() *SystemInformation {
	if d, ok := gdmi[SMBIOSStructureTypeSystem]; ok {
		return d.(*SystemInformation)
	}
	return nil
}

func GetBaseboardInformation() *BaseboardInformation {
	if d, ok := gdmi[SMBIOSStructureTypeBaseBoard]; ok {
		return d.(*BaseboardInformation)
	}
	return nil
}

func GetGDMI() map[SMBIOSStructureType]interface{} {
	return gdmi
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
