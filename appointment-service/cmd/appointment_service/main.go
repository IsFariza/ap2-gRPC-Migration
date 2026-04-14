package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/IsFariza/ap2-gRPC-Migration/appointment-service/appointment_proto"
	"github.com/IsFariza/ap2-gRPC-Migration/appointment-service/internal/client"
	"github.com/IsFariza/ap2-gRPC-Migration/appointment-service/internal/repository"
	transportgrpc "github.com/IsFariza/ap2-gRPC-Migration/appointment-service/internal/transport/grpc"
	"github.com/IsFariza/ap2-gRPC-Migration/appointment-service/internal/usecase"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}

	mongoURI := os.Getenv("MONGO_DB")
	grpcPort := os.Getenv("PORT")
	doctorAddr := os.Getenv("DOCTOR_ADDR")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mongoClient, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("MongoDB Connection Error: %v", err)
	}

	conn, err := grpc.Dial(
		doctorAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to initialize Doctor Service client: %v", err)
	}
	defer conn.Close()

	docClient := client.NewDoctorClient(conn)
	repo := repository.NewAppointmentRepository(mongoClient)
	apptUsecase := usecase.NewAppointmentUsecase(repo, docClient)
	handler := transportgrpc.NewAppointmentHandler(apptUsecase)

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterAppointmentServiceServer(server, handler)

	go func() {
		log.Printf("Appointment service starting on port %s", grpcPort)
		if err := server.Serve(lis); err != nil {
			log.Fatalf("gRPC Server Error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down gracefully...")

	server.GracefulStop()

	disconnectCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := mongoClient.Disconnect(disconnectCtx); err != nil {
		log.Fatalf("MongoDB Disconnect Error: %v", err)
	}

	log.Println("Service stopped safely.")
}
