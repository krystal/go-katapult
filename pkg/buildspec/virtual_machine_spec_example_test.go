package buildspec

import (
	"fmt"
)

//nolint:lll
func ExampleVirtualMachineSpec_JSON() {
	spec := &VirtualMachineSpec{
		DataCenter: &DataCenter{Permalink: "london"},
		Resources: &Resources{
			Package: &Package{Permalink: "rock-3"},
		},
		DiskTemplate: &DiskTemplate{Permalink: "templates/ubuntu-18-04"},
		Hostname:     "web-3",
	}
	x, _ := spec.JSON()

	fmt.Println(string(x))
	// Output:
	// {"data_center":{"permalink":"london"},"resources":{"package":{"permalink":"rock-3"}},"disk_template":{"permalink":"templates/ubuntu-18-04"},"hostname":"web-3"}
}

//nolint:lll
func ExampleVirtualMachineSpec_XML() {
	spec := &VirtualMachineSpec{
		DataCenter: &DataCenter{Permalink: "london"},
		Resources: &Resources{
			Package: &Package{Permalink: "rock-3"},
		},
		DiskTemplate: &DiskTemplate{Permalink: "templates/ubuntu-18-04"},
		Hostname:     "web-3",
	}
	x, _ := spec.XML()

	fmt.Println(string(x))
	// Output:
	// <?xml version="1.0" encoding="UTF-8"?>
	// <VirtualMachineSpec><DataCenter by="permalink">london</DataCenter><Resources><Package by="permalink">rock-3</Package></Resources><DiskTemplate><DiskTemplate by="permalink">templates/ubuntu-18-04</DiskTemplate></DiskTemplate><Hostname><Hostname>web-3</Hostname></Hostname></VirtualMachineSpec>
}

func ExampleVirtualMachineSpec_YAML() {
	spec := &VirtualMachineSpec{
		DataCenter: &DataCenter{Permalink: "london"},
		Resources: &Resources{
			Package: &Package{Permalink: "rock-3"},
		},
		DiskTemplate: &DiskTemplate{Permalink: "templates/ubuntu-18-04"},
		Hostname:     "web-3",
	}
	x, _ := spec.YAML()

	fmt.Println(string(x))
	// Output:
	// data_center:
	//   permalink: london
	// resources:
	//   package:
	//     permalink: rock-3
	// disk_template:
	//   permalink: templates/ubuntu-18-04
	// hostname: web-3
}
