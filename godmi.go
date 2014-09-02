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
	"fmt"
	"os"
	"strconv"
	"sync"
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

func (h dmiHeader) OutOfBandRemoteAccess() *OutOfBandRemoteAccess {
	data := h.data
	return &OutOfBandRemoteAccess{
		ManufacturerName: h.FieldString(int(data[0x04])),
		Connections:      NewOutOfBandRemoteAccessConnections(data[0x05]),
	}
}

func (h dmiHeader) SystemBootInformation() *SystemBootInformation {
	data := h.data
	return &SystemBootInformation{
		BootStatus: SystemBootInformationStatus(data[0x0A]),
	}
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

func (h dmiHeader) ManagementDevice() *ManagementDevice {
	data := h.data
	return &ManagementDevice{
		Description: h.FieldString(int(data[0x04])),
		Type:        ManagementDeviceType(data[0x05]),
		Address:     u32(data[0x06:0x0A]),
		AddressType: ManagementDeviceAddressType(data[0x0A]),
	}
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

func (h dmiHeader) Inactive() *Inactive {
	return &Inactive{}
}

func (h dmiHeader) EndOfTable() *EndOfTable {
	return &EndOfTable{}
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

func (h dmiHeader) newType() (interface{}, error) {
	t := h.SMType
	newfn, err := getTypeFunc(t)
	if err != nil {
		return nil, err
	}
	return newfn(h), nil
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

func (e entryPoint) StructureTable() map[SMBIOSStructureType]interface{} {
	tmem, err := e.StructureTableMem()
	if err != nil {
		return nil
	}
	m := make(map[SMBIOSStructureType]interface{})
	for hd := newdmiHeader(tmem); hd != nil; hd = hd.Next() {
		newtype, err := hd.newType()
		if err != nil {
			continue
		}
		m[hd.SMType] = newtype
	}
	return m
}

type dmiTyper interface {
	String() string
}

type newFunction func(d dmiHeader) dmiTyper

type typeFunc map[SMBIOSStructureType]newFunction

var g_typeFunc = make(typeFunc)
var g_wg sync.WaitGroup

func addTypeFunc(t SMBIOSStructureType, f newFunction) {
	g_typeFunc[t] = f
}

func getTypeFunc(t SMBIOSStructureType) (fn newFunction, err error) {
	fn, ok := g_typeFunc[t]
	if !ok {
		return fn, fmt.Errorf("type %s have no NewFunction", t)
	}
	return fn, nil
}

func Init() {
	eps, err := newEntryPoint()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		panic(err)
	}
	gdmi = eps.StructureTable()
}

func GetCacheInformation() *CacheInformation {
	if d, ok := gdmi[SMBIOSStructureTypeCache]; ok {
		return d.(*CacheInformation)
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
