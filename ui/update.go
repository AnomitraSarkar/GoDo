package ui

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"
)

var ErrNotGitRepo = errors.New("not a git repository")

func getCurrentBranch() (string, error) {
	branchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	var stderr bytes.Buffer
	branchCmd.Stderr = &stderr
	branchOut, err := branchCmd.Output()
	if err != nil {
		if stderr.Len() > 0 && strings.Contains(stderr.String(), "not a git repository") {
			return "", ErrNotGitRepo
		}
		return "", err
	}
	return strings.TrimSpace(string(branchOut)), nil
}

func checkForUpdate() (bool, error) {
	branch, err := getCurrentBranch()
	if err == ErrNotGitRepo {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	fetchCmd := exec.Command("git", "fetch", "origin", branch)
	fetchCmd.Stdout = nil
	fetchCmd.Stderr = nil
	if err := fetchCmd.Run(); err != nil {
		return false, err
	}

	localCmd := exec.Command("git", "rev-parse", branch)
	localOut, err := localCmd.Output()
	if err != nil {
		return false, err
	}
	localCommit := strings.TrimSpace(string(localOut))

	remoteCmd := exec.Command("git", "rev-parse", "origin/"+branch)
	remoteOut, err := remoteCmd.Output()
	if err != nil {
		return false, err
	}
	remoteCommit := strings.TrimSpace(string(remoteOut))

	return localCommit != remoteCommit, nil
}

func updateApplication() error {
	branch, err := getCurrentBranch()
	if err != nil {
		return err
	}
	cmd := exec.Command("git", "merge", "--ff-only", "origin/"+branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}