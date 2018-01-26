package cqrs

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/influx6/faux/fmtwriter"

	"github.com/gokit/cqrskit/static"
	"github.com/influx6/moz/ast"
	"github.com/influx6/moz/gen"
)

type CommandEventPair struct {
	Command ast.StructDeclaration
	Event   ast.StructDeclaration
}

// ESCQRSGen generates scaffolding code for
func ESCQRSGen(toPackage string, an ast.AnnotationDeclaration, str ast.StructDeclaration, declr ast.PackageDeclaration, pkg ast.Package) ([]gen.WriteDirective, error) {
	structName := str.Object.Name.Name
	structNameLower := strings.ToLower(structName)

	packageName := fmt.Sprintf("%scqrs", structNameLower)

	var pairs []CommandEventPair

	for _, annotation := range str.Annotations {
		if annotation.Name != "CommandEvent" {
			continue
		}

		// If it does not have associated Command key, skip...
		if _, ok := annotation.Params["Command"]; !ok {
			continue
		}

		// If it does not have associated Event key, skip...
		if _, ok := annotation.Params["Event"]; !ok {
			continue
		}

		// Find struct declaration for event name else skip.
		eventStructName := annotation.Param("Event")
		eventStruct, ok := pkg.StructFor(eventStructName)
		if !ok {
			continue
		}

		// Find struct declaration for command name else skip.
		commandStructName := annotation.Param("Command")
		commandStruct, ok := pkg.StructFor(commandStructName)
		if !ok {
			continue
		}

		pairs = append(pairs, CommandEventPair{
			Command: commandStruct,
			Event:   eventStruct,
		})
	}

	readWriteRepo := gen.Package(
		gen.Name(packageName),
		gen.Imports(
			gen.Import("context", ""),
			gen.Import(declr.Path, ""),
			gen.Import("github.com/gokit/cqrskit/internal/cqrs", ""),
		),
		gen.Block(
			gen.SourceTextWith(
				"cqrskit:read-write-repository",
				string(static.MustReadFile("read-write-repository.tml", true)),
				gen.ToTemplateFuncs(
					ast.ASTTemplatFuncs,
					template.FuncMap{},
				),
				struct {
					Str   ast.StructDeclaration
					Pkg   ast.PackageDeclaration
					An    ast.AnnotationDeclaration
					Pairs []CommandEventPair
				}{
					An:    an,
					Str:   str,
					Pkg:   declr,
					Pairs: pairs,
				},
			),
		),
	)

	documents := []gen.WriteDirective{
		{
			FileName: fmt.Sprintf("%s-read-write-repo.cqrs.go", structNameLower),
			Writer:   fmtwriter.New(readWriteRepo, true, true),
			Dir:      packageName,
		},
	}

	return documents, nil
}
