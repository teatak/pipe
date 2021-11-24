package sections

type endpoint struct {
	Method  string   `yaml:"method,omitempty"`
	Servers []string `yaml:"servers,omitempty"`
}

type endpointSection map[string]*endpoint

func (s *endpointSection) SectionName() string {
	return "endpoint"
}

var Endpoint = &endpointSection{}
