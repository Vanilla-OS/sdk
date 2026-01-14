package types

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

// UserInfo represents information about a user
type UserInfo struct {
	// UID is the user's ID
	UID string `json:"uid"`

	// GID is the user's group ID
	GID string `json:"gid"`

	// Username is the user's username
	Username string `json:"username"`

	// Name is the user's name
	Name string `json:"name"`

	// HomeDir is the user's home directory
	HomeDir string `json:"home_directory"`

	// Shell is the user's configured shell
	Shell string `json:"shell"`
}
