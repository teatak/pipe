package reload

import (
	"flag"
	"fmt"
	"syscall"

	"github.com/teatak/pipe/cmd/cli"
)

const help = `Usage: pipe reload

Options:

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
	c.flags = flag.NewFlagSet("reload", flag.ContinueOnError)

	c.flags.Usage = func() {
		fmt.Println(c.Help())
	}
}

func (c *cmd) Run(args []string) int {
	if err := c.flags.Parse(args); err != nil {
		return 1
	}
	c.Reload()
	return 0
}

func (c *cmd) Reload() {
	pid := cli.GetPid()
	if pid == 0 {
		fmt.Println("can't find pipe")
	} else {
		if p, find := cli.ProcessExist(pid); find {
			err := p.Signal(syscall.SIGUSR2)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("reload pipe success")
			}
		}
	}
}

func (c *cmd) Synopsis() string {
	return "Reload Pipe config"
}

func (c *cmd) Help() string {
	return help
}
