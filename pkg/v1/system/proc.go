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
	"strconv"
	"strings"

	"github.com/vanilla-os/sdk/pkg/v1/system/types"
)

// GetProcessList returns a list of all processes running on the system.
//
// Example:
//
//	processes, err := system.GetProcessList()
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
//
//	for _, process := range processes {
//		fmt.Printf("PID: %d\n", process.PID)
//		fmt.Printf("Name: %s\n", process.Name)
//		fmt.Printf("State: %s\n", process.State)
//	}
func GetProcessList() ([]types.Process, error) {
	var processes []types.Process

	// Read the list of files in the /proc directory
	files, err := os.ReadDir("/proc")
	if err != nil {
		return nil, fmt.Errorf("error reading /proc directory: %v", err)
	}

	// Iterate over each file in /proc and check if it's a process
	for _, file := range files {
		pid, err := strconv.Atoi(file.Name())
		if err == nil {
			// If the file name is a number, it's a process
			process, err := GetProcessInfo(pid)
			if err == nil {
				processes = append(processes, *process)
			}
		}
	}

	return processes, nil
}

// GetProcessInfo returns information about a specific process.
//
// Example:
//
//	process, err := system.GetProcessInfo(1)
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
//
//	fmt.Printf("PID: %d\n", process.PID)
//	fmt.Printf("Name: %s\n", process.Name)
//	fmt.Printf("State: %s\n", process.State)
func GetProcessInfo(pid int) (*types.Process, error) {
	process := &types.Process{PID: pid}

	// Read information from the /proc/{pid}/stat file
	statPath := fmt.Sprintf("/proc/%d/stat", pid)
	statContent, err := os.ReadFile(statPath)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %v", statPath, err)
	}

	// Extract information from /proc/{pid}/stat content
	fields := strings.Fields(string(statContent))
	if len(fields) < 24 {
		return nil, fmt.Errorf("insufficient fields in %s", statPath)
	}

	process.Name = fields[1][1 : len(fields[1])-1]
	process.State = fields[2]
	process.PPID, _ = strconv.Atoi(fields[3])
	process.Priority, _ = strconv.Atoi(fields[17])
	process.Nice, _ = strconv.Atoi(fields[18])
	process.Threads, _ = strconv.Atoi(fields[19])

	// Read information from the /proc/{pid}/status file
	statusPath := fmt.Sprintf("/proc/%d/status", pid)
	statusContent, err := os.ReadFile(statusPath)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %v", statusPath, err)
	}

	// Extract information from /proc/{pid}/status content
	lines := strings.Split(string(statusContent), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			switch fields[0] {
			case "Uid:":
				process.UID, _ = strconv.Atoi(fields[1])
			case "Gid:":
				process.GID, _ = strconv.Atoi(fields[1])
			}
		}
	}

	return process, nil
}

// KillProcess terminates a process given its PID.
//
// Example:
//
//	err := system.KillProcess(2024)
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
func KillProcess(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("unable to find process with PID %d: %v", pid, err)
	}

	err = process.Kill()
	if err != nil {
		return fmt.Errorf("unable to kill process with PID %d: %v", pid, err)
	}

	return nil
}
