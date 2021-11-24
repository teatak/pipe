package sections

import "github.com/teatak/config"

func init() {
	Load()
}

func Load() {
	config.LoadConfig()
	config.Load(Riff)
	config.Load(Endpoint)
	config.Load(Servers)
}
