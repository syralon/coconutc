package entschema

import (
	"bytes"
	"errors"
	"go/format"
	"io"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
)

type FieldType string

const (
	Bool    FieldType = "bool"
	Int     FieldType = "int"
	Int8    FieldType = "int8"
	Int16   FieldType = "int16"
	Int32   FieldType = "int32"
	Int64   FieldType = "int64"
	Uint    FieldType = "uint"
	Uint8   FieldType = "uint8"
	Uint16  FieldType = "uint16"
	Uint32  FieldType = "uint32"
	Uint64  FieldType = "uint64"
	Float32 FieldType = "float32"
	Float64 FieldType = "float64"
	Time    FieldType = "time"
	String  FieldType = "string"
)

func (t FieldType) String() string {
	return string(t)
}

type Field struct {
	Name string
	Type FieldType
}

type Fields []*Field

func (f *Fields) add(field *Field) {
	for _, v := range *f {
		if v.Name == field.Name {
			*v = *field
			return
		}
	}
	*f = append(*f, field)
}

type Schema struct {
	Name      string
	Fields    Fields
	Overwrite bool
}

var (
	validateSingle = regexp.MustCompile("^[A-Z]\\w+$")
)

func Parse(s string) (*Schema, error) {
	if validateSingle.MatchString(s) {
		return &Schema{Name: s, Fields: defaultFields()}, nil
	}
	n1 := strings.Index(s, "(")
	n2 := strings.Index(s, ")")
	if n1 < 0 || n2 < 0 {
		return nil, errors.New("invalid schema")
	}
	name := s[:n1]
	if name == "" {
		return nil, errors.New("invalid schema")
	}
	var overwrite = name[len(name)-1] == '-'
	if overwrite {
		name = name[:len(name)-1]
	}
	fields := strings.Split(s[n1+1:n2], ",")
	schema := &Schema{Name: name, Fields: defaultFields(), Overwrite: overwrite}
	for _, field := range fields {
		ty := strings.Split(field, ":")
		f := &Field{Name: strcase.ToSnake(ty[0])}
		if len(ty) > 1 {
			f.Type = FieldType(ty[1])
		} else if strings.HasSuffix(f.Name, "_id") {
			f.Type = Int
		} else if strings.HasSuffix(f.Name, "_at") {
			f.Type = Time
		} else {
			f.Type = String
		}
		schema.Fields.add(f)
	}
	return schema, nil
}

func Parses(args []string) ([]*Schema, error) {
	var schemas = make([]*Schema, 0, len(args))
	for _, arg := range args {
		s, err := Parse(arg)
		if err != nil {
			return nil, err
		}
		schemas = append(schemas, s)
	}
	return schemas, nil
}

func defaultFields() Fields {
	return Fields{
		{Name: "id", Type: Int},
		{Name: "created_at", Type: Time},
		{Name: "updated_at", Type: Time},
	}
}

func (s *Schema) WriteFile(target string) error {
	filename := path.Join(target, strings.ToLower(s.Name)+".go")
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return s.writeFile(filename)
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	data, err = format.Source(data)
	if err != nil {
		return err
	}
	data, err = s.replace(data)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (s *Schema) replace(data []byte) ([]byte, error) {
	var reg, err = regexp.Compile("func \\(*" + s.Name + "\\) Fields\\(\\) \\[]ent.Field {")
	if err != nil {
		return nil, err
	}
	// func ({{ .Name }}) Fields() []ent.Field {
	lines := strings.Split(string(data), "\n")
	var start, end int
	for n, line := range lines {
		if reg.MatchString(line) {
			start = n
			continue
		}
		if start == 0 {
			continue
		}
		if line[0] == '}' {
			end = n
			break
		}
	}
	if start == 0 || end <= start {
		return nil, errors.New("invalid ent schema file")
	}

	buf := bytes.NewBuffer(nil)
	for _, line := range lines[:start] {
		buf.WriteString(line + "\n")
	}

	if err = schemaFieldsTemplate.Execute(buf, s); err != nil {
		return nil, err
	}

	for _, line := range lines[end+1:] {
		buf.WriteString(line + "\n")
	}
	return format.Source(buf.Bytes())
}

func (s *Schema) writeFile(filename string) error {
	_ = os.MkdirAll(path.Dir(filename), os.ModePerm)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	return s.Execute(file)
}

func (s *Schema) Execute(w io.Writer) error {
	return schemaTemplate.Execute(w, s)
}
