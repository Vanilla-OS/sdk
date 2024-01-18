package types

// DNSInfo represents DNS (Domain Name System) information
type DNSInfo struct {
	Hostname    string   `json:"hostname"`
	IPAddresses []string `json:"ip_addresses"`
}
