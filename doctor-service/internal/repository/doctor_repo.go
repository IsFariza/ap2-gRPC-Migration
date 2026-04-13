package repository

import (
	"context"

	"github.com/IsFariza/ap2-gRPC-Migration/doctor-service/internal/model"
	"github.com/IsFariza/ap2-gRPC-Migration/doctor-service/internal/model/interfaces"
	"github.com/IsFariza/ap2-gRPC-Migration/doctor-service/internal/repository/dao"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type doctorRepository struct {
	collection *mongo.Collection
}

func NewDoctorRepository(client *mongo.Client) interfaces.DoctorRepository {
	return &doctorRepository{
		collection: client.Database("doctor_db").Collection("doctors"),
	}
}

func (r *doctorRepository) Create(ctx context.Context, doctor *model.Doctor) error {
	doc := dao.FromDomain(doctor)

	res, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	if objId, ok := res.InsertedID.(primitive.ObjectID); ok {
		doctor.ID = objId.Hex()
	}
	return nil
}

func (r *doctorRepository) GetById(ctx context.Context, id string) (*model.Doctor, error) {
	var doc dao.DoctorDoc
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, model.ErrInvalidID
	}

	filter := bson.M{"_id": objID}

	if err := r.collection.FindOne(ctx, filter).Decode(&doc); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, model.ErrDoctorNotFound
		}
		return nil, err
	}

	return doc.ToDomain(), nil
}

func (r *doctorRepository) GetAll(ctx context.Context) ([]*model.Doctor, error) {
	var docs []dao.DoctorDoc

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	result := make([]*model.Doctor, len(docs))
	for i, doc := range docs {
		result[i] = doc.ToDomain()
	}

	return result, nil
}

func (r *doctorRepository) GetByEmail(ctx context.Context, email string) (*model.Doctor, error) {
	var doc dao.DoctorDoc
	filter := bson.M{"email": email}
	if err := r.collection.FindOne(ctx, filter).Decode(&doc); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, model.ErrDoctorNotFound
		}
		return nil, err
	}
	return doc.ToDomain(), nil
}
