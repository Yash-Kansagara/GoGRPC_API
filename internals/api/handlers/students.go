package handlers

import (
	"context"

	"github.com/Yash-Kansagara/GoGRPC_API/internals/repositories/mongodb"
	pb "github.com/Yash-Kansagara/GoGRPC_API/proto/gen"
)

// GRPC Handler for adding students
func (s *Server) AddStudents(ctx context.Context, req *pb.Students) (*pb.Students, error) {
	return mongodb.AddStudentsToDB(ctx, req)
}

// GRPC Handler for deleting students
func (s *Server) DeleteStudents(ctx context.Context, req *pb.StudentIds) (*pb.DeleteStudentResponse, error) {
	return mongodb.DeleteStudentsFromDB(ctx, req)
}

// GRPC Handler for getting students
func (s *Server) GetStudents(ctx context.Context, req *pb.GetStudentsReq) (*pb.Students, error) {
	return mongodb.GetStudentsFromDB(ctx, req)
}

// GRPC Handler for updating students
func (s *Server) UpdateStudents(ctx context.Context, req *pb.Students) (*pb.Students, error) {
	return mongodb.UpdateStudentsInDB(ctx, req)
}
