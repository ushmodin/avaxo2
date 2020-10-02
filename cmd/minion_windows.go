package cmd

import (
	"fmt"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
)

type minionWinService struct {
	elog debug.Log
}

func (m *minionWinService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	changes <- svc.Status{State: svc.StartPending}
	time.Sleep(100 * time.Millisecond)
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	loop := true
	for loop {
		c := <-r
		switch c.Cmd {
		case svc.Interrogate:
			changes <- c.CurrentStatus
			time.Sleep(100 * time.Millisecond)
			changes <- c.CurrentStatus
			m.elog.Info(1, fmt.Sprintf("Interrogate"))
		case svc.Stop, svc.Shutdown:
			loop = false
			m.elog.Info(1, fmt.Sprintf("Stop"))
		default:
			m.elog.Error(1, fmt.Sprintf("unexpected control request #%d", c))
		}
	}

	changes <- svc.Status{State: svc.StopPending}
	return
}

func runWinSrv(name string) error {
	elog, err := eventlog.Open(name)
	if err != nil {
		return err
	}
	defer elog.Close()

	elog.Info(1, fmt.Sprintf("starting %s service", name))
	err = svc.Run(name, &minionWinService{})
	if err != nil {
		elog.Error(1, fmt.Sprintf("%s service failed: %v", name, err))
		return err
	}
	elog.Info(1, fmt.Sprintf("%s service stopped", name))
	return nil
}
