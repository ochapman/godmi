package godmi_test

import (
	"bufio"
	"bytes"
	"fmt"
	. "github.com/ochapman/godmi"
	"log"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func dmidecode(arg ...string) string {
	output, err := exec.Command("dmidecode", arg...).Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(output)
}

func init() {
	Init()
}

func dmidecode_s(kw string) string {
	output := dmidecode("-s", kw)
	return strings.TrimSpace(output)
}

func dmidecode_t(kw string) string {
	var output string
	dd := dmidecode("-q", "-t", kw)
	// Remove empty line
	r := bytes.NewReader([]byte(dd))
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if line != "" {
			if len(output) > 0 {
				output = output + "\n" + line
			} else {
				output = line
			}
		}
	}
	return output
}

func compare(m map[string]string, t *testing.T) {
	for k, v := range m {
		dmiv := dmidecode_s(k)
		if dmiv != v {
			t.Errorf("%s: \n[godmi]: %s\n[dmidecode]: %s\n", k, v, dmiv)
		}
	}
}
func checkInfo(i interface{}, kw string, t *testing.T) {
	rv := reflect.ValueOf(i)
	if rv.IsNil() {
		dv := dmidecode_s(kw)
		if len(dv) == 0 {
			t.Skip("dmidecode and godmi has no data")
		} else {
			t.Errorf("dmidecode has %s, but godmi has no data\n", dv)
		}
	}
}

/*
dmidecode command has following STRING keywords:
  bios-vendor
  bios-version
  bios-release-date

  system-manufacturer
  system-product-name
  system-version
  system-serial-number
  system-uuid

  baseboard-manufacturer
  baseboard-product-name
  baseboard-version
  baseboard-serial-number
  baseboard-asset-tag

  chassis-manufacturer
  chassis-type
  chassis-version
  chassis-serial-number
  chassis-asset-tag

  processor-family
  processor-manufacturer
  processor-version
  processor-frequency
*/

func TestBIOS(t *testing.T) {
	bi := GetBIOSInformation()
	checkInfo(bi, "bios-vendor", t)
	m := map[string]string{
		"bios-vendor":       bi.Vendor,
		"bios-version":      bi.BIOSVersion,
		"bios-release-date": bi.ReleaseDate,
	}

	compare(m, t)
}

func TestSystem(t *testing.T) {
	si := GetSystemInformation()
	if si == nil {
		t.Skip("GetSystemInformation() is nil")
	}
	checkInfo(si, "system-manufacturer", t)
	m := map[string]string{
		"system-manufacturer":  si.Manufacturer,
		"system-product-name":  si.ProductName,
		"system-version":       si.Version,
		"system-serial-number": si.SerialNumber,
		"system-uuid":          si.UUID,
	}
	compare(m, t)
}

func TestBaseboard(t *testing.T) {
	bi := GetBaseboardInformation()
	if bi == nil {
		t.Skip("GetBaseBoardInformation() is nil")
	}
	checkInfo(bi, "baseboard-manufacturer", t)
	m := map[string]string{
		"baseboard-manufacturer":  bi.Manufacturer,
		"baseboard-product-name":  bi.ProductName,
		"baseboard-version":       bi.Version,
		"baseboard-serial-number": bi.SerialNumber,
		"baseboard-asset-tag":     bi.AssetTag,
	}
	compare(m, t)
}

func TestChassis(t *testing.T) {
	ci := GetChassisInformation()
	if ci == nil {
		t.Skip("GetChassisInformation() is nil")
	}
	checkInfo(ci, "chassis-manufacturer", t)
	m := map[string]string{
		"chassis-manufacturer":  ci.Manufacturer,
		"chassis-type":          ci.Type.String(),
		"chassis-version":       ci.Version,
		"chassis-serial-number": ci.SerialNumber,
		"chassis-asset-tag":     ci.AssetTag,
	}
	compare(m, t)
}

func TestProcessor(t *testing.T) {
	pi := GetProcessorInformation()
	if pi == nil {
		t.Skip("GetProcessorInformation() is nil")
	}
	checkInfo(pi, "processor-family", t)
	m := map[string]string{
		"processor-family":       pi.Family.String(),
		"processor-manufacturer": pi.Manufacturer,
		"processor-version":      pi.Version,
		"processor-frequency":    strconv.Itoa(int(pi.MaxSpeed)),
	}
	compare(m, t)
}

/*
dmidecode has following TYPE keywords:
	bios
	system
	baseboard
	chassis
	processor
	memory
	cache
	connector
	slot
*/

func TestType(t *testing.T) {
	m := map[string]interface{}{
		"bios":      GetBIOSInformation(),
		"system":    GetSystemInformation(),
		"baseboard": GetBaseboardInformation(),
		"chassis":   GetChassisInformation(),
		"processor": GetProcessorInformation(),
		"memory":    GetMemoryDevice(),
		"cache":     GetCacheInformation(),
		"connector": GetPortInformation(),
		"slot":      GetSystemSlot(),
	}
	for k, v := range m {
		dv := dmidecode_t(k)
		vv := reflect.ValueOf(v)
		if vv.IsNil() {
			if len(dv) == 0 {
				t.Logf("%s: dmidecode and godmi has no data", k)
				continue
			} else {
				t.Errorf("%s:\n[godmi]: nil\n[dmidecode]: %s\n", k, dv)
			}
		}
		gv := fmt.Sprintf("%s", v)
		if gv != dv {
			t.Errorf("%s:\n[godmi]:\n%s\n[dmidecode]:\n%s\n", k, gv, dv)
		}
	}
}
