package agent

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/ushmodin/avaxo2/internal/procexec"
)

// Agent provide agent commands
type Agent struct {
	procsMux sync.Mutex
	procs    map[string]*procexec.Proc
}

// DirItem filesystem directory item information
type DirItem struct {
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	Modified string `json:"modified"`
	IsDir    bool   `json:"isDir"`
	Error    string `json:"error"`
}

type ProcInfo struct {
	Cmd      string   `json:"cmd"`
	Args     []string `json:"args"`
	Exited   bool     `json:"exited"`
	ExitCode int      `json:"exitCode"`
	Out      []byte   `json:"out"`
	Created  string   `json:"created"`
}

// NewAgent create new agent
func NewAgent() *Agent {
	return &Agent{
		procs: make(map[string]*procexec.Proc),
	}
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

func (agent *Agent) Exec(cmd string, args ...string) (string, error) {
	proc := procexec.NewProc(cmd, args...)
	if err := proc.Start(); err != nil {
		return "", err
	}
	agent.procsMux.Lock()
	agent.procs[proc.ID] = proc
	agent.procsMux.Unlock()
	return proc.ID, nil
}

func (agent *Agent) ProcInfo(id string) (ProcInfo, error) {
	agent.procsMux.Lock()
	proc, ok := agent.procs[id]
	agent.procsMux.Unlock()

	if !ok {
		return ProcInfo{}, fmt.Errorf("Proc %s not found", id)
	}
	info := ProcInfo{
		Cmd:      proc.Cmd,
		Args:     proc.Args,
		Exited:   proc.Exited(),
		Out:      proc.Out(),
		Created:  proc.Created().Format(time.RFC3339),
		ExitCode: 0,
	}

	if proc.Running() && proc.Exited() {
		ec, _ := proc.ExitCode()
		info.ExitCode = ec
	}

	return info, nil
}

func (agent *Agent) ProcKill(id string) error {
	agent.procsMux.Lock()
	proc, ok := agent.procs[id]
	agent.procsMux.Unlock()

	if !ok {
		return fmt.Errorf("Proc %s not found", id)
	}

	return proc.Kill()
}
