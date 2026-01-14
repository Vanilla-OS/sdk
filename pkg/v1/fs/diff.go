package fs

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

import (
	"os"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/vanilla-os/sdk/pkg/v1/fs/types"
)

// GetFileDiff compares the content of two files and returns the changes
//
// Example:
//
//	diff, err := fs.GetFileDiff("/tmp/batman", "/tmp/robin")
//	if err != nil {
//		fmt.Printf("Error getting file diff: %v", err)
//		return
//	}
//	fmt.Printf("Added lines: %v\n", diff.AddedLines)
//	fmt.Printf("Removed lines: %v\n", diff.RemovedLines)
func GetFileDiff(firstFile, secondFile string) (types.FileDiffInfo, error) {
	firstContent, err := os.ReadFile(firstFile)
	if err != nil {
		return types.FileDiffInfo{}, err
	}

	secondContent, err := os.ReadFile(secondFile)
	if err != nil {
		return types.FileDiffInfo{}, err
	}

	diff := types.FileDiffInfo{}

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(firstContent), string(secondContent), false)
	for _, d := range diffs {
		switch d.Type {
		case diffmatchpatch.DiffInsert:
			diff.AddedLines = append(diff.AddedLines, d.Text)
		case diffmatchpatch.DiffDelete:
			diff.RemovedLines = append(diff.RemovedLines, d.Text)
		}
	}

	return diff, nil
}
