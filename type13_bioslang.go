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
