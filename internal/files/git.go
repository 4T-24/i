package files

import (
	"io"
	"os"
	"path"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/sirupsen/logrus"
)

func GetFiles(repo string, files []string) ([][]byte, error) {
	key, err := ssh.NewPublicKeysFromFile("git", "/key/deploy_key", "")
	if err != nil {
		logrus.Warn("no ssh key found, trying without")
	}

	var auth transport.AuthMethod = nil
	if key != nil {
		auth = key
	}

	tempDir, err := os.MkdirTemp(os.TempDir(), "git_*")
	if err != nil {
		return nil, err
	}

	_, err = git.PlainClone(tempDir, false, &git.CloneOptions{
		URL:  repo,
		Auth: auth,
	})
	if err != nil {
		return nil, err
	}

	var filesData [][]byte = make([][]byte, len(files))

	for i, file := range files {
		path := path.Join(tempDir, file)

		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		b, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}

		filesData[i] = b
	}

	os.RemoveAll(tempDir)

	return filesData, nil
}
