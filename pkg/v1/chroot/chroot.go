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
	"syscall"
)

// EnterChroot enters a chroot and changes the current working directory
// to the new root. It returns a function that exits the chroot and
// changes the current working directory to the old root. Remember to
// call this function after you are done with the chroot to avoid
// leaving the process in a chroot.
//
// Example:
//
//	fExit, err := chroot.EnterChroot("/path/to/new/root")
//	if err != nil {
//		fmt.Printf("Error entering chroot root: %v\n", err)
//		if fExit != nil {
//			fExit()
//		}
//		return
//	}
//	defer fExit()
func EnterChroot(rootFs string) (f func() error, err error) {
	fd, err := os.Open("/")
	if err != nil {
		return nil, err
	}

	closeFunc := func() error {
		defer fd.Close()
		if err := fd.Chdir(); err != nil {
			return err
		}
		return syscall.Chroot(".")
	}

	err = syscall.Chroot(rootFs)
	if err != nil {
		return closeFunc, err
	}

	err = os.Chdir("/")
	if err != nil {
		return closeFunc, err
	}

	return closeFunc, nil
}

// RunChroot runs a function in a chroot and changes the current working
// environment to the new root. It exits the chroot after the function
// returns. The function returns the error returned by the function and
// the error returned by ExitChroot if any.
//
// Example:
//
//	err := chroot.RunChroot("/path/to/new/root", func() error {
//		return exec.Command("ls", "-l").Run()
//	})
//	if err != nil {
//		fmt.Printf("Error running command in chroot root: %v\n", err)
//		return
//	}
func RunChroot(rootFs string, f func() error) (fErr, err error) {
	fd, err := os.Open("/")
	if err != nil {
		return nil, err
	}

	err = syscall.Chroot(rootFs)
	if err != nil {
		fd.Close()
		return nil, err
	}

	err = os.Chdir("/")
	if err != nil {
		fd.Close()
		err = syscall.Chroot(".")
		return nil, err
	}

	fErr = f()

	fd.Chdir()
	fd.Close()

	return fErr, syscall.Chroot(".")
}
