package files

import (
	"io"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/filesystem"
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

	fs := memfs.New()
	storage := filesystem.NewStorage(fs, cache.NewObjectLRU(1*cache.GiByte))
	_, err = git.Clone(storage, fs, &git.CloneOptions{
		URL:  repo,
		Auth: auth,
	})
	if err != nil {
		return nil, err
	}

	var filesData [][]byte = make([][]byte, len(files))

	for i, file := range files {
		file, err := fs.Open(file)
		if err != nil {
			return nil, err
		}

		b, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}

		filesData[i] = b
	}

	return filesData, nil
}
