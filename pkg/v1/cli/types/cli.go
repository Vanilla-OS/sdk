package types

// CLIOptions contains options for creating a new command for the CLI.
type CLIOptions struct {
	// Use is the name of the command
	Use string

	// Short is the short description of the command
	Short string

	// Long is the long description of the command
	Long string
}
