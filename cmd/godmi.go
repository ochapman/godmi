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
	ch := godmi.GetChassis()
	fmt.Println(ch.Manufacturer)
}
