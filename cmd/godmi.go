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
	"os"
)

func main() {
	eps, err := godmi.NewSMBIOS_EPS()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
	m := eps.StructureTable()
	system := m[godmi.SMBIOSStructureTypeSystem].(godmi.SystemInformation)
	fmt.Println(system.UUID)
	fmt.Println(system.ProductName)
	//fmt.Printf("%2X", m)
}
