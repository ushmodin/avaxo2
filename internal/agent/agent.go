package agent

import (
	"io"
	"os"
	"time"
)

// Agent provide agent commands
type Agent struct {
}

// DirItem filesystem directory item information
type DirItem struct {
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	Modified string `json:"modified"`
	IsDir    bool   `json:"isDir"`
	Error    string `json:"error"`
}

// NewAgent create new agent
func NewAgent() *Agent {
	return &Agent{}
}

// ReadDir get directory listing
func (agent *Agent) ReadDir(path string) ([]DirItem, error) {
	d, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	names, err := d.Readdirnames(-1)
	d.Close()
	res := make([]DirItem, 0, len(names))
	for _, filename := range names {
		dirItem := DirItem{
			Name: filename,
			Size: 0,
		}
		fi, err := os.Lstat(path + "/" + filename)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			dirItem.Error = err.Error()
		} else {
			dirItem.IsDir = fi.IsDir()
			if !fi.IsDir() {
				dirItem.Size = fi.Size()
			}
			dirItem.Modified = fi.ModTime().Format(time.RFC3339)
		}
		res = append(res, dirItem)
	}
	return res, nil
}

func (agent *Agent) GetFile(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

func (agent *Agent) PutFile(path string, mode os.FileMode, reader io.Reader) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := io.Copy(f, reader); err != nil {
		return err
	}
	return nil
}
