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
	PubSubErrNilClient       = "PUBSUB_ERR_NIL_MESSAGE"
	PubSubErrPublish         = "PUBSUB_ERR_PUBLISH"
	PubSubErrPubSubNotImpl   = "PUBSUB_ERR_PUBSUB_NOT_IMPL"
	PubSubErrSharedDecode    = "PUBSUB_ERR_SHARED_DECODE"
	PubSubErrSharedEncode    = "PUBSUB_ERR_SHARED_ENCODE"
	PubSubErrSharedMarshal   = "PUBSUB_ERR_SHARED_MARSHAL"
	PubSubErrSharedRead      = "PUBSUB_ERR_SHARED_READ"
	PubSubErrSharedUnmarshal = "PUBSUB_ERR_SHARED_UNMARSHAL"
	PubSubErrSubscribe       = "PUBSUB_ERR_SUBSCRIBE"
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

		catalog.MustSet(PubSubErrNameName, "name. It should be like `v1.meta.created` or `v1.meta.created.queue`")
		catalog.MustSet(PubSubErrNilClient, "get client, it's nil. Call `New`")
		catalog.MustSet(PubSubErrPublish, "publish")
		catalog.MustSet(PubSubErrPubSubNotImpl, "not implemented")
		catalog.MustSet(PubSubErrSharedDecode, "decode")
		catalog.MustSet(PubSubErrSharedEncode, "encode")
		catalog.MustSet(PubSubErrSharedMarshal, "marshal")
		catalog.MustSet(PubSubErrSharedRead, "read")
		catalog.MustSet(PubSubErrSharedUnmarshal, "unmarshal")
		catalog.MustSet(PubSubErrSubscribe, "subscribe")

		singleton = catalog
	})

	return singleton
}
