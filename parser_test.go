package fmgo

import (
	"bytes"
	"testing"
)

func TestParse(t *testing.T) {

	t.Run("can parse package name", func(t *testing.T) {
		src := `package main
type A struct{}
`
		s, _, err := Parse(bytes.NewBufferString(src), "A")
		if err != nil {
			t.Fatal(err)
		}
		if s != "main" {
			t.Errorf("must be main but %s", s)
		}
	})

	t.Run("can parse struct with given type name", func(t *testing.T) {
		src := `package main
type A struct{}
`
		_, node, err := Parse(bytes.NewBufferString(src), "A")
		if err != nil {
			t.Fatal(err)
		}
		if node.Name.Name != "A" {
			t.Errorf("type name must be A but %s", node.Name.Name)
		}
	})

	t.Run("return err if given typename is not found", func(t *testing.T) {
		src := `package main
type A struct{}
`
		_, _, err := Parse(bytes.NewBufferString(src), "xxx")
		if ok, _ := IsTypeNotFoundWithGivenName(err); !ok {
			t.Fatal(err)
		}
	})

}
