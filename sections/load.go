package sections

import (
	"github.com/teatak/config"
)

func init() {
	Load()
}

func Load() {
	config.Load(&Backends)
	config.Load(&Servers)
}
