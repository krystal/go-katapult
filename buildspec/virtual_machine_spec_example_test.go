package buildspec_test

import (
	"fmt"
	"github.com/krystal/go-katapult/buildspec"
)

//nolint:lll
func ExampleVirtualMachineSpec_JSON() {
	spec := &buildspec.VirtualMachineSpec{
		DataCenter: &buildspec.DataCenter{Permalink: "london"},
		Resources: &buildspec.Resources{
			Package: &buildspec.Package{Permalink: "rock-3"},
		},
		DiskTemplate: &buildspec.DiskTemplate{
			Permalink: "templates/ubuntu-18-04",
		},
		Hostname: "web-3",
	}
	out, _ := spec.JSON()

	fmt.Println(string(out))
	// Output:
	// {"data_center":{"permalink":"london"},"resources":{"package":{"permalink":"rock-3"}},"disk_template":{"permalink":"templates/ubuntu-18-04"},"hostname":"web-3"}
}

//nolint:lll
func ExampleVirtualMachineSpec_JSONIndent() {
	spec := &buildspec.VirtualMachineSpec{
		DataCenter: &buildspec.DataCenter{Permalink: "london"},
		Resources: &buildspec.Resources{
			Package: &buildspec.Package{Permalink: "rock-3"},
		},
		DiskTemplate: &buildspec.DiskTemplate{
			Permalink: "templates/ubuntu-18-04",
		},
		Hostname: "web-3",
	}
	out, _ := spec.JSONIndent("      ", "  ")

	fmt.Println("Spec: " + string(out))
	// Output:
	// Spec: {
	//         "data_center": {
	//           "permalink": "london"
	//         },
	//         "resources": {
	//           "package": {
	//             "permalink": "rock-3"
	//           }
	//         },
	//         "disk_template": {
	//           "permalink": "templates/ubuntu-18-04"
	//         },
	//         "hostname": "web-3"
	//       }
}

//nolint:lll
func ExampleVirtualMachineSpec_XML() {
	spec := &buildspec.VirtualMachineSpec{
		DataCenter: &buildspec.DataCenter{Permalink: "london"},
		Resources: &buildspec.Resources{
			Package: &buildspec.Package{Permalink: "rock-3"},
		},
		DiskTemplate: &buildspec.DiskTemplate{
			Permalink: "templates/ubuntu-18-04",
		},
		Hostname: "web-3",
	}
	out, _ := spec.XML()

	fmt.Println(string(out))
	// Output:
	// <?xml version="1.0" encoding="UTF-8"?>
	// <VirtualMachineSpec><DataCenter by="permalink">london</DataCenter><Resources><Package by="permalink">rock-3</Package></Resources><DiskTemplate><DiskTemplate by="permalink">templates/ubuntu-18-04</DiskTemplate></DiskTemplate><Hostname><Hostname>web-3</Hostname></Hostname></VirtualMachineSpec>
}

//nolint:lll
func ExampleVirtualMachineSpec_XMLIndent() {
	spec := &buildspec.VirtualMachineSpec{
		DataCenter: &buildspec.DataCenter{Permalink: "london"},
		Resources: &buildspec.Resources{
			Package: &buildspec.Package{Permalink: "rock-3"},
		},
		DiskTemplate: &buildspec.DiskTemplate{
			Permalink: "templates/ubuntu-18-04",
		},
		Hostname: "web-3",
	}
	out, _ := spec.XMLIndent("      ", "  ")

	fmt.Println("Spec:\n" + string(out))
	// Output:
	// Spec:
	//       <VirtualMachineSpec>
	//         <DataCenter by="permalink">london</DataCenter>
	//         <Resources>
	//           <Package by="permalink">rock-3</Package>
	//         </Resources>
	//         <DiskTemplate>
	//           <DiskTemplate by="permalink">templates/ubuntu-18-04</DiskTemplate>
	//         </DiskTemplate>
	//         <Hostname>
	//           <Hostname>web-3</Hostname>
	//         </Hostname>
	//       </VirtualMachineSpec>
}

func ExampleVirtualMachineSpec_YAML() {
	spec := &buildspec.VirtualMachineSpec{
		DataCenter: &buildspec.DataCenter{Permalink: "london"},
		Resources: &buildspec.Resources{
			Package: &buildspec.Package{Permalink: "rock-3"},
		},
		DiskTemplate: &buildspec.DiskTemplate{
			Permalink: "templates/ubuntu-18-04",
		},
		Hostname: "web-3",
	}
	out, _ := spec.YAML()

	fmt.Println(string(out))
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
