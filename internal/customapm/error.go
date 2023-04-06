package customapm

import (
	"context"
	"errors"

	"github.com/WreckingBallStudioLabs/pubsub/internal/logging"
	"github.com/thalesfsp/sypl"
	"github.com/thalesfsp/sypl/fields"
	"github.com/thalesfsp/sypl/level"
	"go.elastic.co/apm"
)

//////
// Vars, consts, and types.
//////

// Outcome is the outcome of a span. It can be either success or failure.
type Outcome string

const (
	// Failure is the outcome when the span failed.
	Failure Outcome = "failure"

	// Success is the outcome when the span succeeded.
	Success Outcome = "success"
)

func (o Outcome) String() string {
	return string(o)
}

//////
// Exported functionalities.
//////

// TraceError is a helper function to trace an error. It will log the error
// with the APM fields, and tell APM that it was an error. It will also set
// the span outcome to failure.
func TraceError(ctx context.Context, l sypl.ISypl, err error) error {
	originalError := err

	// Get the current span from the context, if any, set the outcome.
	span := apm.SpanFromContext(ctx)
	if span != nil {
		span.Outcome = string(Failure)
	}

	// Unwrap any nested errorcatalog. By default, apm.CaptureError() does not
	// automatically unwrap nested errors or extract any additional context or
	// metadata from the error.
	for {
		unwrapped := errors.Unwrap(err)
		if unwrapped == nil {
			break
		}

		err = unwrapped
	}

	// Correlates the transaction, span and log, and logs it.
	l.PrintlnWithOptions(
		level.Error,
		err.Error(),
		sypl.WithFields(logging.ToAPM(ctx, make(fields.Fields))),
	)

	// Tells APM that it that was an error.
	apm.CaptureError(ctx, err).Send()

	return originalError
}
