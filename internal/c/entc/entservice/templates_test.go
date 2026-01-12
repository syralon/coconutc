package entservice

import (
	"fmt"
	"testing"
)

func TestTemplate(t *testing.T) {
	data := &RenderData{
		Module:       "github.com/syralon/example",
		Services:     []string{"Hello", "Boom"},
		ProtoPath:    "proto/syralon/example",
		ProtoPackage: "github.com/syralon/example/proto/syralon/example",
	}
	if err := data.ValidateTemplates(); err != nil {
		t.Error(err)
	}
}

func TestTemplate_service(t *testing.T) {
	data := RenderData{
		Module:   "github.com/syralon/example",
		Services: []string{"Hello", "Boom"},
	}
	b, err := data.Render("service.go.tpl")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(b))
}

func TestGitignore(t *testing.T) {
	data := &RenderData{
		Module:       "github.com/syralon/example",
		Services:     []string{"Hello", "Boom"},
		ProtoPath:    "proto/syralon/example",
		ProtoPackage: "github.com/syralon/example/proto/syralon/example",
	}
	b, err := data.Render("gitignore.tpl")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(b))
}
