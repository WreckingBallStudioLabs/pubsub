package metrics

import (
	"expvar"
	"fmt"
	"os"

	"github.com/WreckingBallStudioLabs/pubsub/internal/logging"
)

//////
// Exported functionalities.
//////

// NewInt creates and initializes a new expvar.Int.
func NewInt(name string) *expvar.Int {
	prefix := os.Getenv("PUBSUB_METRICS_PREFIX")

	if prefix == "" {
		logging.Get().Warnln("PUBSUB_METRICS_PREFIX is not set. Using default (pubsub).")

		prefix = "pubsub"
	}

	counter := expvar.NewInt(
		fmt.Sprintf(
			"%s.%s",
			prefix,
			name,
		),
	)

	counter.Set(0)

	return counter
}
