/*
* File Name:	type1_system.go
* Description:	
* Author:	Chapman Ou <ochapman.cn@gmail.com>
* Created:	2014-08-18 22:52:15
*/

package godmi

import (
	"fmt"
)

type SystemInformationWakeUpType byte

const (
	SystemInformationReserved SystemInformationWakeUpType = iota
	SystemInformationOther
	SystemInformationUnknown
	SystemInformationAPM_Timer
	SystemInformationModem_Ring
	SystemInformationLAN_Remote
	SystemInformationPower_Switch
	SystemInformationPCI_PME
	SystemInformationAC_Power_Restored
)

func (w SystemInformationWakeUpType) String() string {
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

type SystemInformation struct {
	infoCommon
	Manufacturer string
	ProductName  string
	Version      string
	SerialNumber string
	UUID         string
	WakeUpType   SystemInformationWakeUpType
	SKUNumber    string
	Family       string
}

func (s SystemInformation) String() string {
	return fmt.Sprintf("System Information\n"+
		"\tManufacturer: %s\n"+
		"\tProduct Name: %s\n"+
		"\tVersion: %s\n"+
		"\tSerial Number: %s\n"+
		"\tUUID: %s\n"+
		"\tWake-up Type: %s\n"+
		"\tSKU Number: %s\n"+
		"\tFamily: %s",
		s.Manufacturer,
		s.ProductName,
		s.Version,
		s.SerialNumber,
		s.UUID,
		s.WakeUpType,
		s.SKUNumber,
		s.Family)
}

