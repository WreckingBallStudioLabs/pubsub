package pubsub

import (
	"context"
	"expvar"
	"fmt"

	"github.com/WreckingBallStudioLabs/pubsub/internal/customapm"
	"github.com/WreckingBallStudioLabs/pubsub/internal/logging"
	"github.com/WreckingBallStudioLabs/pubsub/internal/metrics"
	"github.com/thalesfsp/status"
	"github.com/thalesfsp/sypl"
	"github.com/thalesfsp/sypl/level"
	"github.com/thalesfsp/validation"
)

//////
// Vars, consts, and types.
//////

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

	// Metrics.
	counterInstantiationFailed *expvar.Int `json:"-" validate:"required,gte=0"`
	counterPublished           *expvar.Int `json:"-" validate:"required,gte=0"`
	counterPublishedFailed     *expvar.Int `json:"-" validate:"required,gte=0"`
	counterSubscribed          *expvar.Int `json:"-" validate:"required,gte=0"`
	counterSubscribedFailed    *expvar.Int `json:"-" validate:"required,gte=0"`
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

// GetPublishedCounter returns the metric.
func (p *PubSub) GetPublishedCounter() *expvar.Int {
	return p.counterPublished
}

// GetPublishedFailedCounter returns the metric.
func (p *PubSub) GetPublishedFailedCounter() *expvar.Int {
	return p.counterPublishedFailed
}

// GetSubscribedCounter returns the metric.
func (p *PubSub) GetSubscribedCounter() *expvar.Int {
	return p.counterSubscribed
}

// GetSubscribedFailedCounter returns the metric.
func (p *PubSub) GetSubscribedFailedCounter() *expvar.Int {
	return p.counterSubscribedFailed
}

//////
// Factory.
//////

// New returns a new pubsub.
func New(ctx context.Context, name string) (*PubSub, error) {
	// pubsub's individual logger.
	logger := logging.Get().New(name).SetTags(Type, name)

	a := &PubSub{
		Logger: logger,
		Name:   name,

		counterInstantiationFailed: metrics.NewInt(fmt.Sprintf("%s.%s.%s.%s", Type, name, "instantiation."+status.Failed, DefaultMetricCounterLabel)),
		counterPublished:           metrics.NewInt(fmt.Sprintf("%s.%s.%s.%s", Type, name, status.Published, DefaultMetricCounterLabel)),
		counterPublishedFailed:     metrics.NewInt(fmt.Sprintf("%s.%s.%s.%s", Type, name, status.Published+"."+status.Failed, DefaultMetricCounterLabel)),
		counterSubscribed:          metrics.NewInt(fmt.Sprintf("%s.%s.%s.%s", Type, name, status.Subscribed, DefaultMetricCounterLabel)),
		counterSubscribedFailed:    metrics.NewInt(fmt.Sprintf("%s.%s.%s.%s", Type, name, status.Subscribed+"."+status.Failed, DefaultMetricCounterLabel)),
	}

	// Validate the pubsub.
	if err := validation.Validate(a); err != nil {
		return nil, customapm.TraceError(ctx, err, logger, a.counterInstantiationFailed)
	}

	a.GetLogger().PrintlnWithOptions(
		level.Debug,
		fmt.Sprintf("%+v %s %s", a.GetName(), Type, status.Created),
		sypl.WithTags(Type, string(status.Initialized), a.GetName()),
	)

	return a, nil
}
