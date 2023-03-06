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
	pkg   *packages.Package
	types typeInfo
	opts  *options
}

func NewGenerator(f string) (*Generator, error) {
	pkg, err := loadPackage(f)
	if err != nil {
		return nil, err
	}

	g := &Generator{
		pkg:   pkg,
		types: getTypeInfo(pkg),
	}

	return g, nil
}

// Generate can be called multiple times with different options.
func (g *Generator) Generate(opts ...Option) ([]cueast.Decl, error) {
	g.opts = defaultOptions
	for _, opt := range opts {
		opt(g.opts)
	}

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
	pkg := &cueast.Package{Name: Ident(g.pkg.Name, false)}

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
