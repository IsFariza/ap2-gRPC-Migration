package dao

import (
	"time"

	"github.com/IsFariza/ap2-gRPC-Migration/doctor-service/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DoctorDoc struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	FullName       string             `bson:"full_name"`
	Specialization string             `bson:"specialization,omitempty"`
	Email          string             `bson:"email"`
	CreatedAt      time.Time          `bson:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at"`
}

func FromDomain(d *model.Doctor) *DoctorDoc {
	doc := &DoctorDoc{
		FullName:       d.FullName,
		Specialization: d.Specialization,
		Email:          d.Email,
		CreatedAt:      d.CreatedAt,
		UpdatedAt:      time.Now(),
	}
	if doc.CreatedAt.IsZero() {
		doc.CreatedAt = time.Now()
	}
	if d.ID != "" {
		if objID, err := primitive.ObjectIDFromHex(d.ID); err == nil {
			doc.ID = objID
		}
	}
	return doc
}

func (d DoctorDoc) ToDomain() *model.Doctor {
	return &model.Doctor{
		ID:             d.ID.Hex(),
		FullName:       d.FullName,
		Specialization: d.Specialization,
		Email:          d.Email,
		CreatedAt:      d.CreatedAt,
		UpdatedAt:      d.UpdatedAt,
	}
}
