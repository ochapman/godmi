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

func newSystemInformation(h dmiHeader) dmiTyper {
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

func GetSystemInformation() *SystemInformation {
	if d, ok := gdmi[SMBIOSStructureTypeSystem]; ok {
		return d.(*SystemInformation)
	}
	return nil
}

func init() {
	addTypeFunc(SMBIOSStructureTypeSystem, newSystemInformation)
}
