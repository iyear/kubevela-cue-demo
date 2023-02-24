package kubecue

import (
	"reflect"
	"strings"
)

type TagOptions struct {
	// basic
	Name     string
	Inline   bool
	Optional bool

	// extension
	Default *string // nil means no default value
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
	ext := parseExtTag(reflect.StructTag(tag).Get(extTag))

	return &TagOptions{
		Name:     name,
		Inline:   opts.Has("inline"),
		Optional: opts.Has("omitempty"),

		Default: ext.Get("default"),
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

func parseExtTag(str string) extTagOptions {
	sep := ";"
	settings := map[string]string{}
	names := strings.Split(str, sep)

	for i := 0; i < len(names); i++ {
		j := i
		if len(names[j]) > 0 {
			for {
				if names[j][len(names[j])-1] == '\\' {
					i++
					names[j] = names[j][0:len(names[j])-1] + sep + names[i]
					names[i] = ""
				} else {
					break
				}
			}
		}

		values := strings.Split(names[j], ":")
		k := strings.TrimSpace(strings.ToLower(values[0]))

		if len(values) >= 2 {
			settings[k] = strings.Join(values[1:], ":")
		} else if k != "" {
			settings[k] = k
		}
	}

	return settings
}

type extTagOptions map[string]string

func (e extTagOptions) Get(key string) *string {
	if v, ok := e[key]; ok {
		return &v
	}
	return nil
}
