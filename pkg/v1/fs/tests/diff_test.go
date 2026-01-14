package tests

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
	"testing"

	"github.com/vanilla-os/sdk/pkg/v1/fs"
)

func TestDiff(t *testing.T) {
	tmpDir := os.TempDir()

	// We expect the diff to be:
	// 1. Added line: fourth line
	// 2. Removed line: third line
	firstFile := filepath.Join(tmpDir, "first.txt")
	secondFile := filepath.Join(tmpDir, "second.txt")

	firstContent := "first line\nsecond line\nthird line\n"
	secondContent := "first line\nsecond line\nthird line\nfourth line\n"

	if err := os.WriteFile(firstFile, []byte(firstContent), 0644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(secondFile, []byte(secondContent), 0644); err != nil {
		t.Fatal(err)
	}

	diff, err := fs.GetFileDiff(firstFile, secondFile)
	if err != nil {
		t.Fatal(err)
	}

	if len(diff.AddedLines) != 1 {
		t.Fatalf("expected 1 added line, got %d", len(diff.AddedLines))
	}

	if diff.AddedLines[0] != "fourth line\n" {
		t.Fatalf("expected fourth line, got %s", diff.AddedLines[0])
	}

	if len(diff.RemovedLines) != 0 {
		t.Fatalf("expected 0 removed lines, got %d", len(diff.RemovedLines))
	}

	// We expect the opposite diff to be:
	// 1. Removed line: fourth line
	diff, err = fs.GetFileDiff(secondFile, firstFile)
	if err != nil {
		t.Fatal(err)
	}

	if len(diff.AddedLines) != 0 {
		t.Fatalf("expected 0 added lines, got %d", len(diff.AddedLines))
	}

	if len(diff.RemovedLines) != 1 {
		t.Fatalf("expected 1 removed line, got %d", len(diff.RemovedLines))
	}

	if diff.RemovedLines[0] != "fourth line\n" {
		t.Fatalf("expected fourth line, got %s", diff.RemovedLines[0])
	}
}
