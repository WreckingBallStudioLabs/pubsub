package customapm

import (
	"context"
	"fmt"

	"go.elastic.co/apm"
)

//////
// Exported functionalities.
//////

// TXFromCtx Creates a new TX if none is found in the context, otherwise reuses
// the existing one.
//
// Note on TX and Span naming:
//
// The name and type of a transaction and span depend on the specific operation
// being performed by the request. Here's an example of how you might name and
// type a transaction and span for an incoming request to insert data into a db:
//
// Transaction:
//
// - Name: "post-user"
// - Type: "request.post.user"
//
// Span:
//
// - Name: "insert-user"
// - Type: "db.sql.insert"
//
// In this example, the tx was named "post-user" to describe the operation being
// performed "request.post.user".
//
// For the span, it was named "insert-user" to describe a specific operation
// being performed by the span. It was also categorized as "db.sql.insert",
// which indicates that it involves a database operation.
//
// If the span type contains two dots, they are assumed to separate the type and
// subtype parts of the span type. The action will not be set in this case.
//
// For example, if you use a span type of "db.sql.insert", this indicates that
// the span represents a database operation of type "db", subtype "sql", and
// action "insert". The StartSpan() method will automatically parse the span
// type string and set the appropriate values for the Type, Subtype, and Action
// fields of the SpanData object.
//
// If you use a span type of "db.sql", this indicates that the span represents
// a database operation of type "db" and subtype "sql". The Action field of
// the SpanData object will be left blank in this case.
//
// If you use a span type of "db", this indicates that the span represents a
// generic database operation of type "db". Both the Subtype and Action fields
// of the SpanData object will be left blank in this case.
func TXFromCtx(ctx context.Context, txName string, txType string) *apm.Transaction {
	tx := apm.TransactionFromContext(ctx)
	if tx == nil {
		tx = apm.DefaultTracer.StartTransaction(txName, txType)
	}

	return tx
}

// Trace will trace an operation. It uses the existing TX otherwise it fallback
// creating a new TX then it creates a new span within the TX.
//
// NOTE: It's up to the developer to call `span.End()`.
//
// Note on TX and Span naming:
//
// The name and type of a transaction and span depend on the specific operation
// being performed by the request. Here's an example of how you might name and
// type a transaction and span for an incoming request to insert data into a db:
//
// Transaction:
//
// - Name: "post-user"
// - Type: "request.post.user"
//
// Span:
//
// - Name: "insert-user"
// - Type: "db.sql.insert"
//
// In this example, the tx was named "post-user" to describe the operation being
// performed "request.post.user".
//
// For the span, it was named "insert-user" to describe a specific operation
// being performed by the span. It was also categorized as "db.sql.insert",
// which indicates that it involves a database operation.
//
// If the span type contains two dots, they are assumed to separate the type and
// subtype parts of the span type. The action will not be set in this case.
//
// For example, if you use a span type of "db.sql.insert", this indicates that
// the span represents a database operation of type "db", subtype "sql", and
// action "insert". The StartSpan() method will automatically parse the span
// type string and set the appropriate values for the Type, Subtype, and Action
// fields of the SpanData object.
//
// If you use a span type of "db.sql", this indicates that the span represents
// a database operation of type "db" and subtype "sql". The Action field of
// the SpanData object will be left blank in this case.
//
// If you use a span type of "db", this indicates that the span represents a
// generic database operation of type "db". Both the Subtype and Action fields
// of the SpanData object will be left blank in this case.
func Trace(
	ctx context.Context,
	what, nameOf, operation string,
) (context.Context, *apm.Span) {
	//////
	// APM.
	//////

	tx := TXFromCtx(ctx, nameOf, what)

	ctx = apm.ContextWithTransaction(ctx, tx)

	span, ctx := apm.StartSpan(
		ctx,
		fmt.Sprintf("%s.%s", nameOf, operation),
		fmt.Sprintf("%s.%s.%s", what, nameOf, operation),
	)

	return ctx, span
}
