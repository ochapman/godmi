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
)

func main() {
	si := godmi.GetSystemInformation()
	fmt.Println(si.UUID)
	fmt.Println(si.ProductName)
	bi := godmi.GetBIOSInformation()
	fmt.Println(bi.Vendor)
	bo := godmi.GetBaseboardInformation()
	fmt.Println(bo.Manufacturer)
	ch := godmi.GetChassisInformation()
	fmt.Println(ch.Manufacturer)
	ca := godmi.GetCacheInformation()
	if ca != nil {
		fmt.Println(ca)
	}
	if d := godmi.GetPortInformation(); d != nil {
		fmt.Println(d)
	}
}
