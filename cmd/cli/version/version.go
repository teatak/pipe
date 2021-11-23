package version

import (
	"fmt"

	"github.com/teatak/pipe/common"
)

type cmd struct {
	version string
}

func New(version string) *cmd {
	return &cmd{version: version}
}

func (c *cmd) Run(_ []string) int {
	fmt.Printf(common.Name+" version %s, %s build %s-%s\n", c.version, common.Type, common.GitBranch, common.GitSha)
	return 0
}

func (c *cmd) Synopsis() string {
	return "Prints the Pipe version"
}

func (c *cmd) Help() string {
	return ""
}
