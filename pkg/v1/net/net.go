package net

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/vanilla-os/sdk/pkg/v1/net/types"
)

// GetNetworkInterfaces returns a list of network interfaces available on the
// system.
//
// Example:
//
//	interfaces, err := net.GetNetworkInterfaces()
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
//
//	for _, iface := range interfaces {
//		fmt.Printf("Name: %s\n", iface.Name)
//		fmt.Printf("Hardware address: %s\n", iface.HardwareAddr)
//		fmt.Printf("IP addresses: %v\n", iface.IPAddresses)
//	}
func GetNetworkInterfaces() ([]types.NetworkInterfaceInfo, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("error getting network interfaces: %v", err)
	}

	// To keep consistency, we convert the net.Interface type to our own
	// NetworkInterfaceInfo type
	var interfaceInfoList []types.NetworkInterfaceInfo
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, fmt.Errorf("error getting addresses for interface %s: %v", iface.Name, err)
		}

		var ipAddresses []string
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if ok && !ipNet.IP.IsLoopback() {
				ipAddresses = append(ipAddresses, ipNet.IP.String())
			}
		}

		interfaceInfo := types.NetworkInterfaceInfo{
			Name:         iface.Name,
			HardwareAddr: iface.HardwareAddr.String(),
			IPAddresses:  ipAddresses,
		}

		interfaceInfoList = append(interfaceInfoList, interfaceInfo)
	}

	return interfaceInfoList, nil
}

// GetInterfaceIPAddresses returns the IP addresses associated with a
// specific network interface.
//
// Example:
//
//	ipAddresses, err := net.GetInterfaceIPAddresses("eth0")
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
//
//	for _, ip := range ipAddresses {
//		fmt.Printf("IP address: %s\n", ip)
//	}
func GetInterfaceIPAddresses(interfaceName string) ([]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range interfaces {
		if iface.Name == interfaceName {
			addresses, err := iface.Addrs()
			if err != nil {
				return nil, err
			}

			ipAddresses := make([]string, 0)
			for _, addr := range addresses {
				switch v := addr.(type) {
				case *net.IPNet:
					ipAddresses = append(ipAddresses, v.IP.String())
				case *net.IPAddr:
					ipAddresses = append(ipAddresses, v.IP.String())
				}
			}

			return ipAddresses, nil
		}
	}

	return nil, errors.New("network interface not found")
}

// ResolveIPAddress resolves the IP address of a given hostname.
//
// Example:
//
//	ipAddress, err := net.ResolveIPAddress("vanilla.org")
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
//
//	fmt.Printf("IP address: %s\n", ipAddress)
func ResolveIPAddress(hostname string) (string, error) {
	ipAddr, err := net.ResolveIPAddr("ip", hostname)
	if err != nil {
		return "", err
	}

	// Check if the IP address is a loopback address, if so it means that
	// the hostname does not resolve to a valid IP address
	if ipAddr.IP.IsLoopback() {
		return "", errors.New("hostname resolves to a loopback address")
	}

	return ipAddr.IP.String(), nil
}

// CheckInternetConnectivity checks if the system has internet connectivity
// by attempting to resolve the IP address of a well-known hostname
// (e.g., Google's public DNS).
//
// Example:
//
//	if net.CheckInternetConnectivity() {
//		fmt.Println("Internet connectivity is available")
//	} else {
//		fmt.Println("Internet connectivity is not available")
//	}
//
// Note: this function is a simple wrapper around ResolveIPAddress().
func CheckInternetConnectivity() bool {
	_, err := ResolveIPAddress("8.8.8.8")
	return err == nil
}

// ResolveMACAddress resolves the MAC address for a given IP address.
//
// Example:
//
//	macAddress, err := net.ResolveMACAddress("
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
//
//	fmt.Printf("MAC address: %s\n", macAddress)
func ResolveMACAddress(ipAddress string) (string, error) {
	ipAddr := net.ParseIP(ipAddress)
	if ipAddr == nil {
		return "", fmt.Errorf("invalid IP address: %s", ipAddress)
	}

	arpTable, err := readARPTable()
	if err != nil {
		return "", err
	}

	if macAddr, ok := arpTable[ipAddr.String()]; ok {
		return macAddr, nil
	}

	return "", fmt.Errorf("MAC address not found for IP: %s", ipAddress)
}

// readARPTable reads the ARP table to obtain the mapping of IP addresses
// to MAC addresses.
func readARPTable() (map[string]string, error) {
	arpTable := make(map[string]string)

	file, err := os.Open("/proc/net/arp")
	if err != nil {
		return nil, fmt.Errorf("error opening /proc/net/arp: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 4 && fields[2] != "00:00:00:00:00:00" {
			arpTable[fields[0]] = fields[2]
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading /proc/net/arp: %v", err)
	}

	return arpTable, nil
}

// GetDNSInfo resolves a hostname to its corresponding IP addresses.
//
// Example:
//
//	dnsInfo, err := net.GetDNSInfo("vanilla.org")
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
//
//	fmt.Printf("Hostname: %s\n", dnsInfo.Hostname)
//	fmt.Printf("IP addresses: %v\n", dnsInfo.IPAddresses)
func GetDNSInfo(hostname string) (*types.DNSInfo, error) {
	ips, err := net.LookupIP(hostname)
	if err != nil {
		return nil, fmt.Errorf("error resolving hostname %s: %v", hostname, err)
	}

	dnsInfo := &types.DNSInfo{
		Hostname:    hostname,
		IPAddresses: make([]string, len(ips)),
	}

	for i, ip := range ips {
		dnsInfo.IPAddresses[i] = ip.String()
	}

	return dnsInfo, nil
}

// IsLocalNetworkIP checks if the given IP address belongs to a local network.
//
// Example:
//
//	if net.IsLocalNetworkIP("
//		fmt.Println("IP address belongs to a local network")
//	} else {
//		fmt.Println("IP address does not belong to a local network")
//	}
func IsLocalNetworkIP(ipAddress string) bool {
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return false
	}

	// Check if the IP address belongs to a private network range
	return ip.IsLoopback() || ip.IsPrivate()
}

// CheckPortStatus checks the status of a given port on a specific host.
//
// Example:
//
//	import netTypes "github.com/vanilla-os/sdk/pkg/v1/net/types"
//
//	portStatus := net.CheckPortStatus("vanilla.org", 80)
//	if portStatus == netTypes.PortOpen {
//		fmt.Println("Port is open")
//	} else {
//		fmt.Println("Port is closed")
//	}
func CheckPortStatus(host string, port int) types.PortStatus {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return types.PortClosed
	}
	defer conn.Close()

	return types.PortOpen
}

// GetDefaultGateway retrieves the default gateway IP address.
//
// Example:
//
//	gateway, err := net.GetDefaultGateway()
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
//
//	fmt.Printf("Default gateway: %s\n", gateway)
//
// Note: this function relies on the 'ip' command
func GetDefaultGateway() (string, error) {
	routes, err := net.InterfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("error retrieving interface addresses: %v", err)
	}

	for _, route := range routes {
		ipNet, ok := route.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", errors.New("default gateway not found")
}

// GetActiveConnections returns information about active network connections.
//
// Example:
//
//	connInfos, err := net.GetActiveConnections()
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
//
//	for _, connInfo := range connInfos {
//		fmt.Printf("Local address: %s\n", connInfo.LocalAddr)
//		fmt.Printf("Remote address: %s\n", connInfo.RemoteAddr)
//		fmt.Printf("State: %s\n", connInfo.State)
//	}
func GetActiveConnections() ([]types.ConnectionInfo, error) {
	connInfos := make([]types.ConnectionInfo, 0)

	files := []string{"/proc/net/tcp", "/proc/net/tcp6"}

	for _, filePath := range files {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("error opening %s: %v", filePath, err)
		}
		defer file.Close()

		connInfosFromFile, err := parseConnections(file)
		if err != nil {
			return nil, fmt.Errorf("error parsing %s: %v", filePath, err)
		}

		connInfos = append(connInfos, connInfosFromFile...)
	}

	return connInfos, nil
}

// parseConnections parses the output of /proc/net/tcp and /proc/net/tcp6
// to obtain information about active connections.
func parseConnections(reader io.Reader) ([]types.ConnectionInfo, error) {
	connInfos := make([]types.ConnectionInfo, 0)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip header line
		if strings.HasPrefix(line, "  sl") {
			continue
		}

		fields := strings.Fields(line)

		if len(fields) >= 4 {
			localAddr, remoteAddr, state, err := parseConnectionFields(fields)
			if err != nil {
				return nil, err
			}

			var connState types.ConnState

			switch state {
			case "01":
				connState = types.ConnStateEnstablished
			case "02":
				connState = types.ConnStateSynSent
			case "03":
				connState = types.ConnStateSynRecv
			case "04":
				connState = types.ConnStateFinWait1
			case "05":
				connState = types.ConnStateFinWait2
			case "06":
				connState = types.ConnStateTimeWait
			case "07":
				connState = types.ConnStateClose
			case "08":
				connState = types.ConnStateCloseWait
			case "09":
				connState = types.ConnStateLastAck
			case "0A":
				connState = types.ConnStateListen
			case "0B":
				connState = types.ConnStateClosing
			case "0C":
				connState = types.ConnStateNewSynRecv
			default:
				connState = types.ConnStateUnknown
			}

			connInfo := types.ConnectionInfo{
				LocalAddr:  localAddr,
				RemoteAddr: remoteAddr,
				State:      connState,
			}
			connInfos = append(connInfos, connInfo)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning input: %v", err)
	}

	return connInfos, nil
}

// parseConnectionFields parses the fields of a line in /proc/net/tcp or
// /proc/net/tcp6.
func parseConnectionFields(fields []string) (string, string, string, error) {
	localAddr, remoteAddr, state := fields[1], fields[2], fields[3]

	localIP, localPort, err := parseAddress(localAddr)
	if err != nil {
		return "", "", "", fmt.Errorf("error parsing local address: %v", err)
	}

	remoteIP, remotePort, err := parseAddress(remoteAddr)
	if err != nil {
		return "", "", "", fmt.Errorf("error parsing remote address: %v", err)
	}

	return fmt.Sprintf("%s:%s", localIP, localPort), fmt.Sprintf("%s:%s", remoteIP, remotePort), state, nil
}

// parseAddress parses an address in hexadecimal format and returns the
// corresponding IP address and port.
func parseAddress(address string) (string, string, error) {
	parts := strings.Split(address, ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid address format: %s", address)
	}

	// We have to do some parsing here to get human-readable values
	ipInt, err := strconv.ParseInt(parts[0], 16, 64)
	if err != nil {
		return "", "", fmt.Errorf("error parsing hexadecimal IP address: %v", err)
	}
	ip := fmt.Sprintf("%d.%d.%d.%d", byte(ipInt>>24), byte(ipInt>>16), byte(ipInt>>8), byte(ipInt))

	portInt, err := strconv.ParseInt(parts[1], 16, 64)
	if err != nil {
		return "", "", fmt.Errorf("error parsing hexadecimal port: %v", err)
	}
	port := fmt.Sprintf("%d", portInt)

	return ip, port, nil
}
