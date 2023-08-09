package main

import "text/template"

//nolint:gochecknoglobals
var (
	prologTemplate = template.Must(template.New("prolog").Parse(`
package {{ . }}

import (
	"fmt"
	"strings"
)
`))
	constTemplate = template.Must(template.New("const").Parse(`
const (
	{{- range $c := . }}
	{{ $c.Name }} {{ $c.Type }} = {{ $c.Value }}
	{{- end }}
)
`))

	messageTemplate = template.Must(template.New("message").Parse(`
{{range .Structs}}
{{.DocText "// "}}
type {{.NameNative}} struct {
{{range .Fields}}{{.DocText "\t// "}}
	{{.NameNative}}	{{if .TypeList}}[]{{end}}{{if .TypeRef}}*{{end}}{{.TypeNative}}{{range .TagAdd}} {{.}}{{end}}
{{end}}}
{{end}}}
`))

	initTemplate = template.Must(template.New("init").Parse(`
func init() {
	{{- range $i, $sd := . }}
	parse.AssertUpToDate(&{{ $sd.TableVar }}.s, new({{ $sd.Type }}))
	{{- end }}
}
`))
)
