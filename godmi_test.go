package godmi_test

import (
	. "github.com/ochapman/godmi"
	"log"
	"os/exec"
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

func dmidecode_s(kw string) string {
	output := dmidecode("-s", kw)
	return strings.TrimSpace(output)
}

func compare(m map[string]string, t *testing.T) {
	for k, v := range m {
		dmiv := dmidecode_s(k)
		if dmiv != v {
			t.Errorf("%s: \n[godmi]: %s\n[dmidecode]: %s\n", k, v, dmiv)
		}
	}
}

/*
dmidecode command has following keywords:
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
	if bi == nil {
		t.Skip("GetBIOSInformation failed")
	}
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
		t.Skip("GetSystemInformation failed")
	}
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
		t.Skip("GetBaseboardInformation failed")
	}
	m := map[string]string{
		"baseboard-manufacturer":  bi.Manufacturer,
		"baseboard-product-name":  bi.Product,
		"baseboard-version":       bi.Version,
		"baseboard-serial-number": bi.SerialNumber,
		"baseboard-asset-tag":     bi.AssetTag,
	}
	compare(m, t)
}

func TestChassis(t *testing.T) {
	ci := GetChassisInformation()
	if ci == nil {
		t.Skip("GetChassisInformation failed")
	}
	m := map[string]string{
		"chassis-manufacturer":  ci.Manufacturer,
		"chassis-type":          ci.ChassisType.String(),
		"chassis-version":       ci.Version,
		"chassis-serial-number": ci.SerialNumber,
		"chassis-asset-tag":     ci.AssetTag,
	}
	compare(m, t)
}

func TestProcessor(t *testing.T) {
	pi := GetProcessorInformation()
	if pi == nil {
		t.Skip("GetProcessorInformation failed")
	}
	m := map[string]string{
		"processor-family":       pi.Family.String(),
		"processor-manufacturer": pi.Manufacturer,
		"processor-version":      pi.Version,
		"processor-frequency":    strconv.Itoa(int(pi.MaxSpeed)),
	}
	compare(m, t)
}
