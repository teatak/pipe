package sections

import (
	"github.com/teatak/config"
	"github.com/teatak/pipe/common"
)

func init() {
	Load()
}

func Load() {
	config.LoadConfig(common.BinDir + "/config/app.yml")
	config.Load(&Backends)
	config.Load(&Servers)
}
