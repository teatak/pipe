//go:build windows
// +build windows

package reload

import (
	"flag"
	"fmt"

	"github.com/teatak/pipe/cmd/cli"
	"github.com/teatak/pipe/cmd/cli/daem"
	"github.com/teatak/pipe/cmd/cli/quit"
)

const help = `Usage: reload

  Reload pipe config
  
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
		if _, find := cli.ProcessExist(pid); find {
			//quit
			q := quit.New()
			q.Run([]string{})
			//run
			s := daem.New()
			s.Run([]string{})
			fmt.Println("reload pipe success")
		}
	}
}

func (c *cmd) Synopsis() string {
	return "Reload Pipe config"
}

func (c *cmd) Help() string {
	return help
}
