package tests

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

import (
	"testing"

	"github.com/vanilla-os/sdk/pkg/v1/system"
)

func TestGetAllusers(t *testing.T) {
	users, err := system.GetAllUsers(false)
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	for _, user := range users {
		t.Logf("UID: %s", user.UID)
		t.Logf("GID: %s", user.GID)
		t.Logf("Username: %s", user.Username)
		t.Logf("Name: %s", user.Name)
		t.Logf("HomeDir: %s", user.HomeDir)
		t.Logf("Shell: %s", user.Shell)
	}
}

func TestGetAllGroups(t *testing.T) {
	groups, err := system.GetAllGroups()
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	for _, group := range groups {
		t.Logf("GID: %s", group.GID)
		t.Logf("Name: %s", group.Name)
	}
}

func TestGetUsers(t *testing.T) {
	users := system.GetUsers([]string{}, []string{"1000", "1001"}, false)

	if len(users) == 0 {
		t.Skip("No users found")
		return
	}

	for key, value := range users {
		t.Logf("Key: %s", key)
		t.Logf("Info:")
		t.Logf("\tUID: %s", value.UID)
		t.Logf("\tGID: %s", value.GID)
		t.Logf("\tUsername: %s", value.Username)
	}

}
