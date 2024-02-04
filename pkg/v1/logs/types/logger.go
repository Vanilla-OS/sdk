package types

import "github.com/phuslu/log"

// Logger represents a logger for the application
type Logger struct {
	// Console is the logger used to log messages to the console, use this
	// for any logging the user should see
	Console log.Logger

	// File is the logger used to log messages to the vlogs directory, use
	// this for any internal logging the user doesn't need to see
	File log.Logger
}
