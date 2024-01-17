package tests

import (
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vanilla-os/sdk/pkg/v1/system"
)

func TestGetProcessInfo(t *testing.T) {
	process, err := system.GetProcessInfo(1)
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	t.Logf("PID: %d", process.PID)
	t.Logf("Name: %s", process.Name)
	t.Logf("State: %s", process.State)
	t.Logf("PPID: %d", process.PPID)
	t.Logf("Priority: %d", process.Priority)
	t.Logf("Nice: %d", process.Nice)
	t.Logf("Threads: %d", process.Threads)
	t.Logf("UID: %d", process.UID)
	t.Logf("GID: %d", process.GID)
}

func TestGetProcessList(t *testing.T) {
	processes, err := system.GetProcessList()
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	for i, process := range processes {
		t.Logf("PID: %d", process.PID)
		t.Logf("Name: %s", process.Name)
		t.Logf("State: %s", process.State)
		t.Logf("PPID: %d", process.PPID)
		t.Logf("Priority: %d", process.Priority)
		t.Logf("Nice: %d", process.Nice)
		t.Logf("Threads: %d", process.Threads)
		t.Logf("UID: %d", process.UID)
		t.Logf("GID: %d", process.GID)

		if i == 5 {
			return
		}
	}
}

func TestKillProcess(t *testing.T) {
	cmd := exec.Command("sleep", "10")
	err := cmd.Start()
	assert.NoError(t, err, "failed to start the sleep process")

	pid := cmd.Process.Pid
	t.Logf("Test process started with PID %d", pid)

	done := make(chan error)
	go func() {
		done <- system.KillProcess(pid)
	}()

	time.Sleep(1 * time.Second)

	select {
	case err := <-done:
		assert.NoError(t, err, "error killing the process")
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for process termination")
	}

	err = cmd.Wait()
	assert.Error(t, err, "expected an error indicating the process is killed")
}
