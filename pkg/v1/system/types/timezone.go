package types

// Timezone represents a supported timezone
type Timezone struct {
	// Name is the name of the timezone
	Name string `json:"name"`

	// Location is the location of the timezone file
	Location string `json:"location"`
}
