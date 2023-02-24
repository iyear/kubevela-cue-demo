package kubecue

import (
	cueast "cuelang.org/go/cue/ast"
	gotypes "go/types"
	"unicode"
)

func ident(name string, isDef bool) *cueast.Ident {
	if isDef {
		r := []rune(name)[0]
		name = "#" + name
		if !unicode.Is(unicode.Lu, r) {
			name = "_" + name
		}
	}
	return cueast.NewIdent(name)
}

func basicType(x *gotypes.Basic) cueast.Expr {
	switch t := x.String(); t {
	case "uintptr":
		return ident("uint64", false)
	case "byte":
		return ident("uint8", false)
	default:
		return ident(t, false)
	}
}

func anyLit() cueast.Expr {
	return &cueast.StructLit{Elts: []cueast.Decl{&cueast.Ellipsis{}}}
}
