package run

import (
	"flag"
	"fmt"
	"strings"
)

const synopsis = "Run Pipe"
const help = `Usage: run

  Run pipe service
`

const infoServerPrefix = "[INFO] pipe.server: "

type cmd struct {
	flags *flag.FlagSet
}

func New() *cmd {
	c := &cmd{}
	c.init()
	return c
}

func (c *cmd) init() {
	c.flags = flag.NewFlagSet("run", flag.ContinueOnError)
	c.flags.Usage = func() {
		fmt.Println(c.Help())
	}
}

func (c *cmd) Synopsis() string {
	return synopsis
}

func (c *cmd) Help() string {
	return strings.TrimSpace(help)
}
