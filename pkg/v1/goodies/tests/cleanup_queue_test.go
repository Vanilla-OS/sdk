package tests

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

import (
	"os/exec"
	"testing"

	"github.com/vanilla-os/sdk/pkg/v1/goodies"
)

func TestCleanupQueue(t *testing.T) {
	cleanupQueue := goodies.NewCleanupQueue()

	// Mock task execution
	task1Executed := false
	task1 := func(args ...interface{}) error {
		task1Executed = true
		return nil
	}
	args1 := []interface{}{"arg1", "arg2"}

	// Mock task execution with error
	task2Executed := false
	task2 := func(args ...interface{}) error {
		task2Executed = true
		return exec.Command("invalidcommand").Run()
	}
	args2 := []interface{}{"arg3"}

	// Mock task execution with error (2)
	task3Executed := false
	task3 := func(args ...interface{}) error {
		task3Executed = true
		return exec.Command("invalidcommand").Run()
	}
	args3 := []interface{}{"arg4"}

	// Mock error handler
	errorHandlerExecuted := false
	errorHandler := goodies.ErrorHandlerFn(func(args ...interface{}) error {
		errorHandlerExecuted = true
		return nil
	})

	// Mock error handler with error
	errorHandlerWithErrorExecuted := false
	errorHandlerWithError := goodies.ErrorHandlerFn(func(args ...interface{}) error {
		errorHandlerWithErrorExecuted = true
		return exec.Command("invalidcommand").Run()
	})

	cleanupQueue.Add(task1, args1, 1, &goodies.NoErrorHandler{}, false)
	cleanupQueue.Add(task2, args2, 2, errorHandler, false)
	cleanupQueue.Add(task3, args3, 3, errorHandlerWithError, false)

	err := cleanupQueue.Run()

	if !task1Executed {
		t.Error("Task 1 should have been executed")
	} else {
		t.Log("Task 1 was executed")
	}

	if !task2Executed {
		t.Error("Task 2 should have been executed")
	} else {
		t.Log("Task 2 was executed")
	}

	if !task3Executed {
		t.Error("Task 3 should have been executed")
	} else {
		t.Log("Task 3 was executed")
	}

	if !errorHandlerExecuted {
		t.Error("Error handler should have been executed")
	} else {
		t.Log("Error handler was executed")
	}

	if !errorHandlerWithErrorExecuted {
		t.Error("Error handler with error should have been executed")
	} else {
		t.Log("Error handler with error was executed")
	}

	if err == nil {
		t.Error("Error should not be nil, last task error handler should have failed")
	} else {
		t.Logf("Cleanup queue returned error as expected: %v", err)
	}
}
