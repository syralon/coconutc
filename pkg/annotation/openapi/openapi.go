package openapi

import (
	"errors"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	openapiv3 "github.com/google/gnostic/openapiv3"
	"github.com/syralon/coconutc/pkg/annotation/helper"
)

const (
	schemaAnnotationName = "openapi_annotation"
)

type schema struct {
	Schema *openapiv3.Schema `json:"schema"`
}

func (s *schema) Name() string {
	return schemaAnnotationName
}

func Schema(s *openapiv3.Schema) entc.Annotation {
	return &schema{Schema: s}
}

func GetSchema(annotations gen.Annotations) (*openapiv3.Schema, error) {
	s := &schema{}
	err := helper.GetAnnotations(annotations, schemaAnnotationName, s)
	if errors.Is(err, helper.ErrAnnotationNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return s.Schema, nil
}
