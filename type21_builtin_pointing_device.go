/*
* File Name:	type21_builtin_pointing_device.go
* Description:
* Author:	Chapman Ou <ochapman.cn@gmail.com>
* Created:	2014-08-19
 */
package godmi

import (
	"fmt"
)

type BuiltinPointingDeviceType byte

const (
	BuiltinPointingDeviceTypeOther BuiltinPointingDeviceType = 1 + iota
	BuiltinPointingDeviceTypeUnknown
	BuiltinPointingDeviceTypeMouse
	BuiltinPointingDeviceTypeTrackBall
	BuiltinPointingDeviceTypeTrackPoint
	BuiltinPointingDeviceTypeGlidePoint
	BuiltinPointingDeviceTypeTouchPad
	BuiltinPointingDeviceTypeTouchScreen
	BuiltinPointingDeviceTypeOpticalSensor
)

func (b BuiltinPointingDeviceType) String() string {
	types := [...]string{
		"Other",
		"Unknown",
		"Mouse",
		"Track Ball",
		"Track Point",
		"Glide Point",
		"Touch Pad",
		"Touch Screen",
		"Optical Sensor",
	}
	return types[b-1]
}

type BuiltinPointingDeviceInterface byte

const (
	BuiltinPointingDeviceInterfaceOther BuiltinPointingDeviceInterface = 1 + iota
	BuiltinPointingDeviceInterfaceUnknown
	BuiltinPointingDeviceInterfaceSerial
	BuiltinPointingDeviceInterfacePS2
	BuiltinPointingDeviceInterfaceInfrared
	BuiltinPointingDeviceInterfaceHP_HIL
	BuiltinPointingDeviceInterfaceBusmouse
	BuiltinPointingDeviceInterfaceADB
	BuiltinPointingDeviceInterfaceBusmouseDB_9
	BuiltinPointingDeviceInterfaceBusmousemicro_DIN
	BuiltinPointingDeviceInterfaceUSB
)

func (b BuiltinPointingDeviceInterface) String() string {
	interfaces := [...]string{
		"Other",
		"Unknown",
		"Serial",
		"PS/2",
		"Infrared",
		"HP-HIL",
		"Bus mouse",
		"ADB (Apple Desktop Bus)",
		"Bus mouse DB-9",
		"Bus mouse micro-DIN",
		"USB",
	}
	return interfaces[b-1]
}

type BuiltinPointingDevice struct {
	infoCommon
	Type            BuiltinPointingDeviceType
	Interface       BuiltinPointingDeviceInterface
	NumberOfButtons byte
}

func (b BuiltinPointingDevice) String() string {
	return fmt.Sprintf("Built-in Pointing Device\n"+
		"\tType: %s\n"+
		"\tInterface: %s\n"+
		"\tNumber of Buttons: %d",
		b.Type,
		b.Interface,
		b.NumberOfButtons,
	)
}

func newBuiltinPointingDevice(h dmiHeader) dmiTyper {
	data := h.data
	return &BuiltinPointingDevice{
		Type:            BuiltinPointingDeviceType(data[0x04]),
		Interface:       BuiltinPointingDeviceInterface(data[0x05]),
		NumberOfButtons: data[0x06],
	}
}

func GetBuiltinPointingDevice() *BuiltinPointingDevice {
	if d, ok := gdmi[SMBIOSStructureTypeBuilt_inPointingDevice]; ok {
		return d.(*BuiltinPointingDevice)
	}
	return nil
}

func init() {
	addTypeFunc(SMBIOSStructureTypeBuilt_inPointingDevice, newBuiltinPointingDevice)
}
