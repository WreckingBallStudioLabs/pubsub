package pubsub

import (
	"context"
	"expvar"
	"fmt"
	"net/http"

	"github.com/WreckingBallStudioLabs/pubsub/internal/customapm"
	"github.com/WreckingBallStudioLabs/pubsub/internal/logging"
	"github.com/WreckingBallStudioLabs/pubsub/internal/metrics"
	"github.com/thalesfsp/customerror"
	"github.com/thalesfsp/status"
	"github.com/thalesfsp/sypl"
	"github.com/thalesfsp/sypl/level"
	"github.com/thalesfsp/validation"
)

// Type is the type of the entity regarding the framework. It is used to for
// example, to identify the entity in the logs, metrics, and for tracing.
const (
	DefaultMetricCounterLabel = "counter"
	Type                      = "pubsub"

	// Operation name.
	OperationPublish   = "publish"
	OperationSubscribe = "subscribe"
)

// PubSub definition.
type PubSub struct {
	// Logger.
	Logger sypl.ISypl `json:"-" validate:"required"`

	// Name of the pubsub type.
	Name string `json:"name" validate:"required,lowercase,gte=1"`

	// Metricp.
	counterPublish   *expvar.Int `json:"-" validate:"required,gte=0"`
	counterSubscribe *expvar.Int `json:"-" validate:"required,gte=0"`
}

//////
// Implements the IMeta interface.
//////

// GetLogger returns the logger.
func (p *PubSub) GetLogger() sypl.ISypl {
	return p.Logger
}

// GetName returns the storage name.
func (p *PubSub) GetName() string {
	return p.Name
}

// GetType returns its type.
func (p *PubSub) GetType() string {
	return Type
}

// GetPublishCounter returns the counterCount metric.
func (p *PubSub) GetPublishCounter() *expvar.Int {
	return p.counterPublish
}

// GetSubscribeCounter returns the counterCount metric.
func (p *PubSub) GetSubscribeCounter() *expvar.Int {
	return p.counterSubscribe
}

// FailedToErrorMsgHandler trace and logs an already wrapped (customerror) error.
func (p *PubSub) FailedToErrorMsgHandler(ctx context.Context, msg string, err error) error {
	return customapm.TraceError(
		ctx,
		p.GetLogger(),
		customerror.NewFailedToError(msg, customerror.WithError(err)),
	)
}

// ErrorHandler trace and logs an already wrapped (customerror) error.
func (p *PubSub) ErrorHandler(ctx context.Context, operation string, err error) error {
	return customapm.TraceError(
		ctx,
		p.GetLogger(),
		err,
	)
}

// FailedToErrorHandler trace and logs a failedTo error.
func (p *PubSub) FailedToErrorHandler(ctx context.Context, operation string, err error) error {
	return customapm.TraceError(
		ctx,
		p.GetLogger(),
		customerror.NewFailedToError(operation, customerror.WithError(err)),
	)
}

// NotFoundErrorHandler trace and logs a notFound error.
func (p *PubSub) NotFoundErrorHandler(ctx context.Context) error {
	return customapm.TraceError(
		ctx,
		p.GetLogger(),
		customerror.NewHTTPError(http.StatusNotFound),
	)
}

//////
// Factory.
//////

// New returns a new pubsub.
func New(name string) (*PubSub, error) {
	// pubsub's individual logger.
	logger := logging.Get().New(name).SetTags(Type, name)

	a := &PubSub{
		Logger: logger,
		Name:   name,

		counterPublish: metrics.NewInt(
			fmt.Sprintf("%s.%s.%s.%s",
				Type,
				name,
				status.Published,
				DefaultMetricCounterLabel)),
		counterSubscribe: metrics.NewInt(
			fmt.Sprintf("%s.%s.%s.%s",
				Type,
				name,
				status.Subscribed,
				DefaultMetricCounterLabel)),
	}

	// Validate the pubsub.
	if err := validation.Validate(a); err != nil {
		return nil, err
	}

	a.GetLogger().PrintlnWithOptions(
		level.Debug,
		fmt.Sprintf("%+v %s %s", a.GetName(), Type, status.Created),
		sypl.WithTags(Type, string(status.Initialized), a.GetName()),
	)

	return a, nil
}
