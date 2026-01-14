//go:build !check_missing_strings
// +build !check_missing_strings

package app

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

// checkMissingStrings is a no-op when the "check_missing_strings" build tag is not set.
func (a *App) checkMissingStrings() {
	// No-op
}
