package i18n

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
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// CheckMissingStrings scans the rootPath for translation calls (e.g., .Trans("KEY"))
// and verifies they exist in the provided locale file (PO format).
//
// It returns a map where the key is the file path and the value is a slice of missing keys.
//
// Example:
//
//	missing, err := i18n.CheckMissingStrings(".", "locales/en/messages.po")
//	if err != nil {
//		log.Fatal(err)
//	}
//	for file, keys := range missing {
//		fmt.Printf("File %s has missing keys: %v\n", file, keys)
//	}
func CheckMissingStrings(rootPath string, localeFile string) (map[string][]string, error) {
	localeData, err := os.Open(localeFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open locale file: %w", err)
	}
	defer localeData.Close()

	definedKeys := make(map[string]bool)
	scanner := bufio.NewScanner(localeData)

	msgidRe := regexp.MustCompile(`^msgid "(.+)"$`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "msgid \"") {
			matches := msgidRe.FindStringSubmatch(line)
			if len(matches) > 1 {
				definedKeys[matches[1]] = true
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading locale file: %w", err)
	}

	sourceRe := regexp.MustCompile(`\.Trans\("([\S.]+)"`)
	prRe := regexp.MustCompile(`pr:([\w\.]+)`)
	missing := make(map[string][]string)

	err = filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" || d.Name() == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".go") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		matches := sourceRe.FindAllStringSubmatch(string(content), -1)
		matches = append(matches, prRe.FindAllStringSubmatch(string(content), -1)...)
		for _, match := range matches {
			if len(match) < 2 {
				continue
			}
			key := match[1]
			if !definedKeys[key] {
				seen := false
				for _, k := range missing[path] {
					if k == key {
						seen = true
						break
					}
				}
				if !seen {
					missing[path] = append(missing[path], key)
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return missing, nil
}
