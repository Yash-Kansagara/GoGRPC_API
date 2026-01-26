package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Yash-Kansagara/GoGRPC_API/internals/api/handlers"
	"github.com/Yash-Kansagara/GoGRPC_API/internals/api/interceptors"
	"github.com/Yash-Kansagara/GoGRPC_API/internals/repositories/mongodb"
	pb "github.com/Yash-Kansagara/GoGRPC_API/proto/gen"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	// connect DB
	mongoClient, err := mongodb.InitializeMongoDBClient()
	if err == nil {
		defer mongoClient.Disconnect(context.Background())
	}

	// create and serve GRPC services

	interceptors.InitRateLimiter()
	interceptors := grpc.ChainUnaryInterceptor(interceptors.RateLimitingInterceptor, interceptors.ResponseTimeInterceptor)
	grpcServer := grpc.NewServer(interceptors)
	serverInstance := &handlers.Server{}
	pb.RegisterExecServiceServer(grpcServer, serverInstance)
	pb.RegisterStudentServiceServer(grpcServer, serverInstance)
	pb.RegisterTeacherServiceServer(grpcServer, serverInstance)
	reflection.Register(grpcServer)
	port := os.Getenv("GRPC_PORT")
	conn, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal("Failed to create tcp connection", err)
	}
	log.Println("Server listening at", fmt.Sprintf(":%s", port))
	err = grpcServer.Serve(conn)
	if err != nil {
		log.Fatal("Failed serving grpc", err)
	}
}
