package dao

import (
	"time"

	"github.com/IsFariza/ap2-gRPC-Migration/appointment-service/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type AppointmentDoc struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	Title       string        `bson:"title"`
	Description string        `bson:"description"`
	DoctorID    string        `bson:"doctor_id"`
	Status      model.Status  `bson:"status"`
	CreatedAt   time.Time     `bson:"created_at"`
	UpdatedAt   time.Time     `bson:"updated_at"`
}

func FromDomain(d *model.Appointment) *AppointmentDoc {
	doc := &AppointmentDoc{
		Title:       d.Title,
		Description: d.Description,
		DoctorID:    d.DoctorID,
		Status:      d.Status,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   time.Now(),
	}
	if doc.CreatedAt.IsZero() {
		doc.CreatedAt = time.Now()
	}
	if d.ID != "" {
		if objID, err := bson.ObjectIDFromHex(d.ID); err == nil {
			doc.ID = objID
		}
	}
	return doc
}

func (d AppointmentDoc) ToDomain() *model.Appointment {
	return &model.Appointment{
		ID:          d.ID.Hex(),
		Title:       d.Title,
		Description: d.Description,
		DoctorID:    d.DoctorID,
		Status:      d.Status,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}
