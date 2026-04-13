package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IsFariza/ap2-gRPC-Migration/doctor-service/internal/repository"
	"github.com/IsFariza/ap2-gRPC-Migration/doctor-service/internal/transport/grpc"
	"github.com/IsFariza/ap2-gRPC-Migration/doctor-service/internal/usecase"
	doctorpb "github.com/IsFariza/ap2-gRPC-Migration/doctor-service/proto"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	libgrpc "google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found")
	}

	uri := os.Getenv("MONGO_DB")
	port := os.Getenv("PORT")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer func() {
		disconnectCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := client.Disconnect(disconnectCtx); err != nil {
			log.Printf("Error disconnecting Mongo: %v", err)
		}
	}()

	repo := repository.NewDoctorRepository(client)
	uc := usecase.NewDoctorUseCase(repo)
	handler := grpc.NewDoctorHandler(uc)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", port, err)
	}

	server := libgrpc.NewServer()
	doctorpb.RegisterDoctorServiceServer(server, handler)

	log.Printf("Doctor Service running on port %s", port)

	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down gracefully...")

	server.GracefulStop()
	log.Println("Server stopped.")
}
