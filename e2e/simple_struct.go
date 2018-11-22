package e2e

//go:generate fm5 -t SimpleStruct -fo -r Interface
type SimpleStruct struct {
	id   string
	name string
}

func (s *SimpleStruct) Run() string {
	return s.name
}

type Interface interface {
	Run() string
}

