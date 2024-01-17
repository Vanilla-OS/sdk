package types

// Process represents information about a process
type Process struct {
	// PID (Process ID) is the unique identifier for the process
	PID int `json:"pid"`

	// Name is the name of the process
	Name string `json:"name"`

	// State is the current state of the process
	State string `json:"state"`

	// PPID (Parent Process ID) is the PID of the parent process
	PPID int `json:"ppid"`

	// Priority is the priority of the process
	Priority int `json:"priority"`

	// Nice (user-space priority) is the nice value of the process
	Nice int `json:"nice"`

	// Threads is the number of threads in the process
	Threads int `json:"threads"`

	// UID (User ID) is the user ID of the process owner
	UID int `json:"uid"`

	// GID (Group ID) is the group ID of the process owner
	GID int `json:"gid"`
}
