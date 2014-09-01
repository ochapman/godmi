/*
* godmi.go
* godmi command
*
* Chapman Ou <ochapman.cn@gmail.com>
*
* Thu Jul 31 22:44:14 CST 2014
 */

package main

import (
	"fmt"
	"github.com/ochapman/godmi"
	"reflect"
)

func main() {
	godmi.Init()
	infos := []interface{}{
		godmi.GetPortInformation(),
		godmi.GetSystemSlot(),
		godmi.GetOnBoardDeviceInformation(),
		godmi.GetBIOSLanguageInformation(),
		godmi.GetChassisInformation(),
		godmi.GetProcessorInformation(),
		godmi.GetCacheInformation(),
		godmi.GetSystemConfigurationOptions(),
		godmi.GetOEMStrings(),
		godmi.GetGroupAssociations(),
		godmi.GetPhysicalMemoryArray(),
		godmi.GetMemoryDevice(),
		godmi.Get_32BitMemoryErrorInformation(),
		godmi.GetBuiltinPointingDevice(),
		godmi.GetPortableBattery(),
		godmi.GetSystemReset(),
		godmi.GetHardwareSecurity(),
		godmi.GetSystemPowerControls(),
		godmi.GetVoltageProbe(),
		godmi.GetCoolingDevice(),
		godmi.GetTemperatureProbe(),
		godmi.GetElectricalCurrentProbe(),
		godmi.GetOutOfBandRemoteAccess(),
		godmi.GetSystemBootInformation(),
		godmi.Get_64BitMemoryErrorInformation(),
		godmi.GetManagementDevice(),
		godmi.GetManagementDeviceComponent(),
		godmi.GetManagementDeviceThresholdData(),
		godmi.GetMemoryChannel(),
		godmi.GetIPMIDeviceInformation(),
		godmi.GetSystemPowerSupply(),
		godmi.GetAdditionalInformation(),
		godmi.GetOnBoardDevicesExtendedInformation(),
		godmi.GetManagementControllerHostInterface(),
		godmi.GetBIOSInformation(),
		godmi.GetSystemInformation(),
		godmi.GetBaseboardInformation(),
	}
	for _, info := range infos {
		rv := reflect.ValueOf(info)
		if rv.IsNil() {
			continue
		}
		fmt.Println(info)
	}
}
