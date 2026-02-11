package entservice

import (
	"context"
	"fmt"
	"path"
	"strings"

	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/syralon/coconutc/internal/tools/text"
	"github.com/syralon/coconutc/pkg/annotation/entproto"
)

type repositoryBuilder struct {
	*BuildOptions

	file          *jen.File
	name          string
	node          *gen.Type
	entityPackage string
	txPackage     string
}

func (b *repositoryBuilder) build(_ context.Context) error {
	b.structs()
	for _, method := range b.APIOptions.Method.Methods() {
		switch method {
		case entproto.GET:
			b.get()
		case entproto.LIST:
			b.list()
		case entproto.CREATE:
			b.create()
		case entproto.UPDATE:
			b.update()
		case entproto.DELETE:
			b.delete()
		default:
		}
	}
	b.set()
	b.edges()
	b.edge()

	return nil
}

func (b *repositoryBuilder) structs() {
	defer b.file.Line()

	name := strcase.ToLowerCamel(b.node.Name)
	repName := name + "Repository"
	b.file.Type().Id(repName).Interface(
		jen.Id(b.node.Name).Params(ctxVar()).Params(jen.Op("*").Qual(b.EntPackage, b.node.Name+"Client")),
	).Line()
	b.file.Type().Id(b.name).Struct(
		jen.Id(repName),
	)
	b.file.Func().Id("New" + b.name).Params(
		jen.Id("repo").Op("*").Qual(b.txPackage, "Repository"),
	).Op("*").Id(b.name).Block(
		jen.Return(jen.Op("&").Id(b.name).Block(
			jen.Id(repName).Op(":").Id("repo").Op(","),
		)),
	)
}

func (b *repositoryBuilder) fn(fn *jen.Statement) *jen.Statement {
	return b.file.Func().Op("(").Id("rep").Op("*").Id(b.name).Op(")").Add(fn)
}
func (b *repositoryBuilder) fnGet() *jen.Statement {
	return jen.Id("Get").
		Params(
			jen.Id("ctx").Qual(pkgContext, "Context"),
			jen.Id("id").Id(b.node.IDType.String()),
		).
		Call(jen.Op("*").Qual(b.entityPackage, b.node.Name), jen.Error())
}

func (b *repositoryBuilder) fnList(paginatorName string) *jen.Statement {
	return jen.Id("List").
		Params(
			jen.Id("ctx").Qual(pkgContext, "Context"),
			jen.Id("options").Op("*").Qual(b.ProtoPackage, b.node.Name+"Options"),
			jen.Id("paginator").Op("*").Qual(pkgCoconutField, paginatorName),
		).
		Call(
			jen.Index().Op("*").Qual(b.entityPackage, b.node.Name),
			jen.Op("*").Qual(pkgCoconutField, paginatorName),
			jen.Error(),
		)
}

func (b *repositoryBuilder) fnCreate() *jen.Statement {
	return jen.Id("Create").
		Params(
			jen.Id("ctx").Qual(pkgContext, "Context"),
			jen.Id("create").Op("*").Qual(b.ProtoPackage, b.node.Name+"Create"),
		).
		Call(jen.Op("*").Qual(b.entityPackage, b.node.Name), jen.Error())
}

func (b *repositoryBuilder) fnUpdate() *jen.Statement {
	return jen.Id("Update").
		Params(
			jen.Id("ctx").Qual(pkgContext, "Context"),
			jen.Id("id").Id(b.node.IDType.String()),
			jen.Id("data").Op("*").Qual(b.ProtoPackage, b.node.Name+"Update"),
		).
		Call(jen.Error())
}

func (b *repositoryBuilder) fnDelete() *jen.Statement {
	return jen.Id("Delete").
		Params(
			jen.Id("ctx").Qual(pkgContext, "Context"),
			jen.Id("id").Id(b.node.IDType.String()),
		).
		Call(jen.Error())
}

func (b *repositoryBuilder) fnSet(v *gen.Field, opts entproto.FieldOptions) *jen.Statement {
	typ := jen.Id(v.Type.String())
	if opts.ProtoEnum {
		typ = jen.Qual(b.ProtoPackage, strcase.ToScreamingSnake(b.node.Name+"_"+v.Name))
	}
	return jen.Id(fmt.Sprintf("Set%s", text.ProtoPascal(v.Name))).
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("id").Id(b.node.IDType.Type.String()),
			jen.Id("value").Add(typ),
		).Error()

}
func (b *repositoryBuilder) fnEdges(edge *gen.Edge, paginatorName string) *jen.Statement {
	name := fmt.Sprintf("List%s", text.ProtoPascal(edge.Name))
	return jen.Id(name).
		Params(
			jen.Id("ctx").Qual(pkgContext, "Context"),
			jen.Id("id").Id(b.node.IDType.String()),
			jen.Id("options").Op("*").Qual(b.ProtoPackage, fmt.Sprintf("%sOptions", edge.Type.Name)),
			jen.Id("paginator").Op("*").Qual(pkgCoconutField, paginatorName),
		).
		Call(
			jen.Index().Op("*").Qual(b.entityPackage, edge.Type.Name),
			jen.Op("*").Qual(pkgCoconutField, paginatorName),
			jen.Error(),
		)
}

func (b *repositoryBuilder) fnEdge(edge *gen.Edge) *jen.Statement {
	name := fmt.Sprintf("Get%s", text.ProtoPascal(edge.Name))
	return jen.Id(name).
		Params(
			jen.Id("ctx").Qual(pkgContext, "Context"),
			jen.Id("id").Id(b.node.IDType.String()),
		).
		Call(
			jen.Op("*").Qual(b.entityPackage, edge.Type.Name),
			jen.Error(),
		)
}

func (b *repositoryBuilder) set() {
	for _, v := range b.node.Fields {
		if v.Name == "created_at" || v.Name == "updated_at" {
			continue
		}
		if v.Immutable {
			continue
		}
		fieldOpts, err := entproto.GetFieldOptions(v.Annotations)
		if err != nil {
			panic(err)
		}
		if fieldOpts.Immutable || !fieldOpts.Settable {
			continue
		}
		fnName := fmt.Sprintf("Set%s", text.ProtoPascal(v.Name))

		b.fn(b.fnSet(v, fieldOpts)).
			Block(
				define("_", "err").Id("rep").Dot(b.node.Name).Call(jen.Id("ctx")).Dot("Update").Call().
					Dot(fnName).Call(entType(v.Type.Type, jen.Id("value"), &fieldOpts)).Dot("Save").Call(jen.Id("ctx")),
				jen.Return(jen.Err()),
			).Line()
	}
}

func (b *repositoryBuilder) edges() {
	for _, edge := range b.node.Edges {
		opts, err := entproto.GetAPIOptions(edge.Annotations)
		if err != nil {
			return
		}
		if opts.DisableEdge || edge.Unique {
			continue
		}
		var fields []*jen.Statement
		for _, v := range edge.Type.Fields {
			fieldOpts, err := entproto.GetFieldOptions(v.Annotations)
			if err != nil {
				return
			}
			if fieldOpts.Sensitive || !fieldOpts.Filterable {
				continue
			}
			fields = append(fields, jen.Id("options").Dot("Get"+text.ProtoPascal(v.Name)).Call().Dot("Selector").Call(
				jen.Qual(path.Join(b.EntPackage, strings.ToLower(edge.Type.Name)), fmt.Sprintf("Field%s", text.EntPascal(v.Name))),
			))
		}
		b.fn(b.fnEdges(edge, opts.PaginatorStyle.String())).Block(
			define("query").Add(b.queryEdge()).Dot(fmt.Sprintf("Query%s", text.EntPascal(edge.Name))).Call().Dot("Where").Call(
				jen.Qual(pkgCoconutField, "Selectors").Index(jen.Qual(path.Join(b.EntPackage, "predicate"), edge.Type.Name)).Add(calls(fields...)).Op("..."),
			),
			jen.Line(),
			b.paginator(opts.PaginatorStyle),
			jen.Line(),
			define("data", "err").Id("query").Dot("All").Call(jen.Id("ctx")),
			jen.If(jen.Id("err").Op("!=").Nil()).Block(
				jen.Return(jen.Nil(), jen.Id("paginator"), jen.Id("err")),
			),
			jen.Line(),
			jen.Return(
				jen.Qual(pkgXSlices, "Trans").Call(jen.Id("data"), jen.Qual(b.entityPackage, "New"+edge.Type.Name)),
				jen.Id("paginator"),
				jen.Nil(),
			),
		)
	}
}

func (b *repositoryBuilder) edge() {
	for _, edge := range b.node.Edges {
		opts, err := entproto.GetAPIOptions(edge.Annotations)
		if err != nil {
			return
		}
		if opts.DisableEdge || !edge.Unique {
			continue
		}
		b.fn(b.fnEdge(edge)).
			Block(
				define("query").Add(b.queryEdge()).Dot(fmt.Sprintf("Query%s", text.EntPascal(edge.Name))).Call(),
				define("data", "err").Id("query").Dot("First").Call(jen.Id("ctx")),
				ifErr(),
				jen.Return(
					jen.Qual(b.entityPackage, "New"+edge.Type.Name).Call(jen.Id("data")),
					jen.Err(),
				),
			)
	}
}

func (b *repositoryBuilder) get() {
	defer b.file.Line()
	b.fn(b.fnGet()).
		Block(
			define("data", "err").Id("rep").Dot(b.node.Name).Call(jen.Id("ctx")).Dot("Query()").Dot("Where").
				Call(jen.Qual(path.Join(b.EntPackage, strings.ToLower(b.node.Name)), "ID").Call(jen.Id("id"))).
				Dot("First").Call(jen.Id("ctx")),
			jen.Return(
				jen.Qual(b.entityPackage, "New"+b.node.Name).Call(jen.Id("data")),
				jen.Err(),
			),
		)
}

func (b *repositoryBuilder) list() {
	defer b.file.Line()
	paginatorName := b.APIOptions.PaginatorStyle.String()

	var fields []*jen.Statement
	for _, v := range b.node.Fields {
		fieldOpts, err := entproto.GetFieldOptions(v.Annotations)
		if err != nil {
			return
		}
		if fieldOpts.Sensitive || !fieldOpts.Filterable {
			continue
		}
		fields = append(fields, jen.Id("options").Dot("Get"+text.ProtoPascal(v.Name)).Call().Dot("Selector").Call(
			jen.Qual(path.Join(b.EntPackage, strings.ToLower(b.node.Name)), fmt.Sprintf("Field%s", text.EntPascal(v.Name))),
		))
	}

	b.fn(b.fnList(paginatorName)).
		Block(
			define("query").Id("rep").Dot(b.node.Name).Call(jen.Id("ctx")).Dot("Query()").Dot("Where").Call(
				jen.Qual(pkgCoconutField, "Selectors").Index(jen.Qual(path.Join(b.EntPackage, "predicate"), b.node.Name)).Add(calls(fields...)).Op("..."),
			),
			jen.Line(),
			b.paginator(b.APIOptions.PaginatorStyle),
			jen.Line(),
			define("data", "err").Id("query").Dot("All").Call(jen.Id("ctx")),
			jen.If(jen.Id("err").Op("!=").Nil()).Block(
				jen.Return(jen.Nil(), jen.Id("paginator"), jen.Id("err")),
			),
			jen.Line(),
			jen.Return(
				jen.Qual(pkgXSlices, "Trans").Call(jen.Id("data"), jen.Qual(b.entityPackage, "New"+b.node.Name)),
				jen.Id("paginator"),
				jen.Nil(),
			),
		)
}

func (b *repositoryBuilder) create() {
	defer b.file.Line()
	create := define("op").Id("rep").Dot(b.node.Name).Call(jen.Id("ctx")).Dot("Create()")
	for _, v := range b.node.Fields {
		fieldOpts, err := entproto.GetFieldOptions(v.Annotations)
		if err != nil {
			return
		}

		if v.Name == "created_at" || v.Name == "updated_at" {
			continue
		}
		val := jen.Id("create").Dot(fmt.Sprintf("Get%s()", text.ProtoPascal(v.Name)))
		if v.Type.Type == field.TypeTime {
			val = val.Dot("AsTime()")
		}
		val = wrap(v, val, &fieldOpts)
		create = create.Op(".").Id("\n").Id(fmt.Sprintf("Set%s", text.EntPascal(v.Name))).Call(val)
	}

	b.fn(b.fnCreate()).
		Block(
			create,
			define("data", "err").Id("op").Dot("Save").Call(jen.Id("ctx")),
			ifErr(),
			jen.Return(
				jen.Qual(b.entityPackage, "New"+b.node.Name).Call(jen.Id("data")),
				jen.Err(),
			),
		)
}

func (b *repositoryBuilder) update() {
	defer b.file.Line()
	var fields []jen.Code
	for _, v := range b.node.Fields {
		if v.Name == "created_at" || v.Name == "updated_at" {
			continue
		}
		if v.Immutable {
			continue
		}
		fieldOpts, err := entproto.GetFieldOptions(v.Annotations)
		if err != nil {
			panic(err)
		}
		if fieldOpts.Immutable {
			continue
		}
		val := jen.Id("data").Dot(fmt.Sprintf("Get%s()", text.ProtoPascal(v.Name)))
		if v.Type.Type == field.TypeTime {
			val = val.Dot("AsTime()")
		}
		val = wrap(v, val, &fieldOpts)
		fields = append(
			fields,
			jen.If(jen.Id("data").Dot(text.ProtoPascal(v.Name)).Op("!=").Nil()).Block(
				jen.Id("update").Dot(fmt.Sprintf("Set%s", text.EntPascal(v.Name))).Call(val),
			).Line(),
		)
	}

	b.fn(b.fnUpdate()).
		Block(
			define("update").Id("rep").Dot(b.node.Name).Call(jen.Id("ctx")).Dot("UpdateOneID").Call(jen.Id("id")),
			jen.Add(fields...),
			define("_", "err").Id("update").Dot("Save").Call(jen.Id("ctx")),
			jen.Return(jen.Err()),
		)
}

func (b *repositoryBuilder) delete() {
	defer b.file.Line()
	b.fn(b.fnDelete()).
		Block(
			jen.Return(jen.Id("rep").Dot(b.node.Name).Call(jen.Id("ctx")).Dot("DeleteOneID").Call(jen.Id("id"))).Dot("Exec").Call(jen.Id("ctx")),
		)
}

func (b *repositoryBuilder) queryEdge() *jen.Statement {
	return jen.Id("rep").Dot(b.node.Name).Call(jen.Id("ctx")).Dot("Query").Call().Dot("Where").Call(
		jen.Qual(path.Join(b.EntPackage, strings.ToLower(b.node.Name)), "ID").Call(jen.Id("id")),
	)
}

func (b *repositoryBuilder) paginator(style entproto.PaginatorStyle) *jen.Statement {
	if style == entproto.InfinitePaginator {
		return b.infinitePaginator()
	}
	return b.classicalPaginator()
}

func (b *repositoryBuilder) classicalPaginator() *jen.Statement {
	return jen.If(jen.Id("paginator").Op("!=").Nil()).Block(
		define("total", "err").Id("query").Dot("Count").Call(jen.Id("ctx")),
		jen.If(jen.Id("err").Op("!=").Nil()).Block(
			jen.Return(jen.Nil(), jen.Id("paginator"), jen.Id("err")),
		),
		jen.Id("paginator").Dot("Total").Op("=").Int64().Call(jen.Id("total")),
		assign("query").Id("query").Dot("Order").Call(jen.Id("paginator").Dot("OrderSelector").Call()).
			Dot("Offset").Call(jen.Int().Call(jen.Id("paginator").Dot("GetLimit()").Op("*").Call(jen.Id("paginator").Dot("GetPage()").Op("-").Id("1")))),
	)
}

func (b *repositoryBuilder) infinitePaginator() *jen.Statement {
	return jen.If(jen.Id("paginator").Op("!=").Nil()).Block(
		assign("query").Id("query").Dot("Order").Call(jen.Qual(path.Join(b.EntPackage, strings.ToLower(b.node.Name)), "ByID").Call(
			jen.Qual(pkgEntSql, "OrderDesc()"),
		)).Dot("Limit").Call(jen.Int().Call(jen.Id("paginator").Dot("GetLimit()"))),
		jen.If(
			define("sequence").Id("paginator").Dot("GetSequence()").Op(";").Id("sequence").Op(">").Id("0"),
		).Block(
			assign("query").Id("query").Dot("Where").Call(
				jen.Qual(path.Join(b.EntPackage, strings.ToLower(b.node.Name)), "IDLT").Call(jen.Id("sequence")),
			),
		),
	)
}

func RepositoryBuilder(pkgName, entityPkg, txPackage string) BuildFunc {
	return func(ctx context.Context, node *gen.Type, opts *BuildOptions) (*jen.File, error) {
		file := jen.NewFile(pkgName)
		file.HeaderComment("Code generated by coconutc. DO NOT EDIT.")
		file.HeaderComment("https://github.com/syralon/coconutc")
		file.HeaderComment("ent." + node.Name)
		file.ImportAlias(opts.ProtoPackage, "pb")
		b := &repositoryBuilder{
			BuildOptions:  opts,
			name:          fmt.Sprintf("%sRepository", node.Name),
			node:          node,
			file:          file,
			entityPackage: entityPkg,
			txPackage:     txPackage,
		}
		err := b.build(ctx)
		if err != nil {
			return nil, err
		}
		return file, nil
	}
}

func RepositoryInterfaceBuilder(pkgName, entityPkg string) BuildFunc {
	return func(ctx context.Context, node *gen.Type, opts *BuildOptions) (*jen.File, error) {
		file := jen.NewFile(pkgName)
		file.HeaderComment("Code generated by coconutc. DO NOT EDIT.")
		file.HeaderComment("https://github.com/syralon/coconutc")
		file.ImportAlias(opts.ProtoPackage, "pb")
		b := &repositoryBuilder{
			BuildOptions:  opts,
			file:          file,
			name:          fmt.Sprintf("%sRepository", node.Name),
			node:          node,
			entityPackage: entityPkg,
		}
		var funcs []jen.Code
		for _, method := range b.APIOptions.Method.Methods() {
			switch method {
			case entproto.GET:
				funcs = append(funcs, b.fnGet())
			case entproto.LIST:
				funcs = append(funcs, b.fnList(opts.APIOptions.PaginatorStyle.String()))
			case entproto.CREATE:
				funcs = append(funcs, b.fnCreate())
			case entproto.UPDATE:
				funcs = append(funcs, b.fnUpdate())
			case entproto.DELETE:
				funcs = append(funcs, b.fnDelete())
			default:
			}
		}

		for _, v := range b.node.Fields {
			if v.Name == "created_at" || v.Name == "updated_at" {
				continue
			}
			if v.Immutable {
				continue
			}
			fieldOpts, err := entproto.GetFieldOptions(v.Annotations)
			if err != nil {
				panic(err)
			}
			if fieldOpts.Immutable || !fieldOpts.Settable {
				continue
			}
			funcs = append(funcs, b.fnSet(v, fieldOpts))
		}

		for _, edge := range b.node.Edges {
			edgeOpts, err := entproto.GetAPIOptions(edge.Annotations)
			if err != nil {
				return nil, err
			}
			if edgeOpts.DisableEdge {
				continue
			}
			if edge.Unique {
				funcs = append(funcs, b.fnEdge(edge))
			} else {
				funcs = append(funcs, b.fnEdges(edge, edgeOpts.PaginatorStyle.String()))
			}
		}
		file.Type().Id(b.name).Interface(funcs...)
		return file, nil
	}
}
