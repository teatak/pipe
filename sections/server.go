package sections

//load section in load.go file
type server struct {
	Listen string    `yaml:"listen,omitempty"`
	SSL    bool      `yaml:"ssl,omitempty"`
	Domain []*domain `yaml:"domain,omitempty"`
}

type domain struct {
	Name     string      `yaml:"name,omitempty"`
	CertFile string      `yaml:"certFile,omitempty"`
	KeyFile  string      `yaml:"keyFile,omitempty"`
	Return   string      `yaml:"return,omitempty"`
	Location []*location `yaml:"location,omitempty"`
}

type location struct {
	Path string `yaml:"path,omitempty"`
	To   string `yaml:"to,omitempty"`
}

type serverSection []*server

func (s *serverSection) SectionName() string {
	return "server"
}

var Server = &serverSection{}
