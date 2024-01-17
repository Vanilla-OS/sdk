# Getting Started

This guide will walk you through the steps to get started with the SDK in a
Go application. In this example, we'll demonstrate how to use the
`GetSystemInfo` function from the `v1/system` package to retrieve system
information.

## Prerequisites

The SDK is tested on Vanilla OS 2 and above. Ensure that you have a compatible
environment before proceeding. For your convenience, we have provided a Docker
image of Vanilla OS 2 with Go installed along with other useful tools and
libraries. You can use easily use this image in Apx, using the `vanilla-dev`
stack, by issuing the following command:

```bash
apx subsystems new --name sdk-demo --stack vanilla-dev --init
```

the `--init` flag is optional, it will initialize systemd in the container,
allowing you to test SDK features that require systemd.

If you prefer to use your own environment, ensure that you have Go installed.
Please note that the SDK is meant to be used in Vanilla OS only, and may not
work as expected in other environments, do not expect support for issues
encountered in other operating systems.

## Create a Go Application

Start by creating a new Go application. Open your favorite text editor or
integrated development environment (IDE) and create a new file named `main.go`
in an empty directory. Add the following code:

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
	fmt.Printf("Machine Type: %s\n", systemInfo.MachineType)
}
```

## Importing the SDK

In the `main.go` file, we import the `system` package from the Vanilla OS SDK
since we aim to retrieve system information.

```go
import (
	"fmt"
	"github.com/vanilla-os/sdk/pkg/v1/system"
)
```

If you need to use other packages from the SDK, you can import them in the same
way, by specifying the package path.

## Using GetSystemInfo

Call the `GetSystemInfo` function from the `system` package to retrieve system
information.

```go
systemInfo, err := system.GetSystemInfo()
if err != nil {
	fmt.Printf("Error getting system information: %v\n", err)
	return
}
```

By convention, the majority of functions in the SDK return an error as the last
return value. Always handle the error returned by the function, even if it is
unlikely to occur.

## Accessing System Information

The `GetSystemInfo` function returns a `SystemInfo` struct, which contains
information about the system. You can access the information by accessing the
struct fields.

```go
fmt.Printf("Operating System: %s\n", systemInfo.OS)
fmt.Printf("Version: %s\n", systemInfo.Version)
fmt.Printf("Machine Type: %s\n", systemInfo.MachineType)
```

## Running the Application

Save the `main.go` file and run the application using the following command:

```bash
go run main.go
```

You should see the system information printed to the console. If you are
performing this tutorial in Apx, you will notice that the Machine Type is
`container`, this is because the application is running in a Docker container.

Congratulations! You have successfully integrated the Vanilla OS SDK into
your Go application and used the `system` package to retrieve and display
system information. Feel free to explore other packages and functionalities
provided by the SDK for a more in-depth integration into your projects; each
public function in the SDK is documented and offers a detailed description of
its purpose along with example usage.
