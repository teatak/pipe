package sections

type server struct {
	Listen string    `yaml:"listen,omitempty"`
	Domain []*domain `yaml:"domain,omitempty"`
}

type domain struct {
	Name     string      `yaml:"name,omitempty"`
	CertFile string      `yaml:"certFile,omitempty"`
	KeyFile  string      `yaml:"keyFile,omitempty"`
	Location []*location `yaml:"location,omitempty"`
}

type location struct {
	Path string `yaml:"path,omitempty"`
	To   string `yaml:"to,omitempty"`
}

type serverSection []*server

func (s *serverSection) SectionName() string {
	return "servers"
}

var Servers = &serverSection{}
