package gen

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/hashicorp/go-hclog"
	"github.com/krystal/go-katapult/apischema"
	"mvdan.cc/gofumpt/format"
)

type Generator struct {
	PkgName           string
	OutputDir         string
	SchemaIncludePath string
	SchemaExcludePath string
	SchemaFiles       []string

	Logger hclog.Logger
}

func (g *Generator) newFile(packageName string) *jen.File {
	f := jen.NewFile(packageName)
	f.ImportName("github.com/krystal/go-katapult", "katapult")
	f.ImportName("github.com/krystal/go-katapult/core", "core")

	return f
}

func (g *Generator) loadSchema(filename string) (*apischema.Schema, error) {
	schemaFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	schema := &apischema.Schema{}
	err = json.NewDecoder(schemaFile).Decode(schema)
	if err != nil {
		return nil, err
	}

	return schema, nil
}

func (g *Generator) katapult(name string) *jen.Statement {
	base := &jen.Statement{}
	if strings.HasPrefix(name, "*") {
		base = jen.Id("*")
		name = name[1:]
	}

	if g.PkgName == "katapult" {
		return base.Id(name)
	}

	return base.Qual("github.com/krystal/go-katapult", name)
}

func (g *Generator) core(name string) *jen.Statement {
	base := &jen.Statement{}
	if strings.HasPrefix(name, "*") {
		base = jen.Id("*")
		name = name[1:]
	}

	if g.PkgName == "core" {
		return base.Id(name)
	}

	return base.Qual("github.com/krystal/go-katapult/core", name)
}

func (g *Generator) render(f *jen.File) ([]byte, error) {
	var langVersion string
	out, err := exec.Command(
		"go", "list", "-m", "-f", "{{.GoVersion}}",
	).Output()
	out = bytes.TrimSpace(out)
	if err == nil && len(out) > 0 {
		langVersion = string(out)
	}

	var buf bytes.Buffer
	err = f.Render(&buf)
	if err != nil {
		return nil, err
	}

	b, err := format.Source(buf.Bytes(), format.Options{
		LangVersion: langVersion,
		ExtraRules:  true,
	})
	if err != nil {
		return nil, err
	}

	return b, nil
}
