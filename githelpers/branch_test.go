package githelpers

import (
	"os"
	"strings"
	"testing"

	"github.com/mcdonaldseanp/charlie/localexec"
	"github.com/mcdonaldseanp/charlie/utils"
)

func initRepo(t *testing.T) {
	_, _, airr := localexec.ExecReadOutput("git", "init")
	if airr != nil {
		t.Fatalf("Failed to create test repo: %s\n", airr)
	}
	f, err := os.Create("test.txt")
	if err != nil {
		t.Fatalf("Failed to create test repo: %s\n", err)
	}
	f.Close()
	_, _, airr = localexec.ExecReadOutput("git", "add", "--all")
	if airr != nil {
		t.Fatalf("Failed to create test repo: %s\n", airr)
	}
	_, _, airr = localexec.ExecReadOutput("git", "commit", "-m", "initial", "--no-gpg-sign")
	if airr != nil {
		t.Fatalf("Failed to create test repo: %s\n", airr)
	}
	_, _, airr = localexec.ExecReadOutput("git", "checkout", "-B", "main")
	if airr != nil {
		t.Fatalf("Failed to create test repo: %s\n", airr)
	}
}

func getBranch(t *testing.T) string {
	result, _, airr := localexec.ExecReadOutput("git", "rev-parse", "--abbrev-ref", "HEAD")
	if airr != nil {
		t.Fatalf("Failed to read branch: %s\n", airr)
	}
	return strings.Trim(result, " \n\r\t")
}

func fakeCommit(dir string, t *testing.T) {
	fakeFile(dir, t)
	_, _, airr := localexec.ExecReadOutput("git", "add", "--all")
	if airr != nil {
		t.Fatalf("Failed to create new commit: %s\n", airr)
	}
	_, _, airr = localexec.ExecReadOutput("git", "commit", "-m", "fake", "--no-gpg-sign")
	if airr != nil {
		t.Fatalf("Failed to create new commit: %s\n", airr)
	}
}

func fakeFile(dir string, t *testing.T) {
	f, err := os.CreateTemp(dir, "")
	if err != nil {
		t.Fatalf("Failed to create file: %s\n", err)
	}
	f.Close()
}

// Test that checking out the already checked out branch does nothing
func TestBranchSame(t *testing.T) {
	originaldir, _ := os.Getwd()
	testdir, _ := os.MkdirTemp("", "charlie_testing")
	defer os.RemoveAll(testdir)
	defer os.Chdir(originaldir)
	os.Chdir(testdir)
	initRepo(t)
	// Actually run the thing
	airr := SetBranch("main", false, false)
	if airr != nil {
		t.Fatalf("Setting the git branch failed: %s\n", airr)
	}
	branch_now := getBranch(t)
	if branch_now != "main" {
		t.Fatalf("Branch is not correctly set, should be 'main', is '%s'\n", branch_now)
	}
}

// Test that checking out an empty branch or one that doesn't exist fails
// but does not change the current branch
func TestBranchNoExist(t *testing.T) {
	originaldir, _ := os.Getwd()
	testdir, _ := os.MkdirTemp("", "charlie_testing")
	defer os.RemoveAll(testdir)
	defer os.Chdir(originaldir)
	os.Chdir(testdir)
	initRepo(t)
	// Actually run the thing
	airr := SetBranch("", false, false)
	if airr == nil {
		t.Fatalf("Setting the git branch was supposed to fail!")
	}
	// The branch should remain the same
	branch_now := getBranch(t)
	if branch_now != "main" {
		t.Fatalf("Branch is not correctly set, should be 'main', is '%s'\n", branch_now)
	}
	// run again with a real name
	airr = SetBranch("non_existant", false, false)
	if airr == nil {
		t.Fatalf("Setting the git branch was supposed to fail!")
	}
	// The branch should remain the same
	branch_now = getBranch(t)
	if branch_now != "main" {
		t.Fatalf("Branch is not correctly set, should be 'main', is '%s'\n", branch_now)
	}
}

// Test that checking out a separate branch works
func TestBranchDifferent(t *testing.T) {
	originaldir, _ := os.Getwd()
	testdir, _ := os.MkdirTemp("", "charlie_testing")
	defer os.RemoveAll(testdir)
	defer os.Chdir(originaldir)
	os.Chdir(testdir)
	initRepo(t)
	_, _, airr := localexec.ExecReadOutput("git", "checkout", "-B", "different")
	if airr != nil {
		t.Fatalf("Creating second git branch failed: %s\n", airr)
	}
	// Check to make sure we are actually on a new branch
	branch_now := getBranch(t)
	if branch_now != "different" {
		t.Fatalf("Branch is not correctly set, should be 'different', is '%s'\n", branch_now)
	}
	// Create a fake commit so the new branch is ahead of main
	fakeCommit(testdir, t)
	// Actually run the thing
	airr = SetBranch("main", false, false)
	if airr != nil {
		t.Fatalf("Setting the git branch failed: %s\n", airr)
	}
	branch_now = getBranch(t)
	if branch_now != "main" {
		t.Fatalf("Branch is not correctly set, should be 'main', is '%s'\n", branch_now)
	}
}

// Test that checking out a separate branch when the worktree is dirty fails
func TestBranchDirty(t *testing.T) {
	originaldir, _ := os.Getwd()
	testdir, _ := os.MkdirTemp("", "charlie_testing")
	defer os.RemoveAll(testdir)
	defer os.Chdir(originaldir)
	os.Chdir(testdir)
	initRepo(t)
	_, _, airr := localexec.ExecReadOutput("git", "checkout", "-B", "different")
	if airr != nil {
		t.Fatalf("Creating second git branch failed: %s\n", airr)
	}
	// Check to make sure we are actually on a new branch
	branch_now := getBranch(t)
	if branch_now != "different" {
		t.Fatalf("Branch is not correctly set, should be 'different', is '%s'\n", branch_now)
	}
	// Create a fake commit so the new branch is ahead of main
	fakeCommit(testdir, t)
	// Create an uncommitted file
	fakeFile(testdir, t)
	// Actually run the thing
	airr = SetBranch("main", false, false)
	if airr == nil {
		t.Fatalf("Setting the git branch was supposed to fail!")
	}
	branch_now = getBranch(t)
	if branch_now != "different" {
		t.Fatalf("Branch is not correctly set, should be 'different', is '%s'\n", branch_now)
	}
}

// Test that checking out with --clear set will clean the branch and switch
func TestBranchClear(t *testing.T) {
	originaldir, _ := os.Getwd()
	testdir, _ := os.MkdirTemp("", "charlie_testing")
	defer os.RemoveAll(testdir)
	defer os.Chdir(originaldir)
	os.Chdir(testdir)
	initRepo(t)
	_, _, airr := localexec.ExecReadOutput("git", "checkout", "-B", "different")
	if airr != nil {
		t.Fatalf("Creating second git branch failed: %s\n", airr)
	}
	// Check to make sure we are actually on a new branch
	branch_now := getBranch(t)
	if branch_now != "different" {
		t.Fatalf("Branch is not correctly set, should be 'different', is '%s'\n", branch_now)
	}
	// Create a fake commit so the new branch is ahead of main
	fakeCommit(testdir, t)
	// Create an uncommitted file
	fakeFile(testdir, t)
	// Actually run the thing
	airr = SetBranch("main", true, false)
	if airr != nil {
		t.Fatalf("Setting the git branch failed: %s\n", airr)
	}
	branch_now = getBranch(t)
	if branch_now != "main" {
		t.Fatalf("Branch is not correctly set, should be 'different', is '%s'\n", branch_now)
	}
	// Check that the work tree is clean again
	wt, airr := utils.OpenWorktree()
	if airr != nil {
		t.Fatalf("Getting the work tree failed: %s\n", airr)
	}
	clean, airr := utils.WorkTreeClean(wt)
	if airr != nil {
		t.Fatalf("Getting the work tree failed: %s\n", airr)
	}
	if !clean {
		t.Fatalf("Working tree is not clean")
	}
}
