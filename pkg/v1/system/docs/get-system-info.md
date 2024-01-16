# Getting system info

To get system information, use the `GetSystemInfo` function from the
`v1/system` package.

The following example shows how to get system information:

```go
package main

import (
	"fmt"
	"github.com/vanilla-os/sdk/pkg/v1/system"
)

func main() {
    // Call the GetSystemInfo function to retrieve system information
    systemInfo, err := system.GetSystemInfo()
    if err != nil {
        fmt.Printf("Error getting system information: %v\n", err)
        return
    }

    // Access and print the MachineType from the retrieved system information
    fmt.Printf("Operating System: %s\n", systemInfo.OS)
    fmt.Printf("Version: %s\n", systemInfo.Version)
    fmt.Printf("Codename: %s\n", systemInfo.Codename)
    fmt.Printf("Arch: %s\n", systemInfo.Arch)
    fmt.Printf("Machine Type: %s\n", systemInfo.MachineType)
}
```

The output of the above example is similar to the following:

```text
Operating System: Vanilla OS
Version: 2
Codename: focalorchid
Arch: x86_64
Machine Type: container
```

The machine type can be one of the following:

- `container`: The application is running in a container.
- `vm`: The application is running in a virtual machine.
- `baremetal`: The application is running on bare metal.

the `baremetal` is assumed by default if none of the above is detected.
