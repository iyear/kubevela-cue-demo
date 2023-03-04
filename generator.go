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
	// immutable
	pkg   *packages.Package
	types typeInfo

	anyTypes map[string]struct{}
}

var defaultAnyTypes = []string{
	"map[string]interface{}",
	"map[string]any",
	"interface{}",
	"any",
}

func NewGenerator(f string) (*Generator, error) {
	pkg, err := loadPackage(f)
	if err != nil {
		return nil, err
	}

	types := getTypeInfo(pkg)

	g := &Generator{
		pkg:      pkg,
		types:    types,
		anyTypes: make(map[string]struct{}),
	}

	g.RegisterAny(defaultAnyTypes...)

	return g, nil
}

func (g *Generator) Generate() ([]cueast.Decl, error) {
	var decls []cueast.Decl
	for _, syntax := range g.pkg.Syntax {
		for _, decl := range syntax.Decls {
			if d, ok := decl.(*goast.GenDecl); ok {
				t, err := g.convertDecls(d)
				if err != nil {
					return nil, err
				}
				decls = append(decls, t...)
			}
		}
	}

	return decls, nil
}

func (g *Generator) Write(w io.Writer, decls []cueast.Decl) error {
	pkg := &cueast.Package{Name: ident(g.pkg.Name, false)}

	f := &cueast.File{Decls: []cueast.Decl{pkg}}
	f.Decls = append(f.Decls, decls...)

	if err := astutil.Sanitize(f); err != nil {
		return err
	}

	b, err := cueformat.Node(f, cueformat.Simplify(), cueformat.TabIndent(true))
	if err != nil {
		return err
	}

	_, err = w.Write(b)
	return err
}

func (g *Generator) Package() *packages.Package {
	return g.pkg
}
