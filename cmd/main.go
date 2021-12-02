package main

import (
	"fmt"
	"os"

	"github.com/teatak/pipe/cli"
	"github.com/teatak/pipe/cmd/cli/daem"
	"github.com/teatak/pipe/cmd/cli/quit"
	"github.com/teatak/pipe/cmd/cli/reload"
	"github.com/teatak/pipe/cmd/cli/run"
	"github.com/teatak/pipe/cmd/cli/version"
	"github.com/teatak/pipe/common"
)

var Commands cli.Commands

func init() {
	Commands = cli.Commands{
		"version": version.New(common.Version),
		"daem":    daem.New(),
		"reload":  reload.New(),
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
 dP""Y88b   d8P''Y8b  .dP""Y88b  o888  o888  .dP""Y88b   d8P''Y8b  .dP""Y88b  
      ]8P' 888    888       ]8P'  888   888        ]8P' 888    888       ]8P' 
    .d8P'  888    888     .d8P'   888   888      .d8P'  888    888     .d8P'  
  .dP'     888    888   .dP'      888   888    .dP'     888    888   .dP'     
.oP     .o '88b  d88' .oP     .o  888   888  .oP     .o '88b  d88' .oP     .o 
8888888888  'Y8bd8P'  8888888888 o888o o888o 8888888888  'Y8bd8P'  8888888888 
`
