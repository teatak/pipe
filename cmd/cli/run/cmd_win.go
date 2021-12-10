//go:build windows
// +build windows

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
	pid := os.Getegid()
	s.Logger.Printf(infoServerPrefix+"start server %v\n", pid)
	setPid(pid)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for {
			sig := <-sigs
			fmt.Println()
			s.Logger.Printf(infoServerPrefix+"get signal %v\n", sig)
			s.Shutdown()
		}
	}()
	<-s.ShutdownCh
	return 0
}
