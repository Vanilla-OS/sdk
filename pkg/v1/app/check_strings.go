//go:build check_missing_strings
// +build check_missing_strings

package app

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/vanilla-os/sdk/pkg/v1/i18n"
)

// checkMissingStrings is a hook used when the "check_missing_strings" build tag is set.
// It scans the project for missing translation keys and exists with a non-zero status if any are found.
func (a *App) checkMissingStrings() {
	rootDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	// Try to locate the English locale file in common paths.
	var localeFile string
	possiblePaths := []string{
		"cmd/locales/en/LC_MESSAGES",
		"cmd/locales/en",
		"locales/en/LC_MESSAGES",
		"locales/en",
		"assets/locales/en/LC_MESSAGES",
		"assets/locales/en",
	}

	for _, p := range possiblePaths {
		path := filepath.Join(rootDir, p)
		matches, _ := filepath.Glob(filepath.Join(path, "*.po"))
		if len(matches) > 0 {
			localeFile = matches[0]
			break
		}
	}

	if localeFile == "" {
		fmt.Println("Error: Could not find any .po file in standard locations")
		os.Exit(1)
	}

	fmt.Printf("Checking missing strings in %s using locale %s\n", rootDir, localeFile)

	missing, err := i18n.CheckMissingStrings(rootDir, localeFile)
	if err != nil {
		fmt.Printf("Error checking strings: %v\n", err)
		os.Exit(1)
	}

	if len(missing) > 0 {
		fmt.Println("Oops, there are missing translation strings:")
		for file, keys := range missing {
			for _, key := range keys {
				fmt.Printf("- %s: Missing '%s'\n", file, key)
			}
		}

		f, _ := os.Create("missing_strings.json")
		defer f.Close()
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		enc.Encode(missing)

		os.Exit(1)
	}

	fmt.Println("No missing translation strings. Good job!")
	os.Exit(0)
}
