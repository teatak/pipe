package sections

//load section in load.go file
type server struct {
	Listen string    `yaml:"listen,omitempty"`
	SSL    bool      `yaml:"ssl,omitempty"`
	Domain []*Domain `yaml:"domain,omitempty"`
}

type Domain struct {
	Name     string      `yaml:"name,omitempty"`
	CertFile string      `yaml:"certFile,omitempty"`
	KeyFile  string      `yaml:"keyFile,omitempty"`
	Location []*Location `yaml:"location,omitempty"`
}

type Location struct {
	Path   string   `yaml:"path,omitempty"`
	Header []string `yaml:"header,omitempty"`
	Proxy  string   `yaml:"proxy,omitempty"`
	Return string   `yaml:"return,omitempty"`
}

type serverSection []*server

func (s *serverSection) SectionName() string {
	return "server"
}

var Server = &serverSection{}
