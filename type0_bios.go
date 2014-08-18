/*
* File Name:	type0_bios.go
* Description:	
* Author:	Chapman Ou <ochapman.cn@gmail.com>
* Created:	2014-08-18 22:45:25
*/

package godmi

import (
	"fmt"
)

type BIOSCharacteristics uint64

func (b BIOSCharacteristics) String() string {
	var s string
	chars := [...]string{
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

	for i := uint32(4); i < 32; i++ {
		if b&(1<<i) != 0 {
			s += "\n\t\t" + chars[i-3]
		}
	}
	return s
}

type BIOSCharacteristicsExt1 byte

func (b BIOSCharacteristicsExt1) String() string {
	var s string
	chars := [...]string{
		"ACPI is supported", /* 0 */
		"USB legacy is supported",
		"AGP is supported",
		"I2O boot is supported",
		"LS-120 boot is supported",
		"ATAPI Zip drive boot is supported",
		"IEEE 1394 boot is supported",
		"Smart battery is supported", /* 7 */
	}

	for i := uint32(0); i < 7; i++ {
		if b&(1<<i) != 0 {
			s += "\n\t\t" + chars[i]
		}
	}
	return s
}

type BIOSCharacteristicsExt2 byte

func (b BIOSCharacteristicsExt2) String() string {
	var s string
	chars := [...]string{
		"BIOS boot specification is supported", /* 0 */
		"Function key-initiated network boot is supported",
		"Targeted content distribution is supported",
		"UEFI is supported",
		"System is a virtual machine", /* 4 */
	}

	for i := uint32(0); i < 5; i++ {
		if b&(1<<i) != 0 {
			s += "\n\t\t" + chars[i]
		}
	}
	return s
}

type BIOSRuntimeSize uint

func (b BIOSRuntimeSize) String() string {
	if (b & 0x3FF) > 0 {
		return fmt.Sprintf("%d Bytes", b)
	}
	return fmt.Sprintf("%d kB", b>>10)
}

type BIOSRomSize byte

func (b BIOSRomSize) String() string {
	return fmt.Sprintf("%d kB", b)
}

type BIOSInformation struct {
	infoCommon
	Vendor                                 string
	BIOSVersion                            string
	StartingAddressSegment                 uint16
	ReleaseDate                            string
	RomSize                                BIOSRomSize
	RuntimeSize                            BIOSRuntimeSize
	Characteristics                        BIOSCharacteristics
	CharacteristicsExt1                    BIOSCharacteristicsExt1
	CharacteristicsExt2                    BIOSCharacteristicsExt2
	SystemBIOSMajorRelease                 byte
	SystemBIOSMinorRelease                 byte
	EmbeddedControllerFirmwareMajorRelease byte
	EmbeddedControllerFirmawreMinorRelease byte
}

func (b BIOSInformation) String() string {
	s := fmt.Sprintf("BIOS Information\n"+
		"\tVendor: %s\n"+
		"\tVersion: %s\n"+
		"\tRelease Date: %s\n"+
		"\tAddress: 0x%4X0\n"+
		"\tRuntime Size: %s\n"+
		"\tROM Size: %s\n"+
		"\tCharacteristics:%s",
		b.Vendor,
		b.BIOSVersion,
		b.ReleaseDate,
		b.StartingAddressSegment,
		b.RuntimeSize,
		b.RomSize,
		b.Characteristics)

	if b.CharacteristicsExt1 != 0 {
		s += b.CharacteristicsExt1.String()
	}
	if b.CharacteristicsExt2 != 0 {
		s += b.CharacteristicsExt2.String()
	}
	return s
}
