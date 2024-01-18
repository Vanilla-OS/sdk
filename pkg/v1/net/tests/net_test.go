package tests

import (
	"testing"

	"github.com/vanilla-os/sdk/pkg/v1/net"
)

const testHostname = "vanillaos.org"
const nonExistentHostname = "iambatman.nanana"

func TestGetNetworkInterfaces(t *testing.T) {
	interfaces, err := net.GetNetworkInterfaces()
	if err != nil {
		t.Errorf("Error getting network interfaces: %v", err)
	}

	if len(interfaces) == 0 {
		t.Skip("No network interfaces found")
	}

	for _, iface := range interfaces {
		t.Logf("Name: %s", iface.Name)
		t.Logf("Hardware address: %s", iface.HardwareAddr)
		t.Logf("IP addresses: %v", iface.IPAddresses)
	}
}

func TestGetInterfaceIPAddresses(t *testing.T) {
	interfaces, _ := net.GetNetworkInterfaces()
	if len(interfaces) == 0 {
		t.Skip("No network interfaces found. Skipping test.")
	}

	ipAddresses, err := net.GetInterfaceIPAddresses(interfaces[0].Name)
	if err != nil {
		t.Errorf("Error getting IP addresses for interface: %v", err)
	}

	if len(ipAddresses) == 0 {
		t.Error("No IP addresses found for the interface")
	}
}

func TestResolveIPAddress(t *testing.T) {
	ipAddress, err := net.ResolveIPAddress(testHostname)
	if err != nil {
		t.Errorf("Error resolving IP address: %v", err)
	}

	// IP address length should be greater than 0
	if len(ipAddress) == 0 {
		t.Error("Invalid IP address")
	}

	// Testing with a nonexistent hostname
	_, err = net.ResolveIPAddress(nonExistentHostname)
	if err == nil {
		t.Errorf("Unexpected success resolving nonexistent hostname: %s", nonExistentHostname)
	}
}

func TestCheckInternetConnectivity(t *testing.T) {
	if !net.CheckInternetConnectivity() {
		t.Error("Internet connectivity check failed")
	}
}

func TestGetDNSInfo(t *testing.T) {
	dnsInfo, err := net.GetDNSInfo(testHostname)
	if err != nil {
		t.Errorf("Error getting DNS info: %v", err)
	}

	if dnsInfo.Hostname != testHostname {
		t.Skipf("Expected hostname %s, got %s", testHostname, dnsInfo.Hostname)
	}

	if len(dnsInfo.IPAddresses) == 0 {
		t.Error("No IP addresses found for the hostname")
	}
}

func TestIsLocalNetworkIP(t *testing.T) {
	interfaces, _ := net.GetNetworkInterfaces()
	if len(interfaces) == 0 {
		t.Skip("No network interfaces found. Skipping test.")
	}

	for _, iface := range interfaces {
		if len(iface.IPAddresses) > 0 {
			ipAddress := iface.IPAddresses[0]
			if !net.IsLocalNetworkIP(ipAddress) {
				t.Errorf("Expected IP address %s to belong to a local network", ipAddress)
			}
		}
	}

	// Testing with a non-local IP address
	if net.IsLocalNetworkIP("8.8.8.8") {
		t.Error("Unexpected success for a non-local IP address")
	}
}

func TestGetDefaultGateway(t *testing.T) {
	gateway, err := net.GetDefaultGateway()
	if err != nil {
		t.Errorf("Error getting default gateway: %v", err)
	}

	// Gateway IP address length should be greater than 0
	if len(gateway) == 0 {
		t.Error("Invalid default gateway IP address")
	}
}

func TestGetActiveConnections(t *testing.T) {
	connInfos, err := net.GetActiveConnections()
	if err != nil {
		t.Errorf("Error getting active connections: %v", err)
	}

	if len(connInfos) == 0 {
		t.Skip("No active connections found")
	}

	for i, connInfo := range connInfos {
		t.Logf("Local address: %s", connInfo.LocalAddr)
		t.Logf("Remote address: %s", connInfo.RemoteAddr)
		t.Logf("State: %s", connInfo.State)

		if i == 5 {
			break
		}
	}
}
