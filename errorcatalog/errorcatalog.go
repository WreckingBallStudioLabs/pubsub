package errorcatalog

import (
	"sync"

	"github.com/thalesfsp/customerror"
)

//////
// Vars, consts, and types.
//////

// Singleton.
var (
	once      sync.Once
	singleton *customerror.Catalog
)

const (
	PubSubErrNameName        = "PUBSUB_ERR_NAME_NAME"
	PubSubErrNatsNilMessage  = "PUBSUB_ERR_NATS_NIL_MESSAGE"
	PubSubErrNatsPublish     = "PUBSUB_ERR_NATS_PUBLISH"
	PubSubErrNatsSubscribe   = "PUBSUB_ERR_NATS_SUBSCRIBE"
	PubSubErrSharedDecode    = "PUBSUB_ERR_SHARED_DECODE"
	PubSubErrSharedEncode    = "PUBSUB_ERR_SHARED_ENCODE"
	PubSubErrSharedMarshal   = "PUBSUB_ERR_SHARED_MARSHAL"
	PubSubErrSharedRead      = "PUBSUB_ERR_SHARED_READ"
	PubSubErrSharedUnmarshal = "PUBSUB_ERR_SHARED_UNMARSHAL"
)

//////
// Exported functionalities.
//////

// Get the singleton.
func Get() *customerror.Catalog {
	// Setup once.
	once.Do(func() {
		catalog, err := customerror.NewCatalog("pubsub")
		if err != nil {
			panic(err)
		}

		//////
		// Add error codes.
		//////

		catalog.MustSet(PubSubErrNameName, "name")
		catalog.MustSet(PubSubErrNatsNilMessage, "get client, it's nil. Call `New`")
		catalog.MustSet(PubSubErrNatsPublish, "publish")
		catalog.MustSet(PubSubErrNatsSubscribe, "subscribe")
		catalog.MustSet(PubSubErrSharedDecode, "decode")
		catalog.MustSet(PubSubErrSharedEncode, "encode")
		catalog.MustSet(PubSubErrSharedMarshal, "marshal")
		catalog.MustSet(PubSubErrSharedRead, "read")
		catalog.MustSet(PubSubErrSharedUnmarshal, "unmarshal")

		singleton = catalog
	})

	return singleton
}
