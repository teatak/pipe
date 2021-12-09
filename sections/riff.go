package sections

//load section in load.go file
type riff struct {
	Url string `yaml:"url,omitempty"`
}

func (s riff) SectionName() string {
	return "riff"
}

var Riff = riff{}
