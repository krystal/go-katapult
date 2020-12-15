package buildspec_test

import (
	"fmt"
	"strings"

	"github.com/krystal/go-katapult/pkg/buildspec"
)

func ExampleFromJSON() {
	r := strings.NewReader(`
{
  "hostname": "web-3",
  "disk_template": {
    "permalink": "templates/ubuntu-18-04"
  },
  "resources": {
    "package": {
      "permalink": "rock-3"
    }
  },
  "data_center": {
    "permalink": "london"
  }
}`)

	spec, _ := buildspec.FromJSON(r)

	fmt.Println(spec.DataCenter.Permalink)
	fmt.Println(spec.Resources.Package.Permalink)
	fmt.Println(spec.DiskTemplate.Permalink)
	fmt.Println(spec.Hostname)
	// Output:
	// london
	// rock-3
	// templates/ubuntu-18-04
	// web-3
}

func ExampleFromXML() {
	r := strings.NewReader(`
<?xml version="1.0" encoding="UTF-8"?>
<VirtualMachineSpec>
  <DataCenter by="permalink">london</DataCenter>
  <Resources>
    <Package by="permalink">rock-3</Package>
  </Resources>
  <DiskTemplate>
    <DiskTemplate by="permalink">templates/ubuntu-18-04</DiskTemplate>
  </DiskTemplate>
  <Hostname>
    <Hostname>web-3</Hostname>
  </Hostname>
</VirtualMachineSpec>`)

	spec, _ := buildspec.FromXML(r)

	fmt.Println(spec.DataCenter.Permalink)
	fmt.Println(spec.Resources.Package.Permalink)
	fmt.Println(spec.DiskTemplate.Permalink)
	fmt.Println(spec.Hostname)
	// Output:
	// london
	// rock-3
	// templates/ubuntu-18-04
	// web-3
}

func ExampleFromYAML() {
	r := strings.NewReader(`
data_center:
  permalink: london
resources:
  package:
    permalink: rock-3
disk_template:
  permalink: templates/ubuntu-18-04
hostname: web-3`)

	spec, _ := buildspec.FromYAML(r)

	fmt.Println(spec.DataCenter.Permalink)
	fmt.Println(spec.Resources.Package.Permalink)
	fmt.Println(spec.DiskTemplate.Permalink)
	fmt.Println(spec.Hostname)
	// Output:
	// london
	// rock-3
	// templates/ubuntu-18-04
	// web-3
}
