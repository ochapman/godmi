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

type BIOSInformation struct {
	Type                                   byte
	Length                                 byte
	Handle                                 uint16
	Vendor                                 string
	BIOSVersion                            string
	StartingAddressSegment                 uint16
	ReleaseDate                            string
	RomSize                                byte
	Characteristics                        uint64
	CharacteristicsExt                     []byte
	SystemBIOSMajorRelease                 byte
	SystemBIOSMinorRelease                 byte
	EmbeddedControllerFirmwareMajorRelease byte
	EmbeddedControllerFirmawreMinorRelease byte
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
	WakeUpType   byte
	SKUNumber    string
	Family       string
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
	FeatureFlags                   byte
	LocationInChassis              string
	ChassisHandle                  uint16
	BoardType                      byte
	NumberOfContainedObjectHandles byte
	ContainedObjectHandles         []byte
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
	FeatureFlagsHostingBoard = iota
	FeatureFlagsAtLeastOneDaughter = 1 << 1
	FeatureFlagsRemovable = 1 << 2
	FeatureFlagsRepleaceable = 1 << 3
	FeatureFlagsHotSwappable = 1 << 4
	//FeatureFlagsReserved = 000b
)

const (
	BoardTypeUnknown = iota + 1
	BoardTypeOther
	BoardTypeServerBlade
	BoardTypeConnectivitySwitch
	BoardTypeSystemManagementModule
	BoardTypeProcessorModule
	BoardTypeIOModule
	BoardTypeMemoryModule
	BoardTypeDaughterboard
	BoardTypeMotherboard
	BoardTypeProcessorMemoryModule
	BoardTypeProcessorIOModule
	BoardTypeInterconnectboard
)

func U16(data []byte) uint16 {
	var u16 uint16
	binary.Read(bytes.NewBuffer(data[0:2]), binary.LittleEndian, &u16)
	return u16 
}

func NewDMIHeader(data []byte) DMIHeader {
	var h uint16
	binary.Read(bytes.NewBuffer(data[2:4]), binary.LittleEndian, &h)
	hd := DMIHeader{Type: data[0], Length: data[1], Handle: h, data: data}
	return hd
}

func NewSMBIOS_EPS() SMBIOS_EPS {
	var eps SMBIOS_EPS
	var u16 uint16
	var u32 uint32

	mem, err := getMem(0xF0000, 0x10000)
	if err != nil {
		return SMBIOS_EPS{}
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
	return eps
}

func (e SMBIOS_EPS) StructrueTableMem() ([]byte, error) {
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
	default:
		fmt.Println("Unknown")
	}
}

func (h DMIHeader) FieldString(offset int) string {
	d := h.data
	index := int(h.Length)
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
	return bi
}

func (bi BIOSInformation) String() string {
	return fmt.Sprintf("BIOS Information\n\tVendor: %s\n\tVersion: %s\n\tAddress: %4X0", bi.Vendor, bi.BIOSVersion, bi.StartingAddressSegment)
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
	//si.UUID
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
	bi.LocationInChassis = h.FieldString(int(data[0x0A]))
	return bi
}

func (e SMBIOS_EPS) StructureTable() {
	tmem, err := e.StructrueTableMem()
	if err != nil {
		return
	}
	//for i := 0, hd := NewDMIHeader(tmem); i < e.NumberOfSM ; i++, hd = hd.Next() {
	hd := NewDMIHeader(tmem)
	for i := 0; i < 3; i++ {
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
	eps := NewSMBIOS_EPS()
	eps.StructureTable()
	//fmt.Printf("%2X", m)
}
