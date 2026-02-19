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

func TestCopyTree(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	if err := os.MkdirAll(filepath.Join(src, "a", "b"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(src, "a", "b", "file.txt"), []byte("hello"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := os.Symlink("../b/file.txt", filepath.Join(src, "a", "link")); err != nil {
		t.Fatalf("symlink: %v", err)
	}

	out := filepath.Join(dst, "out")
	if err := fs.CopyTree(src, out, fs.CopyTreeOptions{Workers: 2}); err != nil {
		t.Fatalf("copytree: %v", err)
	}

	b, err := os.ReadFile(filepath.Join(out, "a", "b", "file.txt"))
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if string(b) != "hello" {
		t.Fatalf("unexpected content: %q", string(b))
	}

	if _, err := os.Readlink(filepath.Join(out, "a", "link")); err != nil {
		t.Fatalf("readlink: %v", err)
	}
}
