package handlers

import (
	"context"

	"github.com/Yash-Kansagara/GoGRPC_API/internals/repositories/mongodb"
	pb "github.com/Yash-Kansagara/GoGRPC_API/proto/gen"
)

// GRPC Handler for adding teachers
func (s *Server) AddTeachers(ctx context.Context, req *pb.Teachers) (*pb.Teachers, error) {

	return mongodb.AddTeachersToDB(ctx, req)
}

// GRPC Handler for deleting teachers
func (s *Server) DeleteTeachers(ctx context.Context, req *pb.TeacherIds) (*pb.DeleteTeacherResponse, error) {
	return mongodb.DeleteTeachersFromDB(ctx, req)
}

// GRPC Handler for getting teachers
func (s *Server) GetTeachers(ctx context.Context, req *pb.GetTeachersReq) (*pb.Teachers, error) {
	return mongodb.GetTeachersFromDB(ctx, req)
}

// GRPC Handler for updating teachers
func (s *Server) UpdateTeachers(ctx context.Context, req *pb.Teachers) (*pb.Teachers, error) {
	return mongodb.UpdateTeachersInDB(ctx, req)
}

// GRPC Handler for getting students by class teacher
func (s *Server) GetStudentsByClassTeacher(ctx context.Context, req *pb.TeacherId) (*pb.Students, error) {
	return mongodb.GetStudentsByClassTeacherFromDB(ctx, req)
}

// GRPC Handler for getting student count by teacher
func (s *Server) GetStudentCountByTeacher(ctx context.Context, req *pb.TeacherId) (*pb.StudentCount, error) {
	return mongodb.GetStudentCountByTeacherFromDB(ctx, req)
}
