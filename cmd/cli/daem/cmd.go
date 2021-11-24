package daem

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/teatak/pipe/cmd/cli"
	"github.com/teatak/pipe/common"
)

const synopsis = "Run Pipe as service "
const help = `Usage: daem

  Run pipe as daemon service
  
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
	c.flags = flag.NewFlagSet("daem", flag.ContinueOnError)

	c.flags.Usage = func() {
		fmt.Println(c.Help())
	}
}

func (c *cmd) Run(args []string) int {
	if err := c.Start(args); err != nil {
		fmt.Println(err)
		return 1
	}
	return 0
}

func (c *cmd) Start(args []string) error {
	if cli.GetPid() != 0 {
		return fmt.Errorf("%s is already running", "pipe")
	}
	command := c.resoveCommand(common.BinDir + "/pipe")
	dir, _ := filepath.Abs(filepath.Dir(command))

	newArgs := append([]string{}, "run")
	newArgs = append(newArgs, args...)
	args = newArgs
	cmd := exec.Command(command, args...)
	cmd.Dir = dir

	out := common.MakeFile(common.BinDir + "/logs/stdout.log")
	cmd.Stdout = out
	cmd.Stderr = out

	err := cmd.Start()
	if err != nil {
		return err
	} else {
		cli.SetPid(cmd.Process.Pid)
		fmt.Println("start pipe success")
	}
	return nil
}

func (c *cmd) resoveCommand(path string) string {
	if filepath.IsAbs(path) {
		return path
	} else {
		if strings.HasPrefix(path, "."+string(os.PathSeparator)) {
			return common.BinDir + path[1:]
		} else {
			return path
		}
	}
}

func (c *cmd) Synopsis() string {
	return synopsis
}

func (c *cmd) Help() string {
	return strings.TrimSpace(help)
}
