package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCreateDeployLockFileCreatesLockContent(t *testing.T) {
	t.Parallel()

	lockFile := filepath.Join(t.TempDir(), ".deploy.test.lock")
	startedAt := time.Date(2026, 5, 18, 10, 20, 30, 0, time.UTC)

	if err := createDeployLockFile(lockFile, startedAt); err != nil {
		t.Fatalf("createDeployLockFile() error = %v", err)
	}

	content, err := os.ReadFile(lockFile)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	text := string(content)
	if !strings.Contains(text, "pid=") {
		t.Fatalf("lock file missing pid: %q", text)
	}
	if !strings.Contains(text, "started=2026-05-18T10:20:30Z") {
		t.Fatalf("lock file missing started timestamp: %q", text)
	}
}

func TestCreateDeployLockFileFailsWhenLockAlreadyExists(t *testing.T) {
	t.Parallel()

	lockFile := filepath.Join(t.TempDir(), ".deploy.test.lock")
	startedAt := time.Date(2026, 5, 18, 10, 20, 30, 0, time.UTC)

	if err := createDeployLockFile(lockFile, startedAt); err != nil {
		t.Fatalf("first createDeployLockFile() error = %v", err)
	}

	err := createDeployLockFile(lockFile, startedAt.Add(time.Second))
	if !errors.Is(err, os.ErrExist) {
		t.Fatalf("second createDeployLockFile() error = %v, want os.ErrExist", err)
	}
}