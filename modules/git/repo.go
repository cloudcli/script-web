package git

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	// "github.com/syncthing/syncthing/lib/sync"
	"path"
	"time"

	"github.com/pkg/errors"
)

type SHA string
type Repo struct {
	commitLock sync.Mutex
	GitDir     string
	WorkDir    string
	Config     map[string]string
}

func (r *Repo) Init(path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		return err
	}

	if stat, err = os.Stat(filepath.Join(path, ".git", "config")); err != nil {
		return err
	}
	r.WorkDir = path

	return nil
}

func (r *Repo) Git(cmdstring string, args ...string) (stdout *bytes.Buffer, err error) {
	cmd, stdout, stderr := Git(cmdstring, args...)
	cmd.Dir = r.WorkDir

	if err := cmd.Start(); err != nil {
		return nil, errors.Errorf("cmd start error: %s", err.Error())
	}

	timeout := 60 * time.Second

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(timeout):
		if cmd.Process != nil && cmd.ProcessState != nil && !cmd.ProcessState.Exited() {
			if killerr := cmd.Process.Kill(); killerr != nil {
				err = fmt.Errorf("fail to kill process: %v", killerr)
				return
			}
		}

		<-done
		err = fmt.Errorf("TIME_OUT")
		return
	case err = <-done:
	}

	if err != nil {
		err = errors.Errorf("%s:%s", err.Error(), stderr.String())
	}

	return
}

func (r *Repo) CatFile(fullpath string, ref string) (io.Reader, error) {
	stdout, err := r.Git("ls-tree", "--full-tree", ref+" "+fullpath)
	if err != nil {
		return nil, err
	}

	if stdout.Len() == 0 {
		return nil, fmt.Errorf("%s is not present in %s", fullpath, r.WorkDir)
	}

	parts := strings.Split(stdout.String(), " ")
	if parts[1] != "blob" {
		return nil, fmt.Errorf("%s is not present in %s", fullpath, r.WorkDir)
	}

	shaname := strings.Split(parts[2], "\t")
	catout, err := r.Git("cat-file", "blob", shaname[0])

	return catout, err
}

func (r *Repo) ListTree(fullpath string, ref string) (lsout io.Reader, err error) {
	fmt.Println(fullpath)
	stat, err := os.Stat(path.Join(r.WorkDir, fullpath))
	if err != nil {
		lsout = nil
		return
	}

	if stat.IsDir() {
		fullpath = fullpath + "/"
	}

	lsout, err = r.Git("ls-tree", "--full-tree", ref, fullpath)
	return
}

func (r *Repo) CommitFile(name string, dir string, content []byte, user string, email string) (ref string, diff string, err error) {
	stat, err := os.Stat(path.Join(r.WorkDir, dir))
	if err != nil {
		if os.IsNotExist(err) {
			return "", "", errors.Errorf("%s is not directory", dir)
		}
		return "", "", err
	}

	if !stat.IsDir() {
		return "", "", errors.Errorf("%s is not directory", dir)
	}

	fullpath := path.Join(r.WorkDir, dir, name)
	_, err = os.Stat(fullpath)
	if os.IsNotExist(err) {
		file, err := os.Create(fullpath)
		if err != nil {
			return "", "", errors.New("create file error")
		}

		defer file.Close()

		_, err = file.Write(content)
		if err != nil {
			return "", "", errors.New("write to file error")
		}

		_, err = r.Git("add", "-A")
		if err != nil {
			return "", "", errors.New("add error")
		}

		diffout, err := r.Git("diff", "HEAD", "--", path.Join(dir, name))
		if err != nil {
			return "", "", errors.Errorf("diff error: %s", err.Error())
		}

		diff = diffout.String()

		_, err = r.Git("commit", fmt.Sprintf("-am\"Create File: %s %s\"", dir, name))

		if err != nil {
			return "", "", errors.New("commit error " + err.Error())
		}

		logout, err := r.Git("log", "-1")
		if err != nil {
			return "", "", errors.New("log error " + err.Error())
		}

		hash, err := getCommitHash(logout.String())
		if err != nil {
			return "", "", err
		}

		return hash, diff, nil
	} else {
		return "", "", errors.Errorf("%s is already exist", path.Join(dir, name))
	}

	return "", "", nil
}

//func (r *Repo) Diff(hash string, filepath) {
//
//}

func (r *Repo) Status() {

}

func getCommitHash(raw string) (string, error) {
	lines := strings.Split(raw, "\n")
	if len(lines) < 3 {
		return "", errors.New("cant find commit info")
	}
	line := lines[0]
	fields := strings.Fields(line)
	return fields[1], nil
}
