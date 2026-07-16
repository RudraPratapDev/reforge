package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// captureStdout runs fn while temporarily redirecting os.Stdout to a pipe,
// returning everything written to stdout during the call.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w
	defer func() {
		os.Stdout = origStdout
	}()

	fn()

	if err := w.Close(); err != nil {
		t.Fatalf("failed to close pipe writer: %v", err)
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("failed to read from pipe: %v", err)
	}

	return buf.String()
}

func TestMain_PrintsStartupMessage(t *testing.T) {
	output := captureStdout(t, main)

	want := "Yep I have started REFORGE !!"
	if !strings.Contains(output, want) {
		t.Errorf("main() output = %q, want it to contain %q", output, want)
	}
}

func TestMain_DoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("main() panicked: %v", r)
		}
	}()

	captureStdout(t, main)
}