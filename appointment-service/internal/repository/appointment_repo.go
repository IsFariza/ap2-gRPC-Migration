package repository

import (
	"context"
	"time"

	"github.com/IsFariza/ap2-gRPC-Migration/appointment-service/internal/model"
	"github.com/IsFariza/ap2-gRPC-Migration/appointment-service/internal/model/interfaces"
	"github.com/IsFariza/ap2-gRPC-Migration/appointment-service/internal/repository/dao"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type appointmentRepo struct {
	collection *mongo.Collection
}

func NewAppointmentRepository(client *mongo.Client) interfaces.AppointmentRepo {
	return &appointmentRepo{
		collection: client.Database("appointment_db").Collection("appointments"),
	}
}

func (r *appointmentRepo) Create(ctx context.Context, appointment *model.Appointment) error {
	appointment.CreatedAt = time.Now()
	appointment.UpdatedAt = time.Now()

	doc := dao.FromDomain(appointment)
	res, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	if objId, ok := res.InsertedID.(bson.ObjectID); ok {
		appointment.ID = objId.Hex()
	}
	return nil
}

func (r *appointmentRepo) GetById(ctx context.Context, id string) (*model.Appointment, error) {
	var doc dao.AppointmentDoc
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, model.ErrInvalidID
	}
	filter := bson.M{"_id": objID}

	if err := r.collection.FindOne(ctx, filter).Decode(&doc); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, model.ErrAppointmentNotFound
		}
		return nil, err
	}
	return doc.ToDomain(), nil
}

func (r *appointmentRepo) GetAll(ctx context.Context) ([]*model.Appointment, error) {
	var docs []dao.AppointmentDoc

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	result := make([]*model.Appointment, len(docs))
	for i, doc := range docs {
		result[i] = doc.ToDomain()
	}

	return result, nil
}

func (r *appointmentRepo) Update(ctx context.Context, id string, newStatus model.Status) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	doc := dao.FromDomain(&model.Appointment{
		ID:     id,
		Status: newStatus,
	})
	update := bson.M{"$set": bson.M{
		"status":     doc.Status,
		"updated_at": doc.UpdatedAt,
	}}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return model.ErrAppointmentNotFound
	}
	return nil
}
