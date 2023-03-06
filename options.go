package kubecue

import goast "go/ast"

type options struct {
	nonNull    bool
	typeFilter func(spec *goast.TypeSpec) bool
	anyTypes   map[string]struct{}
}

var defaultOptions = &options{
	nonNull:    false,
	typeFilter: func(_ *goast.TypeSpec) bool { return true },
	anyTypes: map[string]struct{}{
		"map[string]interface{}": {}, "map[string]any": {},
		"interface{}": {}, "any": {},
	},
}

type Option func(opts *options)

// WithNonNull will not generate null enum for pointer type
func WithNonNull() Option {
	return func(opts *options) {
		opts.nonNull = true
	}
}

// WithTypeFilter will filter out the type that not need to be generated
//
// Default filter will generate all types
func WithTypeFilter(f func(spec *goast.TypeSpec) bool) Option {
	return func(opts *options) {
		opts.typeFilter = f
	}
}

func WithAnyTypes(types ...string) Option {
	return func(opts *options) {
		for _, t := range types {
			opts.anyTypes[t] = struct{}{}
		}
	}
}
