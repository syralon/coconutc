package entschema

import (
	"os"
	"testing"
)

func TestSchemaTemplate(t *testing.T) {
	{
		schema, err := Parse("Example")
		if err != nil {
			t.Fatal(err)
		}
		err = schema.Execute(os.Stdout)
		if err != nil {
			t.Fatal(err)
		}
	}
	{
		schema, err := Parse("Example(id:string)")
		if err != nil {
			t.Fatal(err)
		}
		err = schema.Execute(os.Stdout)
		if err != nil {
			t.Fatal(err)
		}
	}
}
