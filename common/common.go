package common

import (
	"time"

	"github.com/WreckingBallStudioLabs/pubsub/internal/shared"
	"github.com/thalesfsp/status"
	"github.com/thalesfsp/validation"
)

// Common contains common fields across all models.
type Common struct {
	// CreatedAt is the time the record was created.
	CreatedAt time.Time `json:"createdAt,omitempty" form:"createdAt" query:"createdAt"`

	// CreatedBy is the user who created the record.
	CreatedBy string `json:"createdBy,omitempty" form:"createdBy" query:"createdBy" validate:"omitempty,gt=0"`

	// DeleteAt is the time the record was deleted.
	DeleteAt time.Time `json:"deleteAt,omitempty" form:"deleteAt" query:"deleteAt"`

	// DeleteBy is the user who deleted the record.
	DeleteBy string `json:"deleteBy,omitempty" form:"deleteBy" query:"deleteBy" validate:"omitempty,gt=0"`

	// ID is the unique identifier for the record.
	//
	// NOTE: the `id:"uuid"` tag automatically sets with an UUID ONLY if the
	// field is empty.
	ID string `json:"id,omitempty" id:"uuid" form:"id" param:"id" query:"id" db:"id" dbType:"varchar(255)" bson:"_id,omitempty" validate:"omitempty,gt=0"`

	// Queue is the queue to subscribe to, in the form "v1.meta.created.queue".
	// A "queue" is a way to make sure messages are only delivered to one
	// subscriber at a time.
	Queue string `json:"queue" validate:"omitempty,gt=0"`

	// Status is the status of the record.
	Status status.Status `json:"status,omitempty" form:"status" query:"status" validate:"omitempty,gt=0" default:"active"`

	// Topic is the subject to subscribe to, in the form "v1.meta.created".
	// A "topic" is a way to organize messages.
	Topic string `json:"topic" validate:"omitempty,gt=0"`

	// UpdatedAt is the time the record was updated.
	UpdatedAt time.Time `json:"updatedAt,omitempty" form:"updatedAt" query:"updatedAt"`

	// UpdatedBy is the user who updated the record.
	UpdatedBy string `json:"updatedBy,omitempty" form:"updatedBy" query:"updatedBy" validate:"omitempty,gt=0"`
}

//////
// Factory.
//////

// New creates a new Common with initialized default values.
func New() (*Common, error) {
	common := &Common{
		ID:        shared.GenerateUUID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := validation.Validate(common); err != nil {
		return nil, err
	}

	return common, nil
}
