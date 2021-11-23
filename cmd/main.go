package main

import (
	"fmt"
	"os"

	"github.com/teatak/pipe/cli"
	"github.com/teatak/pipe/cmd/cli/daem"
	"github.com/teatak/pipe/cmd/cli/quit"
	"github.com/teatak/pipe/cmd/cli/run"
	"github.com/teatak/pipe/cmd/cli/version"
	"github.com/teatak/pipe/common"
)

var Commands cli.Commands

func init() {
	Commands = cli.Commands{
		"version": version.New(common.Version),
		"daem":    daem.New(),
		"quit":    quit.New(),
		"run":     run.New(),
	}
}
func main() {
	args := os.Args[1:]
	for _, arg := range args {
		if arg == "cheers" {
			fmt.Println(cheers)
			return
		}
		if arg == "--" {
			break
		}

		if arg == "-v" || arg == "--version" {
			args = []string{"version"}
			break
		}
	}

	c := cli.NewCLI("pipe", common.Version)
	c.Args = args
	c.Commands = Commands
	exitCode, err := c.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
	}

	os.Exit(exitCode)
}

const cheers = `
 222222222222222         000000000       1111111        888888888
2:::::::::::::::22     00:::::::::00    1::::::1      88:::::::::88
2::::::222222:::::2  00:::::::::::::00 1:::::::1    88:::::::::::::88
2222222     2:::::2 0:::::::000:::::::0111:::::1   8::::::88888::::::8
            2:::::2 0::::::0   0::::::0   1::::1   8:::::8     8:::::8
            2:::::2 0:::::0     0:::::0   1::::1   8:::::8     8:::::8
         2222::::2  0:::::0     0:::::0   1::::1    8:::::88888:::::8
    22222::::::22   0:::::0 000 0:::::0   1::::l     8:::::::::::::8
  22::::::::222     0:::::0 000 0:::::0   1::::l    8:::::88888:::::8
 2:::::22222        0:::::0     0:::::0   1::::l   8:::::8     8:::::8
2:::::2             0:::::0     0:::::0   1::::l   8:::::8     8:::::8
2:::::2             0::::::0   0::::::0   1::::l   8:::::8     8:::::8
2:::::2       2222220:::::::000:::::::0111::::::1118::::::88888::::::8
2::::::2222222:::::2 00:::::::::::::00 1::::::::::1 88:::::::::::::88
2::::::::::::::::::2   00:::::::::00   1::::::::::1   88:::::::::88
22222222222222222222     000000000     111111111111     888888888
`
