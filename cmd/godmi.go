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
	si := godmi.SystemInformation()
	fmt.Println(si.UUID)
	fmt.Println(si.ProductName)
	bi := godmi.BIOSInformation()
	fmt.Println(bi.Vendor)
	bo := godmi.GetBaseboardInformation()
	fmt.Println(bo.Manufacturer)
	ch := godmi.Chassis()
	fmt.Println(ch.Manufacturer)
}
