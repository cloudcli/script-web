package git

import (
	"os/exec"
	"bytes"
	"os"
	"path/filepath"
	"github.com/pkg/errors"
)

var gitCmd string
var gitRepo *Repo

func GetRepo(path string) (*Repo, error) {
	if gitRepo == nil {
		var err error
		if gitCmd, err = exec.LookPath("git"); err != nil {
			return nil, err
		}
		gitRepo, err = newRepo(path)
		if err != nil {
			return nil, err
		}
		return gitRepo, nil
	}
	return gitRepo, nil
}

func Git(cmdstring string, args ...string) (*exec.Cmd, *bytes.Buffer, *bytes.Buffer) {
	cmdArgs := make([]string, 1)
	cmdArgs[0] = cmdstring
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command(gitCmd, cmdArgs...)
	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
	cmd.Stdout, cmd.Stderr = stdout, stderr
	return cmd, stdout, stderr
}

func newRepo(path string) (repo *Repo, err error) {
	repo = &Repo{}
	stat, err := os.Stat(path)
	if err != nil {
		return
	}

	if !stat.IsDir() {
		err = errors.Errorf("%s is not dir", path)
		return
	}

	if stat, err = os.Stat(filepath.Join(path, ".git", "config")); err != nil {
		cmd, _, _ := Git("init", path)
		cmd.Dir = path
		if err = cmd.Run(); err != nil {
			return
		}
	}

	err = repo.Init(path)

	return
}