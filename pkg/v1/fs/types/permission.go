package types

import "fmt"

// Permission represents file permissions of a file
type Permission struct {
	OwnerRead    bool
	OwnerWrite   bool
	OwnerExecute bool

	GroupRead    bool
	GroupWrite   bool
	GroupExecute bool

	OthersRead    bool
	OthersWrite   bool
	OthersExecute bool
}

func (p Permission) String() string {
	boolToChar := func(b bool) string {
		if b {
			return "r"
		}
		return "-"
	}
	return fmt.Sprintf("%s%s%s%s%s%s%s%s%s",
		boolToChar(p.OwnerRead), boolToChar(p.OwnerWrite), boolToChar(p.OwnerExecute),
		boolToChar(p.GroupRead), boolToChar(p.GroupWrite), boolToChar(p.GroupExecute),
		boolToChar(p.OthersRead), boolToChar(p.OthersWrite), boolToChar(p.OthersExecute),
	)
}
