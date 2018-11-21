package fm5

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
)

//go:generate ./bin/generr -t typeNotFoundWithGivenName -i
type typeNotFoundWithGivenName interface {
	TypeNotFoundWithGivenName() (name string)
}

func Parse(r io.Reader, typename string) (string, *ast.TypeSpec, error) {
	src, err := ioutil.ReadAll(r)
	if err != nil {
		return "", nil, errors.Wrap(err, "parse file readAll failed")
	}
	f, err := parser.ParseFile(token.NewFileSet(), "", string(src), parser.ParseComments)
	// pp.Print(f)
	if err != nil {
		return "", nil, errors.Wrap(err, "parse file with filename is failed")
	}

	var pkgName string
	var ts *ast.TypeSpec
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.File:
			pkgName = x.Name.Name
		case *ast.TypeSpec:
			if _, ok := x.Type.(*ast.StructType); !ok {
				return true
			}
			if x.Name.Name == typename {
				ts = x
				return false
			}
		default:
			return true
		}
		return true
	})

	if ts == nil {
		return "", nil, &TypeNotFoundWithGivenName{Name: typename}
	}

	return pkgName, ts, nil
}
