package kubecue

import goast "go/ast"

func (g *Generator) RegisterAny(types ...string) {
	for _, t := range types {
		g.anyTypes[t] = struct{}{}
	}
}

// RegisterTypeFilter registers a filter to filter out top-level types.
func (g *Generator) RegisterTypeFilter(filter func(spec *goast.TypeSpec) bool) {
	g.typeFilter = filter
}
