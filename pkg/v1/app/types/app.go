package types

// AppOptions contains options for creating a new Vanilla OS application
type AppOptions struct {
	// RDNN is the reverse domain name notation of the application, please
	// see <https://en.wikipedia.org/wiki/Reverse_domain_name_notation> for
	// more information.
	RDNN string

	// Name is the name of the application
	Name string

	// Version is the version of the application
	Version string
}

// App represents a Vanilla OS application
type App struct {
	// Sign is a unique signature for the application
	Sign string

	// RDNN is the reverse domain name notation of the application
	RDNN string

	// Name is the name of the application
	Name string

	// Version is the version of the application
	Version string
}
