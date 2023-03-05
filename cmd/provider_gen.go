package main

import (
	cueast "cuelang.org/go/cue/ast"
	cuetoken "cuelang.org/go/cue/token"
	"fmt"
	"go.uber.org/multierr"
	goast "go/ast"
	"golang.org/x/tools/go/packages"
	"kubecue"
	"os"
	"strings"
)

const (
	TypeProviderFnMap    = "map[string]github.com/kubevela/pkg/cue/cuex/runtime.ProviderFn"
	TypeProvidersParams  = "github.com/kubevela/pkg/cue/cuex/providers.Params"
	TypeProvidersReturns = "github.com/kubevela/pkg/cue/cuex/providers.Returns"
)

type provider struct {
	name    string
	params  string
	returns string
	do      string
}

func providerGen(file string) (rerr error) {
	g, err := kubecue.NewGenerator(file)
	if err != nil {
		return err
	}

	g.RegisterAny(
		"*k8s.io/apimachinery/pkg/apis/meta/v1/unstructured.Unstructured",
		"*k8s.io/apimachinery/pkg/apis/meta/v1/unstructured.UnstructuredList",
	)

	providers, err := extractProviders(g.Package())
	if err != nil {
		return err
	}

	g.RegisterTypeFilter(func(spec *goast.TypeSpec) bool {
		typ := g.Package().TypesInfo.TypeOf(spec.Type)

		if strings.HasPrefix(typ.String(), TypeProvidersParams) ||
			strings.HasPrefix(typ.String(), TypeProvidersReturns) {
			return true
		}

		return false
	})

	decls, err := g.Generate()
	if err != nil {
		return err
	}

	newDecls, err := modifyDecls(g.Package().Name, decls, providers)
	if err != nil {
		return err
	}

	f, err := os.Create(changeExt(file, ".cue"))
	defer multierr.AppendInvoke(&rerr, multierr.Close(f))

	return g.Write(f, newDecls)
}

// extractProviders extracts the providers from the given package
func extractProviders(pkg *packages.Package) (providers []provider, rerr error) {
	// capture panic caused by invalid type assertion or out of range index,
	// so we don't need to check each type assertion and index
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		rerr = fmt.Errorf("extract providers: panic: %v", r)
	// 	}
	// }()

	var funcExpr *goast.CompositeLit
	for k, v := range pkg.TypesInfo.Types {
		if v.Type.String() == TypeProviderFnMap {
			funcExpr = k.(*goast.CompositeLit)
			break
		}
	}

	if funcExpr == nil {
		return nil, fmt.Errorf("no provider function map found like '%s'", TypeProviderFnMap)
	}

	for _, e := range funcExpr.Elts {
		kv := e.(*goast.KeyValueExpr)
		key := kv.Key.(*goast.BasicLit)
		value := kv.Value.(*goast.CallExpr)

		indices := value.Fun.(*goast.IndexListExpr)
		params := indices.Indices[0].(*goast.Ident)
		returns := indices.Indices[1].(*goast.Ident)

		do := value.Args[0].(*goast.Ident)

		providers = append(providers, provider{
			name:    key.Value,
			params:  params.Name,
			returns: returns.Name,
			do:      do.Name,
		})
	}

	return providers, nil
}

func modifyDecls(provider string, old []cueast.Decl, providers []provider) (decls []cueast.Decl, rerr error) {
	defer func() {
		if r := recover(); r != nil {
			rerr = fmt.Errorf("modify decls: panic: %v", r)
		}
	}()

	// map[StructName]StructLit
	mapping := make(map[string]cueast.Expr)
	for _, decl := range old {
		field := decl.(*cueast.Field)
		key := field.Label.(*cueast.BasicLit)

		mapping[unQuote(key.Value)] = field.Value
	}

	providerField := &cueast.Field{
		Label: kubecue.Ident("provider", true),
		Value: cueast.NewString(unQuote(provider)),
	}

	for _, p := range providers {
		params := mapping[p.params].(*cueast.StructLit).Elts
		returns := mapping[p.returns].(*cueast.StructLit).Elts

		doField := &cueast.Field{
			Label: kubecue.Ident("do", true),
			Value: cueast.NewString(unQuote(p.name)),
		}

		pdecls := []cueast.Decl{doField, providerField}
		pdecls = append(pdecls, params...)
		pdecls = append(pdecls, returns...)

		newProvider := &cueast.Field{
			Label: kubecue.Ident(p.do, true),
			Value: &cueast.StructLit{
				Elts: pdecls,
			},
		}
		cueast.SetRelPos(newProvider, cuetoken.NewSection)

		decls = append(decls, newProvider)
	}

	return decls, nil
}

func unQuote(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}
	return s
}
