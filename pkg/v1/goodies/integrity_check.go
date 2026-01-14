package goodies

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ValidatorType defines the protocol for a validation method. Implement this
// to create custom validation methods for your integrity checks.
type ValidatorType interface {
	Hash(data io.Reader) (string, error)
}

// SHA1Validator is an implementation of the validation method using SHA-1
type SHA1Validator struct{}

func (s SHA1Validator) Hash(data io.Reader) (string, error) {
	h := sha1.New()
	if _, err := io.Copy(h, data); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// SHA256Validator is an implementation of the validation method using SHA-256
type SHA256Validator struct{}

func (s SHA256Validator) Hash(data io.Reader) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, data); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// MD5Validator is an implementation of the validation method using MD5
type MD5Validator struct{}

func (m MD5Validator) Hash(data io.Reader) (string, error) {
	h := md5.New()
	if _, err := io.Copy(h, data); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// IntegrityCheckResult represents the result of an integrity check
type IntegrityCheckResult struct {
	TotalRequested int
	Passed         int
	Failed         int
	FailedChecks   []FailedCheck
}

// FailedCheck represents a failed integrity check
type FailedCheck struct {
	ResourcePath  string
	RequestedHash string
	DetectedHash  string
}

// CheckIntegrity performs an integrity check of the provided data using the
// provided validation method and returns the result. If a validation fails
// the result will contain the failed checks.
//
// The prefix parameter specifies the base path to be used for each file check.
//
// Example:
//
//	data := []byte("batman.txt 8c555f537cd1fe3f1239ebf7b6d639bc0d576bda\nrobin.txt 4e58c66d2ef2b9a60c7ea2bc03253d8b01874b52")
//	result, err := CheckIntegrity(data, SHA1Validator{}, "files/")
//
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//	}
//
//	fmt.Printf("Total requested: %d\n", result.TotalRequested)
//	fmt.Printf("Passed: %d\n", result.Passed)
//	fmt.Printf("Failed: %d\n", result.Failed)
//	for _, failedCheck := range result.FailedChecks {
//		fmt.Printf("Resource: %s\n", failedCheck.ResourcePath)
//		fmt.Printf("Requested: %s\n", failedCheck.RequestedHash)
//		fmt.Printf("Detected: %s\n", failedCheck.DetectedHash)
//	}
func CheckIntegrity(data []byte, method ValidatorType, prefix string) (IntegrityCheckResult, error) {
	pairings := make(map[string]string)
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			return IntegrityCheckResult{}, fmt.Errorf("invalid line format: %s", line)
		}
		pairings[parts[0]] = parts[1]
	}

	// Perform the integrity check using the provided validation method
	return checkIntegrity(pairings, method, prefix)
}

// checkIntegrity performs the integrity check of the file data using
// the provided validation method.
func checkIntegrity(pairings map[string]string, method ValidatorType, prefix string) (IntegrityCheckResult, error) {
	var result IntegrityCheckResult

	for resourcePath, requestedHash := range pairings {
		// If a prefix is provided, we have to join it with the resource path to
		// ensure the file is found in the place the developer expects.
		fullPath := filepath.Join(prefix, resourcePath)
		detectedHash, err := calculateHash(fullPath, method)
		if err != nil {
			fmt.Printf("Error calculating hash for %s: %v\n", fullPath, err)
			continue
		}

		result.TotalRequested++
		if requestedHash == detectedHash {
			result.Passed++
		} else {
			result.Failed++
			result.FailedChecks = append(result.FailedChecks, FailedCheck{
				ResourcePath:  fullPath,
				RequestedHash: requestedHash,
				DetectedHash:  detectedHash,
			})
		}
	}

	return result, nil
}

// calculateHash calculates the hash of a file using the provided validation
// method and returns it as a string.
func calculateHash(filePath string, method ValidatorType) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return method.Hash(file)
}
