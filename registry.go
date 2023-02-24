package kubecue

func (g *Generator) RegisterAny(types ...string) {
	for _, t := range types {
		g.anyTypes[t] = struct{}{}
	}
}
