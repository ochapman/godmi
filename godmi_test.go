package godmi_test

import (
	. "github.com/ochapman/godmi"
	"log"
	"os/exec"
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
	m := make(map[string]string, 0)
	m["bios-vendor"] = bi.Vendor
	m["bios-version"] = bi.BIOSVersion
	m["bios-release-date"] = bi.ReleaseDate

	for k, v := range m {
		dmiv := dmidecode_s(k)
		if dmiv != v {
			t.Errorf("%s: \n[godmi]: %s\n[dmidecode]: %s\n", k, v, dmiv)
		}
	}

}

func TestSystem(t *testing.T) {
	si := GetSystemInformation()
	m := make(map[string]string, 0)
	m["system-manufacturer"] = si.Manufacturer
	m["system-product-name"] = si.ProductName
	m["system-version"] = si.Version
	m["system-serial-number"] = si.SerialNumber
	m["system-uuid"] = si.UUID

	for k, v := range m {
		dmiv := dmidecode_s(k)
		if dmiv != v {
			t.Errorf("%s: \n[godmi]: %s\n[dmidecode]: %s\n", k, v, dmiv)
		}
	}
}

func TestBaseboard(t *testing.T) {
	bi := GetBaseboardInformation()
	m := make(map[string]string, 0)
	m["baseboard-manufacturer"] = bi.Manufacturer
	m["baseboard-product-name"] = bi.Product
	m["baseboard-version"] = bi.Version
	m["baseboard-serial-number"] = bi.SerialNumber
	m["baseboard-asset-tag"] = bi.AssetTag

	for k, v := range m {
		dmiv := dmidecode_s(k)
		if dmiv != v {
			t.Errorf("%s: \n[godmi]: %s\n[dmidecode]: %s\n", k, v, dmiv)
		}
	}
}
