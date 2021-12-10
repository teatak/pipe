package sections

//load section in load.go file
type pipe struct {
	Pid string `yaml:"pid,omitempty"`
}

func (s pipe) SectionName() string {
	return "pipe"
}

var Pipe = pipe{}
