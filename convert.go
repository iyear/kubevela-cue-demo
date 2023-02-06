package kubecue

import (
	cueast "cuelang.org/go/cue/ast"
	cuetoken "cuelang.org/go/cue/token"
	"fmt"
	goast "go/ast"
	gotoken "go/token"
	gotypes "go/types"
	"reflect"
	"strconv"
)

func (g *Generator) convertDecls(x *goast.GenDecl) (decls []cueast.Decl) {
	if x.Tok != gotoken.TYPE { // TODO(iyear): currently only support 'type'
		return
	}

	for _, spec := range x.Specs {
		typeSpec, ok := spec.(*goast.TypeSpec)
		if !ok {
			continue
		}

		// only process struct
		typ := g.pkg.TypesInfo.TypeOf(typeSpec.Name)

		if !supportedType(nil, typ) {
			// TODO(iyear): log? panic? ignore?
			fmt.Println("unsupported type:", typ.String())
			continue
		}

		named, ok := typ.(*gotypes.Named)
		if !ok {
			continue
		}
		structType, ok := named.Underlying().(*gotypes.Struct)
		if !ok {
			continue
		}

		field := &cueast.Field{
			Label: cueast.NewString(typeSpec.Name.Name),
			Value: g.makeStructLit(structType),
		}
		cueast.SetRelPos(field, cuetoken.Blank)
		decls = append(decls, field)
	}

	return decls
}

func (g *Generator) convert(typ gotypes.Type) cueast.Expr {
	switch t := typ.(type) {
	case *gotypes.Basic:
		return basicType(t)
	case *gotypes.Named:
		return g.convert(t.Underlying())
	case *gotypes.Struct:
		return g.makeStructLit(t)
	case *gotypes.Pointer:
		return &cueast.BinaryExpr{
			X:  cueast.NewNull(),
			Op: cuetoken.OR,
			Y:  g.convert(t.Elem()),
		}
	case *gotypes.Slice:
		if t.Elem().String() == "byte" {
			return ident("bytes", false)
		}
		return cueast.NewList(&cueast.Ellipsis{Type: g.convert(t.Elem())})
	case *gotypes.Array:
		if t.Elem().String() == "byte" {
			// TODO: no way to constraint lengths of bytes for now, as regexps
			// operate on Unicode, not bytes. So we need
			//     fmt.Fprint(e.w, fmt.Sprintf("=~ '^\C{%d}$'", x.Len())),
			// but regexp does not support that.
			// But translate to bytes, instead of [...byte] to be consistent.
			return ident("bytes", false)
		}
		return &cueast.BinaryExpr{
			X: &cueast.BasicLit{
				Kind:  cuetoken.INT,
				Value: strconv.Itoa(int(t.Len())),
			},
			Op: cuetoken.MUL,
			Y:  cueast.NewList(g.convert(t.Elem())),
		}
	case *gotypes.Map:
		if b, ok := t.Key().Underlying().(*gotypes.Basic); !ok || b.Kind() != gotypes.String {
			panic(fmt.Sprintf("unsupported map key type %T", t.Key()))
		}

		f := &cueast.Field{
			Label: cueast.NewList(ident("string", false)),
			Value: g.convert(t.Elem()),
		}
		return &cueast.StructLit{
			Lbrace: cuetoken.Blank.Pos(),
			Elts:   []cueast.Decl{f},
			Rbrace: cuetoken.Blank.Pos(),
		}
	case *gotypes.Interface:
		return ident("_", false)
	}

	// TODO(iyear): placeholder? panic? error?
	return ident("TODO", false)
}

const tagName = "cue"

func (g *Generator) makeStructLit(x *gotypes.Struct) *cueast.StructLit {
	st := &cueast.StructLit{
		Elts:   make([]cueast.Decl, 0),
		Lbrace: cuetoken.Blank.Pos(),
		Rbrace: cuetoken.Newline.Pos(),
	}

	for i := 0; i < x.NumFields(); i++ {
		field := x.Field(i)

		expr := g.convert(field.Type())

		fieldName := field.Name()
		// TODO(iyear): support more complex tags and usages
		if tag := x.Tag(i); tag != "" {
			if n, ok := reflect.StructTag(tag).Lookup(tagName); ok {
				fieldName = n
			}
		}

		f := &cueast.Field{
			Label: cueast.NewString(fieldName),
			Value: expr,
		}

		st.Elts = append(st.Elts, f)
	}

	return st
}
