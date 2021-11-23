package sections

import "github.com/teatak/config"

type server struct {
	Listen string    `yaml:"listen,omitempty"`
	Domain []*Domain `yaml:"domain,omitempty"`
}

type Domain struct {
	Name     string      `yaml:"name,omitempty"`
	CertFile string      `yaml:"certFile,omitempty"`
	KeyFile  string      `yaml:"keyFile,omitempty"`
	Location []*Location `yaml:"location,omitempty"`
}

type Location struct {
	Path string `yaml:"path,omitempty"`
	To   string `yaml:"to,omitempty"`
}

type serverSection []*server

func (s *serverSection) SectionName() string {
	return "servers"
}

var Servers = serverSection{}

func init() {
	config.Load(&Servers)
}
