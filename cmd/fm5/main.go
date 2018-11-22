package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/akito0107/fm5"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "fm5"
	app.Usage = "factory method generator"
	app.UsageText = "fm5 [OPTIONS]"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "type, t",
			Usage: "struct name (required)",
		},
		cli.BoolFlag{
			Name:  "dryrun",
			Usage: "dryrun (default=false)",
		},
		cli.BoolTFlag{
			Name:  "factory-method, fm",
			Usage: "generate default factory method(default=true)",
		},
		cli.StringFlag{
			Name:  "factory-method-name, fmn",
			Usage: "factory method name(default=New + $typename)",
		},
		cli.BoolFlag{
			Name:  "functional-option, fo",
			Usage: "generate functional option patterns methods(default=false)",
		},
		cli.StringFlag{
			Name:  "functional-option-name, fon",
			Usage: "functional option method name(New + $typename + Options)",
		},
		cli.StringFlag{
			Name:  "return-typename, r",
			Usage: "return typename (if present, applied this type, otherwise, using pointer type of given type)",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(ctx *cli.Context) error {
	typename := ctx.String("type")
	dryrun := ctx.Bool("dryrun")
	fm := ctx.Bool("factory-method")
	fo := ctx.Bool("functional-option")
	o := ctx.String("return-typename")
	if typename == "" {
		return errors.New("type is required")
	}
	fmname := ctx.String("factory-method-name")
	if fmname == "" {
		fmname = "New" + strings.Title(typename)
	}
	fmoname := ctx.String("functional-option-name")
	if fmoname == "" {
		fmoname = "New" + strings.Title(typename) + "Options"
	}

	var filenames []string
	fileinfos, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range fileinfos {
		if strings.HasSuffix(f.Name(), ".go") && !strings.HasSuffix(f.Name(), "_test.go") {
			filenames = append(filenames, f.Name())
		}
	}

	for _, f := range filenames {
		ok, err := generate(f, typename, fmname, fmoname, o, fm, fo, dryrun)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
	}

	return errors.Errorf("typename %s is not found", typename)
}

func generate(filename, typename, fmn, fmon, outtype string, fm, fo, dryrun bool) (bool, error) {
	r, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer r.Close()
	pkgName, ts, err := fm5.Parse(r, typename)
	if ok, _ := fm5.IsTypeNotFoundWithGivenName(err); ok {
		return false, nil
	} else if err != nil {
		return false, err
	}
	g := fm5.NewGenerator(pkgName, typename, ts)
	g.AppendPackage()
	if fm {
		if err := g.AppendDefaultFactory(fmn, outtype); err != nil {
			return false, err
		}
	}

	if fo {
		if err := g.AppendFunctionalOptionType(fmon, outtype); err != nil {
			return false, err
		}
		if err := g.AppendFunctionalOptions(); err != nil {
			return false, err
		}
	}

	if dryrun {
		g.Out(os.Stdout)
	} else {
		filename := fmt.Sprintf("%s_fm.go", ts.Name.Name)
		if _, err := os.Stat(filename); !os.IsNotExist(err) {
			os.Remove(filename)
		}
		f, err := os.Create(filename)
		if err != nil {
			return false, err
		}
		defer f.Close()
		if err := g.Out(f); err != nil {
			return false, err
		}
	}

	return true, nil
}
