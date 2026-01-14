package chroot

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

import (
	"os"
	"path/filepath"
	"syscall"
)

// EnterPivot enters a pivot and changes the current working directory
// to the new root. It returns a function that exits the pivot and
// changes the current working directory to the old root. Remember to
// call this function after you are done with the pivot to avoid
// leaving the process in a pivot.
//
// Example:
//
//	fExit, err := chroot.EnterPivot("/path/to/new/root")
//	if err != nil {
//		fmt.Printf("Error entering pivot root: %v\n", err)
//		if fExit != nil {
//			fExit()
//		}
//		return
//	}
//	defer fExit()
func EnterPivot(rootFs string) (f func() error, err error) {
	fd, err := os.Open("/")
	if err != nil {
		return nil, err
	}

	closeFunc := func() error {
		defer fd.Close()
		if err := fd.Chdir(); err != nil {
			return err
		}
		return syscall.PivotRoot(".", rootFs)
	}

	pivotDir := filepath.Join(rootFs, ".pivot_root")
	err = os.MkdirAll(pivotDir, 0755)
	if err != nil {
		return closeFunc, err
	}

	err = syscall.PivotRoot(rootFs, pivotDir)
	if err != nil {
		return closeFunc, err
	}

	err = os.Chdir("/")
	if err != nil {
		return closeFunc, err
	}

	return closeFunc, nil
}

// RunPivot runs a function in a pivot and changes the current working
// environment to the new root. It exits the pivot after the function
// returns. The function returns the error returned by the function and
// the error returned by ExitPivot if any.
//
// Example:
//
//	err := chroot.RunPivot("/path/to/new/root", func() error {
//		return exec.Command("ls", "-l").Run()
//	})
func RunPivot(rootFs string, f func() error) error {
	fExit, err := EnterPivot(rootFs)
	if err != nil {
		return err
	}
	defer fExit()

	return f()
}

// ExitPivot exits the pivot and changes the current working directory
// to the old root. It returns an error if any.
//
// Example:
//
//	err := chroot.ExitPivot()
func ExitPivot() error {
	return syscall.PivotRoot(".", ".")
}
