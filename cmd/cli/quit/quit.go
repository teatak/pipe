package quit

import (
	"flag"
	"fmt"
	"syscall"
	"time"

	"github.com/teatak/pipe/cmd/cli"
)

const help = `Usage: pipe quit
`

type cmd struct {
	flags *flag.FlagSet
}

func New() *cmd {
	c := &cmd{}
	c.init()
	return c
}

func (c *cmd) init() {
	c.flags = flag.NewFlagSet("quit", flag.ContinueOnError)

	c.flags.Usage = func() {
		fmt.Println(c.Help())
	}
}

func (c *cmd) Run(args []string) int {
	if err := c.flags.Parse(args); err != nil {
		return 1
	}
	c.Quit()
	return 0
}

func (c *cmd) Quit() {
	pid := cli.GetPid()
	if pid == 0 {
		fmt.Println("can't find pipe")
	} else {
		if p, find := cli.ProcessExist(pid); find {
			err := p.Signal(syscall.SIGINT)
			if err != nil {
				fmt.Println(err)
			} else {
				quitStop := make(chan bool)
				go func() {
					for {
						if pid := cli.GetPid(); pid == 0 {
							quitStop <- true
							break
						}
						time.Sleep(1 * time.Second)
					}
				}()
				<-quitStop
				fmt.Println("quit pipe success")
			}
		}
	}
}

func (c *cmd) Synopsis() string {
	return "Quit Pipe"
}

func (c *cmd) Help() string {
	return help
}
