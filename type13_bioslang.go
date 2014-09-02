/*
* File Name:	type13_bioslang.go
* Description:
* Author:	Chapman Ou <ochapman.cn@gmail.com>
* Created:	2014-08-19
 */
package godmi

import (
	"fmt"
)

type BIOSLanguageInformationFlag byte

const (
	BIOSLanguageInformationFlagLongFormat BIOSLanguageInformationFlag = iota
	BIOSLanguageInformationFlagAbbreviatedFormat
)

func NewBIOSLanguageInformationFlag(f byte) BIOSLanguageInformationFlag {
	return BIOSLanguageInformationFlag(f & 0xFE)
}

type BIOSLanguageInformation struct {
	infoCommon
	InstallableLanguage []string
	Flags               BIOSLanguageInformationFlag
	CurrentLanguage     string
}

func (b BIOSLanguageInformation) String() string {
	return fmt.Sprintf("BIOS Language Information:\n"+
		"\tInstallable Languages %s\n"+
		"\tFlags: %s\n"+
		"\tCurrent Language: %s",
		b.InstallableLanguage,
		b.Flags,
		b.CurrentLanguage)
}

func newBIOSLanguageInformation(h dmiHeader) dmiTyper {
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

func GetBIOSLanguageInformation() *BIOSLanguageInformation {
	if d, ok := gdmi[SMBIOSStructureTypeBIOSLanguage]; ok {
		return d.(*BIOSLanguageInformation)
	}
	return nil
}

func init() {
	addTypeFunc(SMBIOSStructureTypeBIOSLanguage, newBIOSLanguageInformation)
}
