package cli

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/teatak/pipe/common"
)

func SetPid(pid int) {
	pidString := []byte(strconv.Itoa(pid))
	_ = os.MkdirAll(common.BinDir+"/run", 0755)
	_ = ioutil.WriteFile(common.BinDir+"/run/pipe.pid", pidString, 0666)
}

func GetPid() int {
	content, err := ioutil.ReadFile(common.BinDir + "/run/pipe.pid")
	if err != nil {
		return 0
	} else {
		pid, _ := strconv.Atoi(strings.Trim(string(content), "\n"))
		if _, find := ProcessExist(pid); find {
			return pid
		} else {
			return 0
		}
	}
}

func ProcessExist(pid int) (*os.Process, bool) {
	process, err := os.FindProcess(pid)
	if err != nil {
		return nil, false
	} else {
		err := process.Signal(syscall.Signal(0))
		if err != nil {
			return nil, false
		}
	}
	return process, true
}
