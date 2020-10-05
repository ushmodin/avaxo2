package procexec

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
	"sync/atomic"
	"time"

	"github.com/smallnest/ringbuffer"
)

var (
	procID = int32(0)
)

type Proc struct {
	ID      string
	Cmd     string
	Args    []string
	out     *ringbuffer.RingBuffer
	cmd     *exec.Cmd
	running bool
	exited  bool
	err     error
	created time.Time
}

func NewProc(cmd string, args ...string) *Proc {
	c := exec.Command(cmd, args...)
	out := ringbuffer.New(4096)
	c.Stdout = out
	c.Stderr = out

	return &Proc{
		ID:      fmt.Sprintf("%d", atomic.AddInt32(&procID, 1)),
		Cmd:     cmd,
		Args:    args,
		out:     out,
		cmd:     c,
		running: false,
		exited:  false,
		created: time.Now(),
	}
}

func (proc *Proc) Start() error {
	if proc.running {
		return errors.New("already runngin")
	}
	if err := proc.cmd.Start(); err != nil {
		return err
	}
	proc.running = true
	go func() {
		if err := proc.cmd.Wait(); err != nil {
			proc.err = err
		}
		proc.exited = true
	}()
	return nil
}

func (proc *Proc) Out() []byte {
	return proc.out.Bytes()
}

func (proc *Proc) Exited() bool {
	if !proc.running {
		return false
	}
	return proc.exited
}

func (proc *Proc) ExitCode() (int, error) {
	if !proc.running {
		return 0, errors.New("not running")
	}

	if !proc.exited {
		return 0, errors.New("still running")
	}
	return proc.cmd.ProcessState.ExitCode(), nil
}

func (proc *Proc) Pid() (int, error) {
	if !proc.running {
		return 0, errors.New("not running")
	}
	return proc.cmd.ProcessState.Pid(), nil
}

func (proc *Proc) Kill() error {
	if !proc.running {
		return errors.New("not running")
	}
	return proc.cmd.Process.Kill()
}

func (proc *Proc) Running() bool {
	return proc.running
}

func (proc *Proc) Created() time.Time {
	return proc.created
}

func (proc *Proc) Tail(w io.Writer) error {
	for !proc.Exited() {
		b := make([]byte, 1024)
		n, err := proc.out.Read(b)

		if err == ringbuffer.ErrIsEmpty {
			time.Sleep(100 * time.Millisecond)
		} else if err != nil {
			return err
		} else if n != 0 {
			if _, err := io.CopyN(w, bytes.NewBuffer(b), int64(n)); err != nil {
				return err
			}
		}
	}
	log.Printf("Finished\n")
	return nil
}
