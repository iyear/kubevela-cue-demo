package kubecue

type options struct {
	nonNull bool
}

var defaultOptions = &options{
	nonNull: false,
}

type Option func(opts *options)

// WithNonNull will not generate null enum for pointer type
func WithNonNull() Option {
	return func(opts *options) {
		opts.nonNull = true
	}
}
