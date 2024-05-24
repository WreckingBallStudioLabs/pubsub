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
	CreatedAt time.Time `form:"createdAt" json:"createdAt,omitempty" query:"createdAt"`

	// CreatedBy is the user who created the record.
	CreatedBy string `form:"createdBy" json:"createdBy,omitempty" query:"createdBy" validate:"omitempty,gt=0"`

	// DeleteAt is the time the record was deleted.
	DeleteAt time.Time `form:"deleteAt" json:"deleteAt,omitempty" query:"deleteAt"`

	// DeleteBy is the user who deleted the record.
	DeleteBy string `form:"deleteBy" json:"deleteBy,omitempty" query:"deleteBy" validate:"omitempty,gt=0"`

	// ID is the unique identifier for the record.
	//
	// NOTE: the `id:"uuid"` tag automatically sets with an UUID ONLY if the
	// field is empty.
	ID string `bson:"_id,omitempty" db:"id" dbType:"varchar(255)" form:"id" id:"uuid" json:"id,omitempty" param:"id" query:"id" validate:"omitempty,gt=0"`

	// Queue is the queue to subscribe to, in the form "v1.meta.created.queue".
	// A "queue" is a way to make sure messages are only delivered to one
	// subscriber at a time.
	Queue string `json:"queue" validate:"omitempty,gt=0"`

	// Status is the status of the record.
	Status status.Status `default:"active" form:"status" json:"status,omitempty" query:"status" validate:"omitempty,gt=0"`

	// Topic is the subject to subscribe to, in the form "v1.meta.created".
	// A "topic" is a way to organize messages.
	Topic string `json:"topic" validate:"omitempty,gt=0"`

	// UpdatedAt is the time the record was updated.
	UpdatedAt time.Time `form:"updatedAt" json:"updatedAt,omitempty" query:"updatedAt"`

	// UpdatedBy is the user who updated the record.
	UpdatedBy string `form:"updatedBy" json:"updatedBy,omitempty" query:"updatedBy" validate:"omitempty,gt=0"`
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
