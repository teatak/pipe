//go:build !windows
// +build !windows

package run

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/teatak/pipe/server"
)

func (c *cmd) Run(args []string) int {
	if err := c.flags.Parse(args); err != nil {
		return 1
	}
	s, err := server.NewServer()
	if err != nil {
		fmt.Println(err)
		return 1
	}
	s.Logger.Printf(infoServerPrefix+"start server %v\n", os.Getpid())
	sigs := make(chan os.Signal, 10)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
	go func() {
		for {
			sig := <-sigs
			s.Logger.Printf(infoServerPrefix+"get signal %v\n", sig)
			if sig == syscall.SIGUSR2 {
				//grace reload
				s.Reload()
			} else {
				s.Shutdown()
			}
		}
	}()
	<-s.ShutdownCh
	return 0
}
