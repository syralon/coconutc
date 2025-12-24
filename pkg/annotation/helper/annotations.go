package helper

import (
	"encoding/json"
	"errors"

	"entgo.io/ent/entc/gen"
)

var (
	ErrAnnotationNotFound = errors.New("annotation not found")
)

func GetAnnotations(annotations gen.Annotations, name string, dst any) error {
	ann, ok := annotations[name]
	if !ok {
		return ErrAnnotationNotFound
	}
	data, err := json.Marshal(ann)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}
