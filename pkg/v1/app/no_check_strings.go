//go:build !check_missing_strings
// +build !check_missing_strings

package app

// checkMissingStrings is a no-op when the "check_missing_strings" build tag is not set.
func (a *App) checkMissingStrings() {
	// No-op
}
