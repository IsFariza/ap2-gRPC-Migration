package usecase

import (
	"context"

	"github.com/IsFariza/ap2-gRPC-Migration/appointment-service/internal/model"
	"github.com/IsFariza/ap2-gRPC-Migration/appointment-service/internal/model/interfaces"
)

type appointmentUsecase struct {
	repo         interfaces.AppointmentRepo
	doctorClient interfaces.DoctorClient
}

func NewAppointmentUsecase(repo interfaces.AppointmentRepo, dc interfaces.DoctorClient) interfaces.AppointmentUsecase {
	return &appointmentUsecase{
		repo:         repo,
		doctorClient: dc,
	}
}

func (uc *appointmentUsecase) Create(ctx context.Context, appt *model.Appointment) error {
	if appt.Title == "" {
		return model.ErrTitleRequired
	}
	if appt.DoctorID == "" {
		return model.ErrDoctorIDRequired
	}

	exists, err := uc.doctorClient.DoctorExists(ctx, appt.DoctorID)
	if err != nil {
		return model.ErrDoctorServiceUnavailable
	}
	if !exists {
		return model.ErrDoctorNotFoundRemote
	}

	appt.Status = model.StatusNew
	return uc.repo.Create(ctx, appt)
}

func (uc *appointmentUsecase) GetById(ctx context.Context, id string) (*model.Appointment, error) {
	if id == "" {
		return nil, model.ErrInvalidID
	}
	return uc.repo.GetById(ctx, id)
}

func (uc *appointmentUsecase) GetAll(ctx context.Context) ([]*model.Appointment, error) {
	return uc.repo.GetAll(ctx)
}

func (uc *appointmentUsecase) Update(ctx context.Context, id string, newStatus model.Status) error {
	currentStatus, err := uc.repo.GetById(ctx, id)
	if id == "" {
		return model.ErrInvalidID
	}
	if err != nil {
		return err
	}

	if currentStatus.Status == model.StatusDone && newStatus == model.StatusNew {
		return model.ErrInvalidStatusTransition
	}
	if newStatus != model.StatusDone && newStatus != model.StatusInProgress &&
		newStatus != model.StatusNew {
		return model.ErrInvalidStatus
	}
	return uc.repo.Update(ctx, id, newStatus)
}
