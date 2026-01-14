package types

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

// CLIOptions contains options for creating a new command for the CLI.
type CLIOptions struct {
	// Use is the name of the command
	Use string

	// Short is the short description of the command
	Short string

	// Long is the long description of the command
	Long string
}
