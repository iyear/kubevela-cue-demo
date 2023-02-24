package kubecue

import (
	cueast "cuelang.org/go/cue/ast"
	cuetoken "cuelang.org/go/cue/token"
	"fmt"
	goast "go/ast"
	gotoken "go/token"
	gotypes "go/types"
	"strconv"
)

func (g *Generator) convertDecls(x *goast.GenDecl) (decls []cueast.Decl, _ error) {
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

		if err := supportedType(nil, typ); err != nil {
			return nil, fmt.Errorf("unsupported type %s: %w", typeSpec.Name.Name, err)
		}

		named, ok := typ.(*gotypes.Named)
		if !ok {
			continue
		}
		st, ok := named.Underlying().(*gotypes.Struct)
		if !ok {
			continue
		}

		lit, err := g.convert(st)
		if err != nil {
			return nil, err
		}

		field := &cueast.Field{
			Label: cueast.NewString(typeSpec.Name.Name),
			Value: lit,
		}
		// there is no doc for typeSpec, so we only add x.Doc
		makeComments(field, &commentUnion{comment: nil, doc: x.Doc})

		cueast.SetRelPos(field, cuetoken.Blank)
		decls = append(decls, field)
	}

	return decls, nil
}

func (g *Generator) convert(typ gotypes.Type) (cueast.Expr, error) {
	if _, ok := g.anyTypes[typ.String()]; ok {
		return anyLit(), nil
	}

	switch t := typ.(type) {
	case *gotypes.Basic:
		return basicType(t), nil
	case *gotypes.Named:
		return g.convert(t.Underlying())
	case *gotypes.Struct:
		return g.makeStructLit(t)
	case *gotypes.Pointer:
		expr, err := g.convert(t.Elem())
		if err != nil {
			return nil, err
		}
		return &cueast.BinaryExpr{
			X:  cueast.NewNull(),
			Op: cuetoken.OR,
			Y:  expr,
		}, nil
	case *gotypes.Slice:
		if t.Elem().String() == "byte" {
			return ident("bytes", false), nil
		}
		expr, err := g.convert(t.Elem())
		if err != nil {
			return nil, err
		}
		return cueast.NewList(&cueast.Ellipsis{Type: expr}), nil
	case *gotypes.Array:
		if t.Elem().String() == "byte" {
			// TODO: no way to constraint lengths of bytes for now, as regexps
			// operate on Unicode, not bytes. So we need
			//     fmt.Fprint(e.w, fmt.Sprintf("=~ '^\C{%d}$'", x.Len())),
			// but regexp does not support that.
			// But translate to bytes, instead of [...byte] to be consistent.
			return ident("bytes", false), nil
		}

		expr, err := g.convert(t.Elem())
		if err != nil {
			return nil, err
		}
		return &cueast.BinaryExpr{
			X: &cueast.BasicLit{
				Kind:  cuetoken.INT,
				Value: strconv.Itoa(int(t.Len())),
			},
			Op: cuetoken.MUL,
			Y:  cueast.NewList(expr),
		}, nil
	case *gotypes.Map:
		if b, ok := t.Key().Underlying().(*gotypes.Basic); !ok || b.Kind() != gotypes.String {
			return nil, fmt.Errorf("unsupported map key type %s of %s", t.Key(), t)
		}

		elem := t.Elem()
		// if map is map[string]interface{}, we treat it as {...}
		if i, ok := elem.Underlying().(*gotypes.Interface); ok && i.Empty() {
			return anyLit(), nil
		}
		expr, err := g.convert(elem)
		if err != nil {
			return nil, err
		}

		f := &cueast.Field{
			Label: cueast.NewList(ident("string", false)),
			Value: expr,
		}
		return &cueast.StructLit{
			Elts: []cueast.Decl{f},
		}, nil
	case *gotypes.Interface:
		return ident("_", false), nil
	}

	return nil, fmt.Errorf("unsupported type %s", typ)
}

func (g *Generator) makeStructLit(x *gotypes.Struct) (*cueast.StructLit, error) {
	st := &cueast.StructLit{
		Elts: make([]cueast.Decl, 0),
	}

	// if num of fields is 1, we don't need braces. Keep it simple.
	if x.NumFields() > 1 {
		st.Lbrace = cuetoken.Blank.Pos()
		st.Rbrace = cuetoken.Newline.Pos()
	}

	err := g.addFields(st, x, map[string]struct{}{})
	if err != nil {
		return nil, err
	}

	return st, nil
}

func (g *Generator) addFields(st *cueast.StructLit, x *gotypes.Struct, names map[string]struct{}) error {
	comments := g.fieldComments(x)

	for i := 0; i < x.NumFields(); i++ {
		field := x.Field(i)

		// skip unexported fields
		if !field.Exported() {
			continue
		}

		// TODO(iyear): support more complex tags and usages
		opts := g.parseTag(x.Tag(i))

		// skip fields with "-" tag
		if opts.Name == "-" {
			continue
		}

		if opts.Name == "" {
			opts.Name = field.Name()
		}

		// process anonymous field with inline tag
		if field.Anonymous() && opts.Inline {
			if t, ok := field.Type().Underlying().(*gotypes.Struct); ok {
				if err := g.addFields(st, t, names); err != nil {
					return err
				}
			}
			continue
		}

		expr, err := g.convert(field.Type())
		if err != nil {
			return err
		}

		// can't decl same field in the same scope
		if _, ok := names[opts.Name]; ok {
			return fmt.Errorf("field '%s' already exists, can not declare duplicate field name", opts.Name)
		}
		names[opts.Name] = struct{}{}

		f := &cueast.Field{
			Label: cueast.NewString(opts.Name),
			Value: expr,
		}

		// process field with optional tag
		if opts.Optional {
			f.Token = cuetoken.COLON
			f.Optional = cuetoken.Blank.Pos()
		}

		makeComments(f, comments[i])

		st.Elts = append(st.Elts, f)
	}

	return nil
}
