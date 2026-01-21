package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootCmd_Version(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		version  string
		expected string
	}{
		{
			name:     "version flag long",
			args:     []string{"--version"},
			version:  "dev",
			expected: "lane version dev",
		},
		{
			name:     "version flag short",
			args:     []string{"-v"},
			version:  "dev",
			expected: "lane version dev",
		},
		{
			name:     "version with custom value",
			args:     []string{"--version"},
			version:  "1.0.0",
			expected: "lane version 1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original version and restore after test
			origVersion, origCommit, origDate := GetVersionInfo()
			defer SetVersionInfo(origVersion, origCommit, origDate)

			SetVersionInfo(tt.version, "test", "test")

			cmd := NewRootCmd()
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			output := buf.String()
			if !strings.Contains(output, tt.expected) {
				t.Errorf("expected output to contain %q, got %q", tt.expected, output)
			}
		})
	}
}

func TestRootCmd_Help(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected []string
	}{
		{
			name: "help flag long",
			args: []string{"--help"},
			expected: []string{
				"lane is a CLI tool",
				"git worktrees and tmux",
				"Usage:",
				"Flags:",
			},
		},
		{
			name: "help flag short",
			args: []string{"-h"},
			expected: []string{
				"lane is a CLI tool",
				"Usage:",
			},
		},
		{
			name: "no args shows help",
			args: []string{},
			expected: []string{
				"lane is a CLI tool",
				"Usage:",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewRootCmd()
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			output := buf.String()
			for _, exp := range tt.expected {
				if !strings.Contains(output, exp) {
					t.Errorf("expected output to contain %q, got %q", exp, output)
				}
			}
		})
	}
}

func TestRootCmd_UnknownCommand(t *testing.T) {
	cmd := NewRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"unknowncommand"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for unknown command, got nil")
	}

	errStr := err.Error()
	// Cobra's error message for unknown positional args with NoArgs validator
	if !strings.Contains(errStr, "unknown command") && !strings.Contains(errStr, "unknown") {
		t.Errorf("expected error to mention unknown, got %q", errStr)
	}
}

func TestRootCmd_VerboseFlag(t *testing.T) {
	// Reset Verbose flag
	Verbose = false
	defer func() { Verbose = false }()

	cmd := NewRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"--verbose", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !Verbose {
		t.Error("expected Verbose flag to be true after --verbose")
	}
}

func TestRootCmd_VerboseFlagShort(t *testing.T) {
	// Reset Verbose flag
	Verbose = false
	defer func() { Verbose = false }()

	cmd := NewRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"-V", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !Verbose {
		t.Error("expected Verbose flag to be true after -V")
	}
}

func TestVersionFull(t *testing.T) {
	// Save original and restore after test
	origVersion, origCommit, origDate := GetVersionInfo()
	defer SetVersionInfo(origVersion, origCommit, origDate)

	SetVersionInfo("1.2.3", "abc123", "2024-01-01T00:00:00Z")

	full := VersionFull()
	expected := "lane version 1.2.3 (commit: abc123, built: 2024-01-01T00:00:00Z)"

	if full != expected {
		t.Errorf("expected %q, got %q", expected, full)
	}
}

func TestRootCmd_HasProperDescriptions(t *testing.T) {
	cmd := NewRootCmd()

	if cmd.Use != "lane" {
		t.Errorf("expected Use to be 'lane', got %q", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("expected Short description to be non-empty")
	}

	if cmd.Long == "" {
		t.Error("expected Long description to be non-empty")
	}

	if !strings.Contains(cmd.Short, "git worktrees") {
		t.Error("expected Short description to mention git worktrees")
	}

	if !strings.Contains(cmd.Short, "tmux") {
		t.Error("expected Short description to mention tmux")
	}
}
