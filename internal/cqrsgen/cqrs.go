package cqrsgen

import (
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/influx6/faux/fmtwriter"

	"github.com/gokit/cqrskit/internal/static"
	"github.com/influx6/moz/ast"
	"github.com/influx6/moz/gen"
)

type MethodEventPair struct {
	Name     string
	TypeName string
	Argument ast.ArgType
	Method   ast.FuncDeclaration
	Def      ast.FunctionDefinition
	Event    ast.StructDeclaration
}

// ESCQRSGen generates scaffolding code for
func ESCQRSGen(toPackage string, an ast.AnnotationDeclaration, str ast.StructDeclaration, declr ast.PackageDeclaration, pkg ast.Package) ([]gen.WriteDirective, error) {
	structName := str.Object.Name.Name
	structNameLower := strings.ToLower(structName)

	//packageName := fmt.Sprintf("%scqrs", structNameLower)

	var methodImports []gen.ImportItemDeclr

	var pairs []MethodEventPair
	methods, _ := declr.MethodFor(str.Object.Name.Name)
	for _, method := range methods {
		if !strings.HasPrefix(method.FuncName, "Handle") {
			continue
		}

		if method.HasAnnotation("@escqrs-skip") {
			log.Printf("Skipping: Function definition for %q in %q\n", method.FuncName, str.Name)
			continue
		}

		def, err := ast.GetFunctionDefinitionFromDeclaration(method, &declr)
		if err != nil {
			log.Printf("Skipping: Unable to get function definition for %q in %q\n", method.FuncName, str.Name)
			continue
		}

		if def.TotalArgs() == 0 {
			log.Printf("Skipping: Unable to use function %q with no arguments in %q as event handler\n", method.FuncName, str.Name)
			continue
		}

		if def.TotalArgs() != 1 {
			log.Printf("Skipping: Unable to use function %q with more than 1 arguments in %q as event handler\n", method.FuncName, str.Name)
			continue
		}

		// Get the argument detail for the argument.
		eventArg := def.Args[0]

		if eventArg.StructObject == nil {
			log.Printf("Skipping: Unable to use function %q with non-struct argument %q in %q\n", method.FuncName, eventArg.ExType, str.Name)
		}

		var typeName string

		// Get name of Event from Handle{{Event}} method
		eventStructName := strings.TrimPrefix(method.FuncName, "Handle")

		var ok bool
		var eventStruct ast.StructDeclaration
		if eventArg.Package != declr.Package {
			pkg, ok = declr.ImportedPackageFor(eventArg.Package)
			if !ok {
				log.Printf("Skipping: Unable to find package for function %q with argument  %q in %q\n", method.FuncName, eventStructName, str.Name)
				continue
			}

			methodImports = append(methodImports, gen.Import(pkg.Path, eventArg.Package))

			structName := strings.TrimPrefix(strings.TrimPrefix(eventArg.ExType, eventArg.Package), ".")
			eventStruct, ok = pkg.StructFor(structName)
			if !ok {
				log.Printf("Skipping: Unable to find struct %q in package %q for function %q in %q\n", structName, pkg.Path, method.FuncName, str.Name)
				continue
			}

			typeName = eventArg.ExType
		} else {
			// Search for event struct declared in package.
			eventStruct, ok = pkg.StructFor(eventStructName)
			if !ok {
				log.Printf("Skipping: Unable to find %q struct definition for method %+q on %q\n", eventStructName, method.FuncName, str.Name)
				continue
			}

			typeName = eventStruct.Name
		}

		if !def.HasReturnType("error") {
			log.Printf("Skipping: Method %q must return an error after receiving %q in %q\n", method.FuncName, eventStructName, str.Name)
			continue
		}

		pairs = append(pairs, MethodEventPair{
			Def:      def,
			Method:   method,
			Argument: eventArg,
			TypeName: typeName,
			Event:    eventStruct,
			Name:     eventStructName,
		})
	}

	methodImports = append(methodImports, gen.Import("context", ""))
	methodImports = append(methodImports, gen.Import("github.com/gokit/cqrskit", ""))

	readWriteRepo := gen.Package(
		gen.Name(declr.Package),
		gen.Imports(methodImports...),
		gen.Block(
			gen.SourceTextWith(
				"cqrskit:appliers",
				string(static.MustReadFile("appliers.tml", true)),
				gen.ToTemplateFuncs(
					ast.ASTTemplatFuncs,
					template.FuncMap{},
				),
				struct {
					Str   ast.StructDeclaration
					Pkg   ast.PackageDeclaration
					An    ast.AnnotationDeclaration
					Pairs []MethodEventPair
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
			FileName: fmt.Sprintf("%s.cqrs.go", structNameLower),
			Writer:   fmtwriter.New(readWriteRepo, true, true),
		},
	}

	return documents, nil
}
