package codegen

import (
	"bytes"
	"go/format"
	"text/template"

	"github.com/deref/exo/inflect"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

type Root struct {
	Package string
	Module
}

type Module struct {
	Interfaces []Interface `hcl:"interface,block"`
	Structs    []Struct    `hcl:"struct,block"`
}

type Interface struct {
	Name    string   `hcl:"name,label"`
	Doc     *string  `hcl:"doc"`
	Methods []Method `hcl:"method,block"`
}

type Method struct {
	Name   string  `hcl:"name,label"`
	Doc    *string `hcl:"doc"`
	Input  []Field `hcl:"input,block"`
	Output []Field `hcl:"output,block"`
}

type Struct struct {
	Name   string  `hcl:"name,label"`
	Doc    *string `hcl:"doc"`
	Fields []Field `hcl:"field,block"`
}

type Field struct {
	Name     string  `hcl:"name,label"`
	Doc      *string `hcl:"doc"`
	Type     string  `hcl:"type,label"`
	Required *bool   `hcl:"required"`
	Nullable *bool   `hcl:"nullable"`
}

func ParseFile(filename string) (*Module, error) {
	var module Module
	if err := hclsimple.DecodeFile(filename, nil, &module); err != nil {
		return nil, err
	}
	err := Validate(module)
	return &module, err
}

func Validate(mod Module) error {
	// TODO: Validate no duplicate names.
	// TODO: All type references.
	return nil
}

func Generate(root *Root) ([]byte, error) {
	tmpl := template.Must(
		template.New("module").
			Funcs(map[string]interface{}{
				"tick":   func() string { return "`" },
				"public": inflect.KebabToPublic,
				"js":     inflect.KebabToJSVar,
			}).
			Parse(moduleTemplate),
	)
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, root); err != nil {
		return nil, err
	}
	bs := buf.Bytes()
	formatted, err := format.Source(bs)
	if err != nil {
		return bs, nil
	}
	return formatted, nil

}

var moduleTemplate = `
// Generated file. DO NOT EDIT.

package {{.Package}}

{{- define "doc" -}}
{{if .Doc}}// {{.Doc}}
{{end}}{{end}}

{{- define "fields" -}}
{{- range . }}
{{template "doc" . -}}
	{{.Name|public}} {{.Type}} {{tick}}json:"{{.Name|js}}"{{tick}}
{{- end}}{{end}}

import (
	"context"
	"net/http"

	josh "github.com/deref/exo/josh/server"
)

{{range .Interfaces}}
{{template "doc" . -}}
type {{.Name|public}} interface {
{{- range .Methods}}
{{template "doc" . -}}
	{{.Name|public}}(context.Context, *{{.Name|public}}Input) (*{{.Name|public}}Output, error)
{{- end}}
}

{{range .Methods}}
type {{.Name|public}}Input struct {
{{template "fields" .Input}}
}

type {{.Name|public}}Output struct {
{{template "fields" .Output}}
}
{{end}}

func New{{.Name|public}}Mux(prefix string, iface {{.Name|public}}) *http.ServeMux {
	b := josh.NewMuxBuilder(prefix)
	Build{{.Name|public}}Mux(b, iface)
	return b.Mux()
}

func Build{{.Name|public}}Mux(b *josh.MuxBuilder, iface {{.Name|public}}) {
{{- range .Methods}}
	b.AddMethod("{{.Name}}", iface.{{.Name|public}})
{{- end}}
}
{{end}}

{{range .Structs}}
{{template "doc" . -}}
type {{.Name|public}} struct {
{{template "fields" .Fields}}
}
{{end}}
`
