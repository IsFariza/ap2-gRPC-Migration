package grpc

import (
	"context"
	"time"

	pb "github.com/IsFariza/ap2-gRPC-Migration/appointment-service/appointment_proto"
	"github.com/IsFariza/ap2-gRPC-Migration/appointment-service/internal/model"
	"github.com/IsFariza/ap2-gRPC-Migration/appointment-service/internal/model/interfaces"
)

type appointmentHandler struct {
	pb.UnimplementedAppointmentServiceServer
	uc interfaces.AppointmentUsecase
}

func NewAppointmentHandler(uc interfaces.AppointmentUsecase) *appointmentHandler {
	return &appointmentHandler{
		uc: uc,
	}
}

func (h *appointmentHandler) CreateAppointment(ctx context.Context, req *pb.CreateAppointmentRequest) (*pb.AppointmentResponse, error) {
	appt := &model.Appointment{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		DoctorID:    req.GetDoctorId(),
	}

	err := h.uc.Create(ctx, appt)
	if err != nil {
		return nil, mapErrorToStatus(err)
	}

	return h.mapToProto(appt), nil
}

func (h *appointmentHandler) GetAppointment(ctx context.Context, req *pb.GetAppointmentRequest) (*pb.AppointmentResponse, error) {
	appt, err := h.uc.GetById(ctx, req.Id)
	if err != nil {
		return nil, mapErrorToStatus(err)
	}

	return h.mapToProto(appt), nil
}

func (h *appointmentHandler) ListAppointments(ctx context.Context, req *pb.ListAppointmentsRequest) (*pb.ListAppointmentsResponse, error) {
	appts, err := h.uc.GetAll(ctx)
	if err != nil {
		return nil, mapErrorToStatus(err)
	}

	var pbAppts []*pb.AppointmentResponse
	for _, appt := range appts {
		pbAppts = append(pbAppts, h.mapToProto(appt))
	}
	return &pb.ListAppointmentsResponse{
		Appointments: pbAppts,
	}, nil
}

func (h *appointmentHandler) UpdateAppointmentStatus(ctx context.Context, req *pb.UpdateStatusRequest) (*pb.AppointmentResponse, error) {
	err := h.uc.Update(ctx, req.Id, model.Status(req.Status))
	if err != nil {
		return nil, mapErrorToStatus(err)
	}

	updated, err := h.uc.GetById(ctx, req.Id)
	if err != nil {
		return nil, mapErrorToStatus(err)
	}
	return h.mapToProto(updated), nil
}

func (h *appointmentHandler) mapToProto(appt *model.Appointment) *pb.AppointmentResponse {
	return &pb.AppointmentResponse{
		Id:          appt.ID,
		Title:       appt.Title,
		Description: appt.Description,
		DoctorId:    appt.DoctorID,
		Status:      string(appt.Status),
		CreatedAt:   appt.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   appt.UpdatedAt.Format(time.RFC3339),
	}
}
