package main

import (
	"os"
	"path/filepath"

	"github.com/gokit/cqrskit/cqrs"

	"github.com/influx6/faux/flags"
	"github.com/influx6/faux/metrics"
	"github.com/influx6/faux/metrics/custom"
	"github.com/influx6/moz/ast"
)

func main() {
	flags.Run("cqrskit", flags.Command{
		Name:      "generate",
		ShortDesc: "Generates a cqrs+es like API for struct types.",
		Desc:      "Generates a EventSourcing+CQRS API from target structs to create ES based services with.",
		Action: func(ctx flags.Context) error {
			force, _ := ctx.GetBool("force")
			dest, _ := ctx.GetString("dest")
			target, _ := ctx.GetString("target")
			verbose, _ := ctx.GetBool("verbose")

			logs := metrics.New()

			if verbose {
				logs = metrics.New(custom.StackDisplay(os.Stderr))
			}

			currentdir, err := os.Getwd()
			if err != nil {
				return err
			}

			if !filepath.IsAbs(dest) {
				dest = filepath.Join(currentdir, dest)
			}

			currentdir = filepath.Join(currentdir, target)

			generators := ast.NewAnnotationRegistryWith(logs)
			generators.Register("@escqrs", cqrs.ESCQRSGen)

			res, err := ast.ParseAnnotations(logs, currentdir)
			if err != nil {
				return err
			}

			return ast.SimplyParse(dest, logs, generators, force, res...)
		},
		Flags: []flags.Flag{
			&flags.BoolFlag{
				Name: "verbose",
				Desc: "verbose logs all operations out to console.",
			},
			&flags.BoolFlag{
				Name: "force",
				Desc: "force regeneration of packages annotation directives.",
			},
			&flags.StringFlag{
				Name:    "dest",
				Default: "./",
				Desc:    "relative destination for package",
			},
			&flags.StringFlag{
				Name:    "target",
				Default: "./",
				Desc:    "-target=./ defines relative path of target for code gen",
			},
		},
	})
}
