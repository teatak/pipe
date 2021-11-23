package run

import (
	"flag"
	"fmt"
	"strings"
)

const synopsis = "Run Pipe"
const help = `Usage: run [options]

  Run pipe service

Options:

  -name       Node name
  -dc         DataCenter name
  -http       Http address of riff (-http 127.0.0.1:8610)
  -rpc        RPC address of riff (-rpc [::]:8630)
  -join       Join RPC address (-join 192.168.1.1:8630,192.168.1.2:8630,192.168.1.3:8630)
`

const infoServerPrefix = "[INFO] pipe.server: "

type cmd struct {
	flags *flag.FlagSet
	name  string
	dc    string
	http  string
	dns   string
	rpc   string
	join  string
}

func New() *cmd {
	c := &cmd{}
	c.init()
	return c
}

func (c *cmd) init() {
	c.flags = flag.NewFlagSet("start", flag.ContinueOnError)
	c.flags.StringVar(&c.http, "http", "", "usage")
	c.flags.StringVar(&c.dns, "dns", "", "usage")
	c.flags.StringVar(&c.rpc, "rpc", "", "usage")
	c.flags.StringVar(&c.name, "name", "", "usage")
	c.flags.StringVar(&c.join, "join", "", "usage")
	c.flags.StringVar(&c.dc, "dc", "", "usage")

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
