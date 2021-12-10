package sections

//load section in load.go file
type Backend struct {
	Name   string   `yaml:"name,omitempty"`
	Mode   string   `yaml:"mode,omitempty"`
	Riff   string   `yaml:"riff,omitempty"`
	Server []string `yaml:"server,omitempty"`
}

type backendSection []*Backend

func (s backendSection) SectionName() string {
	return "backend"
}

var Backends = backendSection{}
