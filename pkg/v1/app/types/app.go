package types

import (
	"embed"

	logsTypes "github.com/vanilla-os/sdk/pkg/v1/logs/types"
	"github.com/vorlif/spreak"
)

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

	// LocalesFS is the file system containing the locales for the application
	LocalesFS embed.FS

	// DefaultLocale is the default locale for the application, this should
	// always be empty unless you want to force a specific locale for the
	// application, for example for testing purposes.
	DefaultLocale string
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

	// Logger is the logger for the application
	Logger logsTypes.Logger

	// LC (Localizer) is the localizer for the application
	LC spreak.Localizer

	// LocalesFS is the file system containing the locales for the application
	LocalesFS embed.FS
}
