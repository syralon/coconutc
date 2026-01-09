package entservice

import (
	"bufio"
	"bytes"
	"embed"
	"errors"
	"fmt"
	"go/format"
	"io"
	"path"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
)

//go:embed templates
var fs embed.FS
var templates *template.Template

func init() {
	funcs := map[string]any{
		"toLower":   strings.ToLower,
		"camel":     strcase.ToLowerCamel,
		"basepath":  path.Base,
		"toPackage": toPackage,
	}
	var err error
	templates, err = template.New("").Funcs(funcs).ParseFS(fs, "templates/*.tpl")
	if err != nil {
		panic(err)
	}
}

func toPackage(s string) string {
	return strings.ReplaceAll(s, "-", "_")
}

type RenderData struct {
	Module          string
	ProtoPath       string
	ProtoPackage    string
	Services        []string
	GatewayServices []string
	overwrite       bool
	verbose         bool
}

func (r *RenderData) ValidateTemplates() error {
	buf := new(bytes.Buffer)
	for _, tpl := range templates.Templates() {
		buf.Reset()
		if err := tpl.Execute(buf, r); err != nil {
			return fmt.Errorf("%s: %w", tpl.Name(), err)
		}
		if strings.HasSuffix(tpl.Name(), ".go.tpl") {
			if _, err := format.Source(buf.Bytes()); err != nil {
				return fmt.Errorf("%s: %w", tpl.Name(), err)
			}
		}
	}
	return nil
}

func (r *RenderData) WithOverwrite(overwrite bool) *RenderData {
	if r.overwrite == overwrite {
		return r
	}
	data := new(RenderData)
	*data = *r
	data.overwrite = overwrite
	return data
}

func (r *RenderData) Render(name string) (data []byte, err error) {
	buf := new(bytes.Buffer)
	if err = templates.ExecuteTemplate(buf, name, r); err != nil {
		return nil, err
	}
	data = buf.Bytes()
	if strings.HasSuffix(name, ".go.tpl") {
		data, err = format.Source(data)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (r *RenderData) RenderAllFile(output string) error {
	buf := new(bytes.Buffer)
	for _, tpl := range templates.Templates() {
		buf.Reset()
		if err := r.renderFile(output, buf, tpl); err != nil {
			return err
		}
	}
	return nil
}

func (r *RenderData) renderFile(output string, buf *bytes.Buffer, tpl *template.Template) error {
	if err := tpl.Execute(buf, r); err != nil {
		return err
	}
	data := buf.Bytes()
	filename, err := r.filenameFromHeader(output, tpl.Name(), data)
	if err != nil {
		return err
	}
	if strings.HasSuffix(tpl.Name(), ".go.tpl") {
		data, err = format.Source(data)
		if err != nil {
			return err
		}
	}
	_, err = autoOverwriteFile(filename, data, r.overwrite)
	if err != nil {
		return err
	}
	return nil
}

func (r *RenderData) filenameFromHeader(output string, tpl string, data []byte) (string, error) {
	reader := bufio.NewReader(bytes.NewReader(data))
	var filename = strings.TrimSuffix(tpl, ".tpl")
	for {
		line, err := reader.ReadString('\n')
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", err
		}
		if len(line) == 0 {
			break
		}
		if strings.HasPrefix(line, "// @file:") {
			filename = strings.TrimSpace(strings.TrimPrefix(line, "// @file:"))
			break
		}
	}
	return path.Join(output, filename), nil
}
