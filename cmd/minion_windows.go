package cmd

import (
	"fmt"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"

	"github.com/ushmodin/avaxo2/internal/minion"
	"github.com/ushmodin/avaxo2/internal/settings"

)

type minionWinService struct {
	elog debug.Log
}

func (m *minionWinService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	changes <- svc.Status{State: svc.StartPending}
	if err := settings.InitSettings(); err != nil {
		m.elog.Error(1, fmt.Sprintf("Error while init setttings: %v", err))
		return false, 1
	}
	server, err := minion.NewServer(
		settings.MinionSettings.Listen,
		settings.MinionSettings.Keyfile,
		settings.MinionSettings.Certfile,
		settings.MinionSettings.Cafile,
	)
	if err != nil {
		m.elog.Error(1, fmt.Sprintf("Error while create minion server: %v", err))
		return false, 1
	}
	serverChan := make(chan error, 1)
	changes <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}

	go func() {
		m.elog.Info(1, fmt.Sprintf("Starting minion server"))
		serverChan <- server.Run()
	}()

	loop := true
	for loop {
		select {
		case err := <-serverChan:
			if err != nil {
				m.elog.Error(1, fmt.Sprintf("Minion server's error %v", err))
			}
			loop = false
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
				time.Sleep(100 * time.Millisecond)
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				loop = false
				if err := server.Stop(); err != nil {
					m.elog.Error(1, fmt.Sprintf("Minion server's error %v", err))
				}
			default:
				m.elog.Error(1, fmt.Sprintf("unexpected control request #%d", c))
			}
		}
	}

	m.elog.Info(1, fmt.Sprintf("Stop"))
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
	err = svc.Run(name, &minionWinService{
		elog: elog,
	})
	if err != nil {
		elog.Error(1, fmt.Sprintf("%s service failed: %v", name, err))
		return err
	}
	elog.Info(1, fmt.Sprintf("%s service stopped", name))
	return nil
}
