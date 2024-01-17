package types

// Peripheral represents a system peripheral.
type Peripheral struct {
	ID   string     `json:"id"`
	Name string     `json:"name"`
	Type DeviceType `json:"type"`
}
