package githelpers

import (
	"testing"
	"strings"
	"os"
	"os/exec"
	"github.com/McdonaldSeanp/charlie/utils"
)

func initRepo(t *testing.T) {
	_, err := utils.ExecReadOutput(exec.Command("git", "init"))
	if err != nil {
		t.Fatalf("Failed to create test repo: %s\n", err)
	}
	f, err := os.Create("test.txt")
	if err != nil {
		t.Fatalf("Failed to create test repo: %s\n", err)
	}
	f.Close()
	_, err = utils.ExecReadOutput(exec.Command("git", "add", "--all"))
	if err != nil {
		t.Fatalf("Failed to create test repo: %s\n", err)
	}
	_, err = utils.ExecReadOutput(exec.Command("git", "commit", "-m", "initial", "--no-gpg-sign"))
	if err != nil {
		t.Fatalf("Failed to create test repo: %s\n", err)
	}
	_, err = utils.ExecReadOutput(exec.Command("git", "checkout", "-B", "main"))
	if err != nil {
		t.Fatalf("Failed to create test repo: %s\n", err)
	}
}

func getBranch(t *testing.T) string {
	result, err := utils.ExecReadOutput(exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD"))
	if err != nil {
		t.Fatalf("Failed to read branch: %s\n", err)
	}
	return strings.Trim(result, " \n\r\t")
}

func fakeCommit(dir string, t *testing.T) {
	fakeFile(dir, t)
	_, err := utils.ExecReadOutput(exec.Command("git", "add", "--all"))
	if err != nil {
		t.Fatalf("Failed to create new commit: %s\n", err)
	}
	_, err = utils.ExecReadOutput(exec.Command("git", "commit", "-m", "fake", "--no-gpg-sign"))
	if err != nil {
		t.Fatalf("Failed to create new commit: %s\n", err)
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
	err := Setgitbranch("main", false)
	if err != nil {
		t.Fatalf("Setting the git branch failed: %s\n", err)
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
	err := Setgitbranch("", false)
	if err == nil {
		t.Fatalf("Setting the git branch was supposed to fail!")
	}
	// The branch should remain the same
	branch_now := getBranch(t)
	if branch_now != "main" {
		t.Fatalf("Branch is not correctly set, should be 'main', is '%s'\n", branch_now)
	}
	// run again with a real name
	err = Setgitbranch("non_existant", false)
	if err == nil {
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
	_, err := utils.ExecReadOutput(exec.Command("git", "checkout", "-B", "different"))
	if err != nil {
		t.Fatalf("Creating second git branch failed: %s\n", err)
	}
	// Check to make sure we are actually on a new branch
	branch_now := getBranch(t)
	if branch_now != "different" {
		t.Fatalf("Branch is not correctly set, should be 'different', is '%s'\n", branch_now)
	}
	// Create a fake commit so the new branch is ahead of main
	fakeCommit(testdir, t)
	// Actually run the thing
	err = Setgitbranch("main", false)
	if err != nil {
		t.Fatalf("Setting the git branch failed: %s\n", err)
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
	_, err := utils.ExecReadOutput(exec.Command("git", "checkout", "-B", "different"))
	if err != nil {
		t.Fatalf("Creating second git branch failed: %s\n", err)
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
	err = Setgitbranch("main", false)
	if err == nil {
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
	_, err := utils.ExecReadOutput(exec.Command("git", "checkout", "-B", "different"))
	if err != nil {
		t.Fatalf("Creating second git branch failed: %s\n", err)
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
	err = Setgitbranch("main", true)
	if err != nil {
		t.Fatalf("Setting the git branch failed: %s\n", err)
	}
	branch_now = getBranch(t)
	if branch_now != "main" {
		t.Fatalf("Branch is not correctly set, should be 'different', is '%s'\n", branch_now)
	}
	// Check that the work tree is clean again
	wt, err := OpenWorktree()
	if err != nil {
		t.Fatalf("Getting the work tree failed: %s\n", err)
	}
	clean, err := WorkTreeClean(wt)
	if err != nil {
		t.Fatalf("Getting the work tree failed: %s\n", err)
	}
	if !clean {
		t.Fatalf("Working tree is not clean")
	}
}
