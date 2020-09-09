package procexec

import (
	"errors"
	"os/exec"
	"time"

	"github.com/google/uuid"

	"github.com/smallnest/ringbuffer"
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

	proc := Proc{
		ID:      uuid.New().String(),
		Cmd:     cmd,
		Args:    args,
		out:     ringbuffer.New(4096),
		cmd:     c,
		running: false,
		exited:  false,
		created: time.Now(),
	}

	c.Stdout = proc.out
	c.Stderr = proc.out

	return &proc
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
