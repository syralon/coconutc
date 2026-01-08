package entschema

import (
	"text/template"

	"entgo.io/ent/entc/gen"
)

const schemaTemplateText = `package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
)

// {{ .Name }} holds the schema definition for the {{ .Name }} entity.
type {{ .Name }} struct {
	ent.Schema
}

// Fields of the {{ .Name }}.
func ({{ .Name }}) Fields() []ent.Field {
    return []ent.Field{
		{{ range .Fields }}field.{{ .Type.String | pascal }}("{{.Name}}"),
		{{ end }}
    }
}

// Edges of the {{ .Name }}.
func ({{ .Name }}) Edges() []ent.Edge {
	return nil
}
`

const (
	schemaFieldsTemplateText = `
func ({{ .Name }}) Fields() []ent.Field {
    return []ent.Field{
		{{ range .Fields }}field.{{ .Type.String | pascal }}("{{.Name}}"),
		{{ end }}
    }
}`
)

var (
	schemaTemplate       *template.Template
	schemaFieldsTemplate *template.Template
)

func init() {
	{
		schemaTemplate = template.New("schema").Funcs(gen.Funcs)
		var err error
		schemaTemplate, err = schemaTemplate.Parse(schemaTemplateText)
		if err != nil {
			panic(err)
		}
	}
	{
		schemaFieldsTemplate = template.New("schema").Funcs(gen.Funcs)
		var err error
		schemaFieldsTemplate, err = schemaFieldsTemplate.Parse(schemaFieldsTemplateText)
		if err != nil {
			panic(err)
		}
	}
}
