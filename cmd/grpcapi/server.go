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
	tokendb "github.com/Yash-Kansagara/GoGRPC_API/internals/repositories/token_memory_db"

	pb "github.com/Yash-Kansagara/GoGRPC_API/proto/gen"
	"github.com/joho/godotenv"
)

func main() {

	// load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	// init in-memory token db
	tokendb.Init()

	// connect DB
	mongoClient, err := mongodb.InitializeMongoDBClient()
	if err == nil {
		defer mongoClient.Disconnect(context.Background())
	}

	// create interceptors
	interceptors.InitRateLimiter()
	interceptors := grpc.ChainUnaryInterceptor(
		interceptors.RateLimitingInterceptor,
		interceptors.ResponseTimeInterceptor,
		interceptors.AuthenticatorInterceptor,
	)

	// create and serve GRPC services
	grpcServer := grpc.NewServer(interceptors)
	serverInstance := &handlers.Server{}
	pb.RegisterExecServiceServer(grpcServer, serverInstance)
	pb.RegisterStudentServiceServer(grpcServer, serverInstance)
	pb.RegisterTeacherServiceServer(grpcServer, serverInstance)

	// for development
	reflection.Register(grpcServer)

	// start server
	port := os.Getenv("GRPC_PORT")
	conn, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal("Failed to create tcp connection", err)
	}
	log.Println("Server listening at", fmt.Sprintf(":%s", port))
	err = grpcServer.Serve(conn) // this will block if succeeds
	if err != nil {
		log.Fatal("Failed serving grpc", err)
	}

	// close in-memory token db
	tokendb.Cache.Close()
}
