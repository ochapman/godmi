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
