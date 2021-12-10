package run

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/teatak/pipe/sections"
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

func setPid(pid int) {
	if sections.Pipe.Pid != "" {
		pidFile := sections.Pipe.Pid
		pidString := []byte(strconv.Itoa(pid))
		os.MkdirAll(filepath.Dir(pidFile), 0755)
		ioutil.WriteFile(pidFile, pidString, 0666)
	}
}
