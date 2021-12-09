package sections

//load section in load.go file
type Server struct {
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
	Return string   `yaml:"return,omitempty"`
}

type serverSection []*Server

func (s serverSection) SectionName() string {
	return "server"
}

var Servers = serverSection{}
