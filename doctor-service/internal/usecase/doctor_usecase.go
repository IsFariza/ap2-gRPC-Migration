package usecase

import (
	"context"

	"github.com/IsFariza/ap2-gRPC-Migration/doctor-service/internal/model"
	"github.com/IsFariza/ap2-gRPC-Migration/doctor-service/internal/model/interfaces"
)

type doctorUseCase struct {
	repo interfaces.DoctorRepository
}

func NewDoctorUseCase(repo interfaces.DoctorRepository) interfaces.DoctorUseCase {
	return &doctorUseCase{
		repo: repo,
	}
}

func (uc *doctorUseCase) Create(ctx context.Context, doc model.Doctor) (*model.Doctor, error) {
	if doc.FullName == "" {
		return nil, model.ErrNameRequired
	}
	if doc.Email == "" {
		return nil, model.ErrEmailRequired
	}

	existing, errEmail := uc.repo.GetByEmail(ctx, doc.Email)
	if existing != nil && errEmail != model.ErrDoctorNotFound {
		return nil, model.ErrEmailUsed
	}
	if existing != nil {
		return nil, model.ErrEmailUsed
	}

	err := uc.repo.Create(ctx, &doc)
	if err != nil {
		return nil, err
	}

	return &doc, nil
}

func (uc *doctorUseCase) GetByID(ctx context.Context, id string) (*model.Doctor, error) {
	doc, err := uc.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, model.ErrDoctorNotFound
	}
	return doc, nil
}

func (uc *doctorUseCase) List(ctx context.Context) ([]*model.Doctor, error) {
	return uc.repo.GetAll(ctx)
}
