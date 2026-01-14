package system

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

import (
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/vanilla-os/sdk/pkg/v1/system/types"
)

// GetAllUsers retrieves information about all users on the system by
// parsing /etc/passwd. If includeNoLogin is true, users with /usr/sbin/nologin
// as their shell will be included.
//
// Example:
//
//	users, err := system.GetAllUsers(false)
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
//
//	for _, user := range users {
//		fmt.Printf("UID: %s\n", user.UID)
//		fmt.Printf("Username: %s\n", user.Username)
//	}
func GetAllUsers(includeNoLogin bool) ([]types.UserInfo, error) {
	etcPasswd, err := os.ReadFile("/etc/passwd")
	if err != nil {
		return nil, fmt.Errorf("error reading /etc/passwd: %v", err)
	}

	var users []types.UserInfo

	for _, line := range strings.Split(string(etcPasswd), "\n") {
		if line == "" {
			continue
		}

		items := strings.Split(line, ":")
		if len(items) < 7 {
			continue
		}

		username := items[0]
		hasLogin := items[6] != "/usr/sbin/nologin"
		shell := items[6]

		if username == "root" || username == "sync" {
			continue
		}

		if !includeNoLogin && !hasLogin {
			continue
		}

		u, err := user.Lookup(username)
		if err != nil {
			return nil, fmt.Errorf("error looking up user %s: %v", u, err)
		}

		users = append(users, types.UserInfo{
			UID:      u.Uid,
			GID:      u.Gid,
			Username: u.Username,
			Name:     u.Name,
			HomeDir:  u.HomeDir,
			Shell:    shell,
		})
	}

	return users, nil
}

// GetUsers retrieves users with the given usernames and UIDs. If useUID is
// true, the returned map will use the UID as the key, otherwise it will use
// the username as the key.
//
// Example:
//
//	users, err := system.GetUsers([]string{"john", "jane"}, []string{"1000"})
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
//
//	for _, user := range users {
//		fmt.Printf("UID: %s\n", user.UID)
//		fmt.Printf("Username: %s\n", user.Username)
//	}
//
// Notes:
//
// if a username and UID refer to the same user, the user will only be
// returned once. If an user does not match any of the given usernames or UIDs,
// it will not be available in the returned map.
func GetUsers(usernames []string, uids []string, useUID bool) map[string]types.UserInfo {
	usersMap := make(map[string]types.UserInfo)

	for _, username := range usernames {
		u, err := user.Lookup(username)
		if err != nil {
			// If the user doesn't exist, continue
			continue
		}

		key := u.Username
		if useUID {
			key = u.Uid
		}

		usersMap[key] = types.UserInfo{
			UID:      u.Uid,
			GID:      u.Gid,
			Username: u.Username,
			Name:     u.Name,
			HomeDir:  u.HomeDir,
			Shell:    u.HomeDir,
		}
	}

	for _, uid := range uids {
		u, err := user.LookupId(uid)
		if err != nil {
			// If the user doesn't exist, continue
			continue
		}

		// Check if the user already exists from the previous loop, since the
		// developer may have passed in a username and UID that refer to the
		// same user
		exists := false
		for _, user := range usersMap {
			if useUID && user.UID == u.Uid {
				exists = true
				break
			} else if !useUID && user.Username == u.Username {
				exists = true
				break
			}
		}

		if exists {
			continue
		}

		key := uid
		if !useUID {
			key = u.Username
		}

		usersMap[key] = types.UserInfo{
			UID:      u.Uid,
			GID:      u.Gid,
			Username: u.Username,
			Name:     u.Name,
			HomeDir:  u.HomeDir,
			Shell:    u.HomeDir,
		}
	}

	return usersMap
}

// GetAllGroups retrieves information about all groups on the system by
// parsing /etc/group.
//
// Example:
//
//	groups, err := system.GetGroups()
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
//
//	for _, group := range groups {
//		fmt.Printf("GID: %s\n", group.GID)
//		fmt.Printf("Name: %s\n", group.Name)
//	}
func GetAllGroups() ([]types.GroupInfo, error) {
	etcGroup, err := os.ReadFile("/etc/group")
	if err != nil {
		return nil, fmt.Errorf("error reading /etc/group: %v", err)
	}

	var groups []types.GroupInfo

	for _, line := range strings.Split(string(etcGroup), "\n") {
		if line == "" {
			continue
		}

		items := strings.Split(line, ":")
		if len(items) < 3 {
			continue
		}

		gid := items[2]
		name := items[0]

		if name == "root" {
			continue
		}

		groups = append(groups, types.GroupInfo{
			GID:  gid,
			Name: name,
		})
	}

	return groups, nil
}
