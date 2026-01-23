package handlers

import pb "github.com/Yash-Kansagara/GoGRPC_API/proto/gen"

type Server struct {
	pb.UnimplementedTeacherServiceServer
	pb.UnimplementedStudentServiceServer
	pb.UnimplementedExecServiceServer
}
