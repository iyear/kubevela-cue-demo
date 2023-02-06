package kubecue

import (
	cueast "cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/ast/astutil"
	cueformat "cuelang.org/go/cue/format"
	goast "go/ast"
	"golang.org/x/tools/go/packages"
	"io"
)

type Generator struct {
	commentMap goast.CommentMap // TODO(iyear): no use in this demo

	// immutable
	pkg   *packages.Package
	types typeInfo // TODO(iyear): no use in this demo
}

func NewGenerator(f string) (*Generator, error) {
	pkg, err := loadPackage(f)
	if err != nil {
		return nil, err
	}

	types := getTypeInfo(pkg)

	return &Generator{
		pkg:   pkg,
		types: types,
	}, nil
}

func (g *Generator) Generate(w io.Writer) error {
	var decls []cueast.Decl

	for _, syntax := range g.pkg.Syntax {
		// g.commentMap = goast.NewCommentMap(g.pkg.Fset, syntax, syntax.Comments)

		for _, decl := range syntax.Decls {
			if d, ok := decl.(*goast.GenDecl); ok {
				decls = append(decls, g.convertDecls(d)...)
			}
		}
	}

	pkg := &cueast.Package{Name: ident(g.pkg.Name, false)}

	f := &cueast.File{Decls: []cueast.Decl{pkg}}
	f.Decls = append(f.Decls, decls...)

	return g.write(w, f)
}

func (g *Generator) write(w io.Writer, f *cueast.File) error {
	if err := astutil.Sanitize(f); err != nil {
		return err
	}

	b, err := cueformat.Node(f, cueformat.Simplify())
	if err != nil {
		return err
	}

	_, err = w.Write(b)
	return err
}
