package kubecue

import (
	"reflect"
	"strings"
)

type TagOptions struct {
	Name     string
	Inline   bool
	Optional bool
}

// TODO(iyear): be customizable
const (
	basicTag = "json"
	extTag   = "cue"
)

func (g *Generator) parseTag(tag string) *TagOptions {
	if tag == "" {
		return &TagOptions{}
	}

	name, opts := parseTag(reflect.StructTag(tag).Get(basicTag))
	// TODO(iyear): support extTag

	return &TagOptions{
		Name:     name,
		Inline:   opts.Has("inline"),
		Optional: opts.Has("omitempty"),
	}
}

type tagOptions string

func parseTag(tag string) (string, tagOptions) {
	tag, opt, _ := strings.Cut(tag, ",")
	return tag, tagOptions(opt)
}

func (o tagOptions) Has(opt string) bool {
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var name string
		name, s, _ = strings.Cut(s, ",")
		if name == opt {
			return true
		}
	}
	return false
}
