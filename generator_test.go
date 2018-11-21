package fm5

import (
	"testing"

	"bytes"

	"go/format"
	"go/token"

	"github.com/andreyvit/diff"
)

func TestGenerator_AppendPackage(t *testing.T) {
	g := NewGenerator("main", "", nil)
	g.AppendPackage()

	exp := "package main\n"
	var buf bytes.Buffer
	format.Node(&buf, token.NewFileSet(), g.f)
	if act := buf.String(); act != exp {
		t.Error(diff.LineDiff(exp, act))
	}
}

func TestGenerator_AppendDefaultFactory(t *testing.T) {

	helper := func(t *testing.T, methodname, typename, src, exp string) {
		t.Helper()
		n, s, err := Parse(bytes.NewBufferString(src), typename)
		if err != nil {
			t.Fatal(err)
		}
		g := NewGenerator(n, typename, s)
		g.AppendPackage()
		if err := g.AppendDefaultFactory(methodname); err != nil {
			t.Fatal(err)
		}
		var buf bytes.Buffer
		format.Node(&buf, token.NewFileSet(), g.f)
		if act := buf.String(); act != exp {
			t.Error(diff.LineDiff(exp, act))
		}
	}

	t.Run("has no member", func(t *testing.T) {
		src := `package main

			type A struct {}
			`

		exp := `package main

func NewA() *A {
	return &A{}
}
`
		helper(t, "NewA", "A", src, exp)
	})

	t.Run("has int member", func(t *testing.T) {
		src := `package main

			type A struct {
				id int
			}`

		exp := `package main

func NewA(id int) *A {
	return &A{id: id}
}
`

		helper(t, "NewA", "A", src, exp)
	})

	t.Run("return multi value value", func(t *testing.T) {
		src := `package main

			type A struct {
				id int
				name string
			}`

		exp := `package main

func NewA(id int, name string) *A {
	return &A{id: id, name: name}
}
`
		helper(t, "NewA", "A", src, exp)
	})
}

func TestGenerator_AppendFunctionalOptionType(t *testing.T) {

	helper := func(t *testing.T, typename, methodname, src, exp string) {
		t.Helper()
		n, s, err := Parse(bytes.NewBufferString(src), typename)
		if err != nil {
			t.Fatal(err)
		}
		g := NewGenerator(n, typename, s)
		g.AppendPackage()
		if err := g.AppendFunctionalOptionType(methodname); err != nil {
			t.Fatal(err)
		}
		var buf bytes.Buffer
		format.Node(&buf, token.NewFileSet(), g.f)
		if act := buf.String(); act != exp {
			t.Error(diff.LineDiff(exp, act))
		}
	}

	t.Run("has no member", func(t *testing.T) {
		src := `package main
type A struct {}
`

		exp := `package main

type AOption func(*A)

func AWithOptions(opts ...AOption) *A {
	i := &A{}
	for _, o := range opts {
		o(i)
	}

	return i
}
`
		helper(t, "A", "AWithOptions", src, exp)
	})
}
