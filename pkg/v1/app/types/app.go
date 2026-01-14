package types

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

import (
	"io/fs"

	cliTypes "github.com/vanilla-os/sdk/pkg/v1/cli/types"
)

// Sign is a unique signature for the application
type Sign string

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
	LocalesFS fs.FS

	// DefaultLocale is the default locale for the application, this should
	// always be empty unless you want to force a specific locale for the
	// application, for example for testing purposes.
	DefaultLocale string

	// CLIOptions contains options for creating the command line interface
	CLIOptions *cliTypes.CLIOptions
}
