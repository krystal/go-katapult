package gen

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/hashicorp/go-hclog"
	"github.com/krystal/go-katapult/apischema"
)

func (g *Generator) Errors() error {
	idMatcher, err := regexp.Compile(g.SchemaPath)
	if err != nil {
		return err
	}

	f := g.newFile(g.PkgName)
	f.Comment(
		"Code generated by github.com/krystal/go-katapult/tools/codegen. " +
			"DO NOT EDIT.",
	).Line()

	for _, filename := range g.SchemaFiles {
		schema, err2 := g.loadSchema(filename)
		if err2 != nil {
			return err2
		}

		var sortedErrors []*apischema.Error
		for _, e := range schema.Errors {
			sortedErrors = append(sortedErrors, e)
		}

		sort.SliceStable(sortedErrors, func(i, j int) bool {
			a := g.errVarName(sortedErrors[i]) + ":" + sortedErrors[i].ID
			b := g.errVarName(sortedErrors[j]) + ":" + sortedErrors[j].ID

			return a < b
		})

		errorCodes := map[string]bool{}
		var errorObjects []*apischema.Error
		for _, e := range sortedErrors {
			_, exists := errorCodes[e.Code]
			if !exists && idMatcher.MatchString(e.ID) {
				errorObjects = append(errorObjects, e)
				errorCodes[e.Code] = true
			}
		}

		for _, e := range errorObjects {
			err2 = g.errVar(f, e)
			if err2 != nil {
				return err2
			}
		}

		for _, e := range errorObjects {
			err3 := g.errStruct(f, e)
			if err3 != nil {
				return err3
			}

			err3 = g.errNewStructFunc(f, e)
			if err3 != nil {
				return err3
			}

			if len(e.Fields) > 0 {
				err3 = g.errStructDetail(f, schema, e)
				if err3 != nil {
					return err3
				}
			}
			f.Line()
		}

		err2 = g.errCastResponseFunc(f, errorObjects)
		if err2 != nil {
			return err2
		}
	}

	b, err := g.render(f)
	if err != nil {
		return err
	}

	outdir, err := filepath.Abs(g.OutputDir)
	if err != nil {
		return err
	}

	g.Logger.Info("writing errors_generated.go",
		"dir", outdir,
		"size", hclog.Fmt("%d bytes", len(b)),
	)
	//nolint:gosec
	err = ioutil.WriteFile(
		filepath.Join(outdir, "errors_generated.go"), b, 0o644,
	)
	if err != nil {
		return err
	}

	return nil
}

func (g *Generator) errCastResponseFunc(
	f *jen.File,
	errs []*apischema.Error,
) error {
	cases := []jen.Code{}
	for _, e := range errs {
		funcName := "New" + g.errStructName(e)
		cases = append(cases, jen.Case(jen.Lit(e.Code)).Return(
			jen.Id(funcName).Call(jen.Id("theError")),
		))
	}

	cases = append(cases, jen.Default().Return(jen.Id("theError")))

	f.Comment("castResponseError casts a *katapult.ResponseError to a more " +
		"specific type based on the error's Code value.")
	f.Func().Id("castResponseError").Params(
		jen.Id("theError").Add(g.katapult("*ResponseError")),
	).Error().Block(
		jen.Switch(jen.Id("theError.Code")).Block(cases...),
	)

	return nil
}

func (g *Generator) errVarName(e *apischema.Error) string {
	return "Err" + strings.TrimSuffix(filepath.Base(e.ID), "Errors")
}

func (g *Generator) errVar(f *jen.File, e *apischema.Error) error {
	var parent *jen.Statement
	switch e.HTTPStatus {
	case http.StatusBadRequest:
		parent = g.katapult("ErrBadRequest")
	case http.StatusUnauthorized:
		parent = g.katapult("ErrUnauthorized")
	case http.StatusForbidden:
		parent = g.katapult("ErrForbidden")
	case http.StatusNotFound:
		parent = g.katapult("ErrResourceNotFound")
	case http.StatusNotAcceptable:
		parent = g.katapult("ErrNotAcceptable")
	case http.StatusConflict:
		parent = g.katapult("ErrConflict")
	case http.StatusUnprocessableEntity:
		parent = g.katapult("ErrUnprocessableEntity")
	case http.StatusTooManyRequests:
		parent = g.katapult("ErrTooManyRequests")
	case http.StatusInternalServerError:
		parent = g.katapult("ErrInternalServerError")
	case http.StatusBadGateway:
		parent = g.katapult("ErrBadGateway")
	case http.StatusServiceUnavailable:
		parent = g.katapult("ErrServiceUnavailable")
	case http.StatusGatewayTimeout:
		parent = g.katapult("ErrGatewayTimeout")
	default:
		parent = g.katapult("ErrUnknown")
	}

	f.Var().Id(g.errVarName(e)).Op("=").Qual("fmt", "Errorf").Call(
		jen.Lit("%w: "+e.Code),
		parent,
	)

	return nil
}

func (g *Generator) errStructName(e *apischema.Error) string {
	name := filepath.Base(e.ID)
	if !strings.HasSuffix(name, "Error") {
		name += "Error"
	}

	return name
}

func (g *Generator) errStruct(f *jen.File, e *apischema.Error) error {
	name := g.errStructName(e)
	detailName := g.errStructDetailName(e)

	var detailField *jen.Statement
	if len(e.Fields) > 0 {
		detailField = jen.Id("Detail").Id("*" + detailName).Tag(
			map[string]string{"json": "detail,omitempty"},
		)
	}

	f.Comment(name + ":")
	f.Comment(e.Description)
	f.Type().Id(name).Struct(
		g.katapult("CommonError"),
		detailField,
	)

	return nil
}

func (g *Generator) errNewStructFunc(f *jen.File, e *apischema.Error) error {
	name := g.errStructName(e)
	detailName := g.errStructDetailName(e)

	funcBody := []jen.Code{
		jen.Return(
			jen.Id("&" + name).Values(jen.Dict{
				jen.Line().Id("CommonError"): g.katapult("NewCommonError").
					Call(
						jen.Line().Id(g.errVarName(e)),
						jen.Line().Lit(e.Code),
						jen.Line().Id("theError.Description").Id(",").Line(),
					),
			}),
		),
	}

	if len(e.Fields) > 0 {
		funcBody = []jen.Code{
			jen.Id("detail").Op(":=").Id("&" + detailName).Values(),
			jen.Id("err").Op(":=").Qual("encoding/json", "Unmarshal").Call(
				jen.Id("theError.Detail"),
				jen.Id("detail"),
			),
			jen.If(jen.Id("err").Op("!=").Nil()).Block(
				jen.Id("detail").Op("=").Nil(),
			).Line(),
			jen.Return(
				jen.Id("&" + name).Values(jen.Dict{
					jen.Id("CommonError"): g.katapult("NewCommonError").Call(
						jen.Line().Id(g.errVarName(e)),
						jen.Line().Lit(e.Code),
						jen.Line().Id("theError.Description").Id(",").Line(),
					),
					jen.Id("Detail"): jen.Id("detail"),
				}),
			),
		}
	}

	f.Func().Id("New" + name).Params(
		jen.Id("theError").Add(g.katapult("*ResponseError")),
	).Id("*" + name).Block(funcBody...)

	return nil
}

func (g *Generator) errStructDetailName(e *apischema.Error) string {
	return g.errStructName(e) + "Detail"
}

func (g *Generator) errStructDetail(
	f *jen.File,
	s *apischema.Schema,
	e *apischema.Error,
) error {
	detailName := g.errStructDetailName(e)

	fields, err := g.structFields(s, e.Fields)
	if err != nil {
		return err
	}

	f.Type().Id(detailName).Struct(fields...)

	return nil
}

func (g *Generator) structFields(
	s *apischema.Schema,
	fields []*apischema.Field,
) ([]jen.Code, error) {
	statements := []jen.Code{}
	for _, field := range fields {
		statement, err := g.structField(s, field)
		if err != nil {
			return nil, err
		}

		statements = append(statements, statement)
	}

	return statements, nil
}

func (g *Generator) fieldName(f *apischema.Field) string {
	return strings.TrimSuffix(filepath.Base(f.ID), "Field")
}

func (g *Generator) structField(
	s *apischema.Schema,
	f *apischema.Field,
) (jen.Code, error) {
	name := g.fieldName(f)
	tag := jen.Tag(map[string]string{
		"json": f.Name + ",omitempty",
	})

	base := &jen.Statement{}
	sliceBase := &jen.Statement{}
	if f.Array {
		base = base.Id("[]")
		sliceBase = sliceBase.Id("[]")
	}
	if f.Null {
		base = base.Id("*")
	}

	switch f.Type {
	case "Apia/Scalars/Boolean", "Rapid/Scalars/Boolean":
		return jen.Id(name).Add(base.Bool().Add(tag)), nil
	case "Apia/Scalars/Decimal", "Rapid/Scalars/Decimal":
		return jen.Id(name).Add(base.Float64().Add(tag)), nil
	case "Apia/Scalars/Integer", "Rapid/Scalars/Integer":
		return jen.Id(name).Add(base.Int().Add(tag)), nil
	case "Apia/Scalars/String", "Rapid/Scalars/String":
		return jen.Id(name).Add(base.String().Add(tag)), nil
	case "Apia/Scalars/UnixTime", "Rapid/Scalars/UnixTime":
		return jen.Id(name).Add(sliceBase.Id("*").Qual(
			"github.com/augurysys/timestamp", "Timestamp",
		)).Add(tag), nil
	case "CoreAPI/Objects/TrashObject":
		return jen.Id(name).Add(
			sliceBase.Add(g.core("*TrashObject")).Add(tag),
		), nil
	default:
		obj, ok := s.Objects[f.Type]
		if !ok {
			return nil, fmt.Errorf(
				"field type %s object not found for field %s",
				f.Type, f.ID,
			)
		}

		g.Logger.Warn(
			"generating anonymous struct",
			"type", f.Type, "field", f.ID,
		)

		fields, err := g.structFields(s, obj.Fields)
		if err != nil {
			return nil, err
		}

		return jen.Id(name).Add(sliceBase.Struct(fields...).Add(tag)), nil
	}
}
