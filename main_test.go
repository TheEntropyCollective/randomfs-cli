package main

import (
	"os/exec"
	"testing"
)

func TestCLIRuns(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "--help")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI did not run: %v, output: %s", err, out)
	}
	if len(out) == 0 {
		t.Error("No output from CLI help")
	}
}
