package grpc

import (
	"context"

	"github.com/IsFariza/ap2-gRPC-Migration/doctor-service/internal/model"
	"github.com/IsFariza/ap2-gRPC-Migration/doctor-service/internal/model/interfaces"
	doctorpb "github.com/IsFariza/ap2-gRPC-Migration/doctor-service/proto"
)

type doctorHandler struct {
	doctorpb.UnimplementedDoctorServiceServer
	uc interfaces.DoctorUseCase
}

func NewDoctorHandler(uc interfaces.DoctorUseCase) doctorpb.DoctorServiceServer {
	return &doctorHandler{uc: uc}
}

func (h *doctorHandler) CreateDoctor(ctx context.Context, req *doctorpb.CreateDoctorRequest) (*doctorpb.DoctorResponse, error) {
	newDoc := model.Doctor{
		FullName:       req.FullName,
		Specialization: req.Specialization,
		Email:          req.Email,
	}

	savedDoc, err := h.uc.Create(ctx, newDoc)
	if err != nil {
		return nil, mapErrorToStatus(err)
	}
	return &doctorpb.DoctorResponse{
		Id:             savedDoc.ID,
		FullName:       savedDoc.FullName,
		Specialization: savedDoc.Specialization,
		Email:          savedDoc.Email,
	}, nil
}

func (h *doctorHandler) GetDoctor(ctx context.Context, req *doctorpb.GetDoctorRequest) (*doctorpb.DoctorResponse, error) {
	doc, err := h.uc.GetByID(ctx, req.Id)
	if err != nil {
		return nil, mapErrorToStatus(err)
	}
	return &doctorpb.DoctorResponse{
		Id:             doc.ID,
		FullName:       doc.FullName,
		Specialization: doc.Specialization,
		Email:          doc.Email,
	}, nil
}

func (h *doctorHandler) ListDoctors(ctx context.Context, req *doctorpb.ListDoctorsRequest) (*doctorpb.ListDoctorsResponse, error) {
	docs, err := h.uc.List(ctx)
	if err != nil {
		return nil, mapErrorToStatus(err)
	}

	var pbDoctors []*doctorpb.DoctorResponse
	for _, d := range docs {
		pbDoctors = append(pbDoctors,
			&doctorpb.DoctorResponse{
				Id:             d.ID,
				FullName:       d.FullName,
				Specialization: d.Specialization,
				Email:          d.Email,
			})
	}

	return &doctorpb.ListDoctorsResponse{
		Doctors: pbDoctors,
	}, nil
}
