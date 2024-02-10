package types

// ConfigOptions is a struct that holds the configuration options
type ConfigOptions struct {
	// Domain is the domain of the configuration file, in the context of a
	// Vanilla OS application, it is the RDNN
	Domain string

	// Path is the path to the configuration file
	Path string

	// Type is the type of the configuration file, e.g. json, yaml
	Type string

	// Prefix is an optional prefix for the Path, for testing purposes
	Prefix string
}
