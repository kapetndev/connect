package recovery

import "context"

var defaultOptions = options{
	recovery: defaultRecoveryFunc,
}

// options describe the full set of options that may be configure to influence
// recovery behaviour.
type options struct {
	recovery RecoveryContextFunc
}

// Option is a function that can configure one or more recovery options.
type Option func(*options)

// RecoveryFunc is a function that recovers from the panic `p` by returning an
// `error`.
type RecoveryFunc func(interface{}) error

// RecoveryContextFunc is a function that recovers from the panic `p` by
// returning an `error`. The context can be used to extract request scoped
// metadata and context values.
type RecoveryContextFunc func(context.Context, interface{}) error

// WithRecovery returns a recovery option to customise a servers recovery
// behaviour from panics.
func WithRecovery(f RecoveryFunc) Option {
	return func(o *options) {
		o.recovery = func(ctx context.Context, p interface{}) error {
			return f(p)
		}
	}
}

// WithRecoveryContext returns a recovery option to customise a servers
// recovery behaviour from panics.
func WithRecoveryContext(f RecoveryContextFunc) Option {
	return func(o *options) {
		o.recovery = f
	}
}

func defaultRecoveryFunc(ctx context.Context, p interface{}) error {
	return nil
}

func applyOptions(opts []Option) options {
	cfg := defaultOptions
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}
