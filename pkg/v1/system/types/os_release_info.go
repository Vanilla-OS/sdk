package types

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

// OSReleaseInfo is a struct that contains information about the OS release
type OSReleaseInfo struct {
	// Name is the name of the operating system
	Name string

	// Version is the version of the operating system
	Version string

	// Codename is how the operating system is referred to internally
	Codename string
}
