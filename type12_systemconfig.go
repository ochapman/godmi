/*
* File Name:	type12_systemconfig.go
* Description:	
* Author:	Chapman Ou <ochapman.cn@gmail.com>
* Created:	2014-08-19
*/
package godmi

import (
	"fmt"
)

type SystemConfigurationOptions struct {
	infoCommon
	Count   byte
	strings string
}

func (s SystemConfigurationOptions) String() string {
	return fmt.Sprintf("System Configuration Option\n\t\t%s", s.strings)
}
