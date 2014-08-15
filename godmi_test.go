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

func TestBIOSVendor(t *testing.T) {
	v := dmidecode("-s", "bios-vendor")
	// v has "\n"
	vendor := strings.TrimSpace(v)
	bi := GetBIOSInformation()
	if vendor != bi.Vendor {
		t.Errorf("bios-vendor godmi: %s dmidecode: %s", bi.Vendor, vendor)
	}
}
