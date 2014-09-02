/*
* File Name:	type11_oem.go
* Description:
* Author:	Chapman Ou <ochapman.cn@gmail.com>
* Created:	2014-08-19
 */
package godmi

import (
	"fmt"
)

type OEMStrings struct {
	infoCommon
	Count   byte
	strings string
}

func (o OEMStrings) String() string {
	return fmt.Sprintf("OEM strings: %s", o.strings)
}

func newOEMStrings(h dmiHeader) dmiTyper {
	var o OEMStrings
	data := h.data
	o.Count = data[0x04]
	for i := byte(0); i < o.Count; i++ {
		o.strings += fmt.Sprintf("strings: %d %s\n\t\t", i, h.FieldString(int(data[i])))
	}
	return &o
}

func GetOEMStrings() *OEMStrings {
	if d, ok := gdmi[SMBIOSStructureTypeOEMStrings]; ok {
		return d.(*OEMStrings)
	}
	return nil
}

func init() {
	addTypeFunc(SMBIOSStructureTypeOEMStrings, newOEMStrings)
}
