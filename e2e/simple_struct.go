package e2e

import "log"

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
