package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestNoArguments(t *testing.T) {
	out, _ := execGo("run", "main.go")

	if out != "Missing hashtag\nexit status 1\n" {
		t.Fatalf(`unexpected output %q`, out)
	}
}

func TestUsageInstructions(t *testing.T) {
	out, _ := execGo("run", "main.go", "-h")

	if !strings.Contains(out, "Usage of") {
		t.Fatalf(`expected usage instructions`)
	}

	if !strings.Contains(out, "-l int") {
		t.Fatalf(`expected -l flag`)
	}

	if !strings.Contains(out, "-s") {
		t.Fatalf(`expected -s flag`)
	}
}

func execGo(args ...string) (string, error) {
	out, err := exec.Command("go", args...).CombinedOutput()

	return string(out), err
}
