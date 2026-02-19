//go:build linux

package fs

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// MountInfo represents a single entry from /proc/self/mountinfo.
//
// Format reference: https://www.kernel.org/doc/Documentation/filesystems/proc.txt
//
// Fields are kept close to the kernel representation to avoid lossy conversions.
// Only a subset is currently exposed as helpers in this package.
type MountInfo struct {
	MountID        int
	ParentID       int
	Major          int
	Minor          int
	Root           string
	MountPoint     string
	Options        []string
	OptionalFields []string
	FSType         string
	Source         string
	SuperOptions   []string
}

func parseMountInfo(r io.Reader) ([]MountInfo, error) {
	entries := make([]MountInfo, 0)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		mi, err := parseMountInfoLine(line)
		if err != nil {
			return nil, err
		}
		entries = append(entries, mi)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func parseMountInfoLine(line string) (MountInfo, error) {
	// Split on the separator " - " (space, dash, space)
	parts := strings.SplitN(line, " - ", 2)
	if len(parts) != 2 {
		return MountInfo{}, fmt.Errorf("invalid mountinfo line: missing separator")
	}

	pre := strings.Fields(parts[0])
	if len(pre) < 6 {
		return MountInfo{}, fmt.Errorf("invalid mountinfo line: too few fields")
	}

	post := strings.Fields(parts[1])
	if len(post) < 3 {
		return MountInfo{}, fmt.Errorf("invalid mountinfo line: too few post-separator fields")
	}

	mountID, err := strconv.Atoi(pre[0])
	if err != nil {
		return MountInfo{}, err
	}
	parentID, err := strconv.Atoi(pre[1])
	if err != nil {
		return MountInfo{}, err
	}

	majMin := strings.SplitN(pre[2], ":", 2)
	if len(majMin) != 2 {
		return MountInfo{}, fmt.Errorf("invalid mountinfo line: invalid major:minor")
	}
	major, err := strconv.Atoi(majMin[0])
	if err != nil {
		return MountInfo{}, err
	}
	minor, err := strconv.Atoi(majMin[1])
	if err != nil {
		return MountInfo{}, err
	}

	root := unescapeMountField(pre[3])
	mountPoint := unescapeMountField(pre[4])
	options := splitComma(pre[5])

	optional := make([]string, 0)
	if len(pre) > 6 {
		for _, f := range pre[6:] {
			optional = append(optional, unescapeMountField(f))
		}
	}

	fsType := post[0]
	source := unescapeMountField(post[1])
	superOptions := splitComma(post[2])

	return MountInfo{
		MountID:        mountID,
		ParentID:       parentID,
		Major:          major,
		Minor:          minor,
		Root:           root,
		MountPoint:     mountPoint,
		Options:        options,
		OptionalFields: optional,
		FSType:         fsType,
		Source:         source,
		SuperOptions:   superOptions,
	}, nil
}

func splitComma(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
}

func unescapeMountField(s string) string {
	// mountinfo uses octal escapes for special characters (e.g. \040 for space)
	if !strings.Contains(s, "\\") {
		return s
	}

	b := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] != '\\' || i+3 >= len(s) {
			b = append(b, s[i])
			continue
		}

		o := s[i+1 : i+4]
		v, err := strconv.ParseInt(o, 8, 32)
		if err != nil {
			b = append(b, s[i])
			continue
		}
		b = append(b, byte(v))
		i += 3
	}
	return string(b)
}

// GetMountInfo returns the current mount table by reading /proc/self/mountinfo.
//
// Example:
//
//	entries, err := fs.GetMountInfo()
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
//	for _, e := range entries {
//		fmt.Printf("%s -> %s (%s)\n", e.Source, e.MountPoint, e.FSType)
//	}
func GetMountInfo() ([]MountInfo, error) {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return parseMountInfo(f)
}

// GetMountpoint returns the mountpoint for a given mount source (e.g. /dev/sda1).
// If the source is not mounted, an empty string is returned.
//
// Example:
//
//	mp, err := fs.GetMountpoint("/dev/sda1")
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return
//	}
//	fmt.Printf("Mountpoint: %s\n", mp)
func GetMountpoint(source string) (string, error) {
	entries, err := GetMountInfo()
	if err != nil {
		return "", err
	}
	for _, e := range entries {
		if e.Source == source {
			return e.MountPoint, nil
		}
	}
	return "", nil
}
