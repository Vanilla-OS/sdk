package tests

import (
	_ "embed"
	"testing"

	"github.com/vanilla-os/sdk/pkg/v1/goodies"
)

//go:embed resources/files_ok.txt
var filesOkContent []byte

//go:embed resources/files_fail.txt
var filesFailContent []byte

func TestIntegrityCheck(t *testing.T) {
	// Perform the integrity check with the correct hash
	resultOk, err := goodies.CheckIntegrity(filesOkContent, goodies.SHA1Validator{}, "resources")
	if err != nil {
		t.Fatalf("Unexpected error during integrity check with files_ok.txt: %v", err)
	} else {
		t.Logf("Total requested: %d\n", resultOk.TotalRequested)
		t.Logf("Passed: %d\n", resultOk.Passed)
		t.Logf("Failed: %d\n", resultOk.Failed)
		for _, failedCheck := range resultOk.FailedChecks {
			t.Logf("Resource: %s\n", failedCheck.ResourcePath)
			t.Logf("Requested: %s\n", failedCheck.RequestedHash)
			t.Logf("Detected: %s\n", failedCheck.DetectedHash)
		}
	}

	if resultOk.Failed > 0 {
		t.Errorf("Expected all checks to pass, but %d checks failed", resultOk.Failed)
	} else {
		t.Logf("All checks passed")
	}

	// Perform the integrity check with the incorrect hash
	resultFail, err := goodies.CheckIntegrity(filesFailContent, goodies.SHA1Validator{}, "resources")
	if err != nil {
		t.Fatalf("Unexpected error during integrity check with files_fail.txt: %v", err)
	} else {
		t.Logf("Total requested: %d\n", resultFail.TotalRequested)
		t.Logf("Passed: %d\n", resultFail.Passed)
		t.Logf("Failed: %d\n", resultFail.Failed)
		for _, failedCheck := range resultFail.FailedChecks {
			t.Logf("Resource: %s\n", failedCheck.ResourcePath)
			t.Logf("Requested: %s\n", failedCheck.RequestedHash)
			t.Logf("Detected: %s\n", failedCheck.DetectedHash)
		}
	}

	if resultFail.Failed == 0 {
		t.Errorf("Expected some checks to fail, but all checks passed")
	} else {
		t.Logf("Some checks failed as expected")
	}
}
