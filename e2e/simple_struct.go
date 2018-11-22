package e2e

import "log"

//go:generate fm5 -t SimpleStruct -fo
type SimpleStruct struct {
	id   string
	name string
}

func (*SimpleStruct) Run() {
	log.Println("test")
}

type Interface interface {
	Run()
}

//go:generate fm5 -t Structure -fo -r Interface
type Structure struct {
	id   int
	name string
}
