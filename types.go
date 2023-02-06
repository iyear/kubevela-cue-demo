package kubecue

import (
	goast "go/ast"
	gotypes "go/types"
	"golang.org/x/tools/go/packages"
)

type typeInfo map[gotypes.Type]*goast.StructType

func getTypeInfo(p *packages.Package) typeInfo {
	m := make(typeInfo)

	for _, f := range p.Syntax {
		goast.Inspect(f, func(n goast.Node) bool {
			switch x := n.(type) {
			case *goast.StructType:
				m[p.TypesInfo.TypeOf(x)] = x
			}
			return true
		})
	}

	return m
}

func supportedType(stack []gotypes.Type, t gotypes.Type) (ok bool) {
	// handle recursive types
	for _, t0 := range stack {
		if t0 == t {
			return true
		}
	}
	stack = append(stack, t)

	t = t.Underlying()
	switch x := t.(type) {
	case *gotypes.Basic:
		return x.String() != "invalid type"
	case *gotypes.Named:
		return true
	case *gotypes.Pointer:
		return supportedType(stack, x.Elem())
	case *gotypes.Slice:
		return supportedType(stack, x.Elem())
	case *gotypes.Array:
		return supportedType(stack, x.Elem())
	case *gotypes.Map:
		if b, ok := x.Key().Underlying().(*gotypes.Basic); !ok || b.Kind() != gotypes.String {
			return false
		}
		return supportedType(stack, x.Elem())
	case *gotypes.Struct:
		// Eliminate structs with fields for which all fields are filtered.
		if x.NumFields() == 0 {
			return true
		}
		for i := 0; i < x.NumFields(); i++ {
			f := x.Field(i)
			if f.Exported() && supportedType(stack, f.Type()) {
				return true
			}
		}
	case *gotypes.Interface:
		return true
	}
	return false
}
