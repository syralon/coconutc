package entservice

import (
	"entgo.io/ent/entc/gen"
	"github.com/dave/jennifer/jen"
)

type builder struct {
	name string
	file *jen.File
	node *gen.Type
}
