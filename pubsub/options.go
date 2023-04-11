package pubsub

import (
	"github.com/thalesfsp/dal/validation"
)

//////
// Vars, consts, and types.
//////

// Func allows to set options.
type Func func(o *Options) error

// Options for operations.
type Options struct {
	// If the operation is synchronous.
	Sync bool `json:"sync" default:"false" env:"PUBSUB_SYNC"`
}

//////
// Exported built-in options.
//////

// WithSync set the sync option.
func WithSync(sync bool) Func {
	return func(o *Options) error {
		o.Sync = sync

		return nil
	}
}

//////
// Factory.
//////

// NewOptions creates Options.
func NewOptions() (*Options, error) {
	o := &Options{}

	if err := validation.Validate(o); err != nil {
		return nil, err
	}

	return o, nil
}
