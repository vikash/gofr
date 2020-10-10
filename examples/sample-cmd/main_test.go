package main

import (
	"testing"

	"github.com/vikash/gofr/pkg/gofr/testUtil"
)

// TestCMDRunWithNoArg checks that if no subcommand is found then error comes on stderr.
func TestCMDRunWithNoArg(t *testing.T) {
	expectedError := "No Command Found!"
	output := testUtil.StderrOutputForFunc(main)
	if output != "No Command Found!" {
		t.Errorf("Expected: %s\n Got: %s", expectedError, output)
	}
}