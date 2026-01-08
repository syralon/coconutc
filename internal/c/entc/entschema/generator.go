package entschema

import "entgo.io/ent/entc/gen"

func Generate(target string, args []string, graph *gen.Graph) error {
	schemas, err := Parses(args)
	if err != nil {
		return err
	}
	nodes := make(map[string]*gen.Type)
	for _, node := range graph.Nodes {
		nodes[node.Name] = node
	}
	for _, schema := range schemas {
		if node, ok := nodes[schema.Name]; ok && !schema.Overwrite {
			for _, field := range node.Fields {
				schema.Fields.add(&Field{
					Name: field.Name,
					Type: FieldType(field.Type.String()),
				})
			}
		}
		if err = schema.WriteFile(target); err != nil {
			return err
		}
	}
	return nil
}
