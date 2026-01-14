package net

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

import (
	"bufio"
	"encoding/hex"
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

		// Determine the status of the network interface by unpacking the
		// flags, looking for the 'up' flag. Default is unknown, if not up
		// then it's down (as per the net package)
		var _status types.NetworkInterfaceStatus
		_status = types.NetworkInterfaceStatusUnknown
		if iface.Flags&net.FlagUp != 0 {
			_status = types.NetworkInterfaceStatusUp
		} else {
			_status = types.NetworkInterfaceStatusDown
		}

		interfaceInfo := types.NetworkInterfaceInfo{
			Name:         iface.Name,
			HardwareAddr: iface.HardwareAddr.String(),
			IPAddresses:  ipAddresses,

			// Following fields are determined by the flags which are
			// unpacked from the net.Interface type for better readability
			// and simplicity in condition checking
			Status:            _status,
			Running:           iface.Flags&net.FlagRunning != 0,
			SupportsBroadcast: iface.Flags&net.FlagBroadcast != 0,
			SupportsMulticast: iface.Flags&net.FlagMulticast != 0,
			IsLoopback:        iface.Flags&net.FlagLoopback != 0,
			IsP2P:             iface.Flags&net.FlagPointToPoint != 0,
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
// Notes:
//
// this function is a simple wrapper around ResolveIPAddress().
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
	// Skip header line
	if scanner.Scan() {
		for scanner.Scan() {
			fields := strings.Fields(scanner.Text())
			// Ensure the line has enough fields and the MAC address is valid
			if len(fields) >= 4 && fields[2] != "00:00:00:00:00:00" {
				arpTable[fields[0]] = fields[2]
			}
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
	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", address, 2*1000*1000*1000)
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
func GetDefaultGateway() (string, error) {
	file, err := os.Open("/proc/net/route")
	if err != nil {
		return "", fmt.Errorf("error opening /proc/net/route: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// Skip header line
	if scanner.Scan() {
		for scanner.Scan() {
			fields := strings.Fields(scanner.Text())
			// Check if the line has enough fields
			// Destination (field 1) should be 00000000
			// Flags (field 3) should include RTF_GATEWAY (0x2)
			if len(fields) >= 4 && fields[1] == "00000000" && (parseHex(fields[3])&0x2) != 0 {
				gatewayHex := fields[2] // Gateway is field 2
				// The gateway address is in little-endian hexadecimal format
				ipBytes, err := hex.DecodeString(gatewayHex)
				if err != nil || len(ipBytes) != 4 {
					fmt.Fprintf(os.Stderr, "Warning: skipping malformed gateway hex '%s' in /proc/net/route\n", gatewayHex)
					continue
				}
				// Reverse bytes for standard big-endian representation
				ip := net.IPv4(ipBytes[3], ipBytes[2], ipBytes[1], ipBytes[0])
				// Ensure the parsed IP is not unspecified (0.0.0.0)
				if !ip.IsUnspecified() {
					return ip.String(), nil
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading /proc/net/route: %v", err)
	}

	// It's also possible the container genuinely has no default route configured.
	return "", errors.New("default gateway not found in /proc/net/route")
}

// parseHex parses a hex string and returns its integer value.
// Returns 0 on error.
func parseHex(s string) int {
	// Use ParseUint as flags are unsigned. Base 16, 64-bit size.
	val, err := strconv.ParseUint(s, 16, 64)
	if err != nil {
		return 0 // Return 0 if parsing fails
	}
	return int(val)
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
		isIPv6 := strings.HasSuffix(filePath, "6")
		file, err := os.Open(filePath)
		if err != nil {
			// Don't fail if one file is missing (e.g., IPv6 disabled)
			if os.IsNotExist(err) {
				continue
			}
			return nil, fmt.Errorf("error opening %s: %v", filePath, err)
		}
		defer file.Close() // Ensure closure in this scope

		connInfosFromFile, err := parseConnections(file, isIPv6)
		if err != nil {
			return nil, fmt.Errorf("error parsing %s: %v", filePath, err)
		}

		connInfos = append(connInfos, connInfosFromFile...)
	}

	return connInfos, nil
}

// parseConnections parses the output of /proc/net/tcp and /proc/net/tcp6
// to obtain information about active connections.
func parseConnections(reader io.Reader, isIPv6 bool) ([]types.ConnectionInfo, error) {
	connInfos := make([]types.ConnectionInfo, 0)

	scanner := bufio.NewScanner(reader)
	// Skip header line
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("error reading header: %v", err)
		}
		return connInfos, nil // Empty file is valid
	}

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		if len(fields) >= 4 { // Need at least sl, local_address, rem_address, st
			localAddr, remoteAddr, state, err := parseConnectionFields(fields, isIPv6)
			if err != nil {
				// Log or handle parsing errors for individual lines?
				// For now, continue to the next line.
				continue
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
func parseConnectionFields(fields []string, isIPv6 bool) (string, string, string, error) {
	if len(fields) < 4 {
		return "", "", "", fmt.Errorf("insufficient fields for parsing connection line")
	}
	localAddrField, remoteAddrField, state := fields[1], fields[2], fields[3]

	localIP, localPort, err := parseAddress(localAddrField, isIPv6)
	if err != nil {
		return "", "", "", fmt.Errorf("error parsing local address '%s': %v", localAddrField, err)
	}

	remoteIP, remotePort, err := parseAddress(remoteAddrField, isIPv6)
	if err != nil {
		return "", "", "", fmt.Errorf("error parsing remote address '%s': %v", remoteAddrField, err)
	}

	// Format address string based on IP version
	localAddrStr := net.JoinHostPort(localIP, localPort)
	remoteAddrStr := net.JoinHostPort(remoteIP, remotePort)

	return localAddrStr, remoteAddrStr, state, nil
}

// parseAddress parses an address in hexadecimal format from /proc/net/tcp*
// and returns the corresponding IP address and port string.
// Handles both IPv4 (8 hex chars) and IPv6 (32 hex chars).
func parseAddress(address string, isIPv6 bool) (string, string, error) {
	parts := strings.Split(address, ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid address format: %s", address)
	}
	hexIP, hexPort := parts[0], parts[1]

	// Parse Port (same for IPv4 and IPv6)
	portInt, err := strconv.ParseUint(hexPort, 16, 16) // Port is 16 bits
	if err != nil {
		return "", "", fmt.Errorf("error parsing hexadecimal port '%s': %v", hexPort, err)
	}
	portStr := strconv.FormatUint(portInt, 10)

	// Parse IP
	ipBytes, err := hex.DecodeString(hexIP)
	if err != nil {
		return "", "", fmt.Errorf("error decoding hexadecimal IP '%s': %v", hexIP, err)
	}

	var ip net.IP
	if isIPv6 {
		if len(ipBytes) != net.IPv6len {
			return "", "", fmt.Errorf("invalid IPv6 hex length for '%s'", hexIP)
		}
		// IPv6 is stored in network byte order (big-endian) in /proc/net/tcp6
		ip = net.IP(ipBytes)
	} else {
		if len(ipBytes) != net.IPv4len {
			return "", "", fmt.Errorf("invalid IPv4 hex length for '%s'", hexIP)
		}
		// IPv4 is stored in host byte order (little-endian) in /proc/net/tcp
		// Reverse bytes for standard big-endian representation
		ip = net.IPv4(ipBytes[3], ipBytes[2], ipBytes[1], ipBytes[0])
	}

	return ip.String(), portStr, nil
}
