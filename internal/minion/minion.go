package minion

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/ushmodin/avaxo2/internal/model"
	"github.com/ushmodin/avaxo2/internal/procexec"
)

// Minion provide minion commands
type Minion struct {
	procsMux sync.Mutex
	procs    map[string]*procexec.Proc
}

// DirItem filesystem directory item information

// NewMinion create new minion
func NewMinion() *Minion {
	return &Minion{
		procs: make(map[string]*procexec.Proc),
	}
}

// ReadDir get directory listing
func (minion *Minion) ReadDir(path string) ([]model.DirItem, error) {
	log.Printf("LS: %s\n", path)
	d, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	names, err := d.Readdirnames(-1)
	d.Close()
	res := make([]model.DirItem, 0, len(names))
	for _, filename := range names {
		dirItem := model.DirItem{
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
	log.Printf("Return %d files\n", len(res))
	return res, nil
}

func (minion *Minion) GetFile(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

func (minion *Minion) PutFile(path string, mode os.FileMode, reader io.Reader) error {
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

func (minion *Minion) Exec(cmd string, args ...string) (string, error) {
	log.Printf("Exec: %s %v", cmd, args)
	proc := procexec.NewProc(cmd, args...)
	if err := proc.Start(); err != nil {
		return "", err
	}
	minion.procsMux.Lock()
	minion.procs[proc.ID] = proc
	minion.procsMux.Unlock()
	return proc.ID, nil
}

func (minion *Minion) ProcInfo(id string) (model.ProcInfo, error) {
	minion.procsMux.Lock()
	proc, ok := minion.procs[id]
	minion.procsMux.Unlock()

	if !ok {
		return model.ProcInfo{}, fmt.Errorf("Proc %s not found", id)
	}
	info := model.ProcInfo{
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

func (minion *Minion) ProcKill(id string) error {
	minion.procsMux.Lock()
	proc, ok := minion.procs[id]
	minion.procsMux.Unlock()

	if !ok {
		return fmt.Errorf("Proc %s not found", id)
	}

	return proc.Kill()
}

func (minion *Minion) ProcPs() []model.ProcPsItem {
	ps := make([]model.ProcPsItem, len(minion.procs))
	i := 0
	for k, v := range minion.procs {
		ps[i].ID = k
		ps[i].Cmd = v.Cmd
		ps[i].Args = v.Args
		ps[i].Created = v.Created().Format(time.RFC3339)
		ps[i].Exited = v.Exited()
		i++
	}
	return ps
}

func (minion *Minion) ProcTail(id string, w io.Writer) error {
	minion.procsMux.Lock()
	proc, ok := minion.procs[id]
	minion.procsMux.Unlock()

	if !ok {
		return fmt.Errorf("Proc %s not found", id)
	}
	proc.Tail(w)
	return nil
}
