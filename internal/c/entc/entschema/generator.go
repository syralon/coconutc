package entschema

import (
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
)

func Generate(target string, args []string, graph *gen.Graph) error {
	schemas, err := Parses(args)
	if err != nil {
		return err
	}
	nodes := make(map[string]*gen.Type)
	if graph != nil {
		for _, node := range graph.Nodes {
			nodes[node.Name] = node
		}
	}
	for _, schema := range schemas {
		if node, ok := nodes[schema.Name]; ok && !schema.Overwrite {
			for _, ff := range node.Fields {
				var t FieldType
				switch ff.Type.Type {
				case field.TypeTime:
					t = Time
				case field.TypeJSON:
					t = "JSON"
				default:
					t = FieldType(ff.Type.String())
				}
				schema.Fields.add(&Field{
					Name: ff.Name,
					Type: t,
				})
			}
		}
		if err = schema.WriteFile(target); err != nil {
			return err
		}
	}
	return nil
}
