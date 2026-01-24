package handlers

import (
	"context"

	"github.com/Yash-Kansagara/GoGRPC_API/internals/repositories/mongodb"
	pb "github.com/Yash-Kansagara/GoGRPC_API/proto/gen"
)

// GRPC Handler for adding execs
func (s *Server) AddExecs(ctx context.Context, req *pb.Execs) (*pb.Execs, error) {
	return mongodb.AddExecsToDB(ctx, req)
}

// GRPC Handler for deleting execs
func (s *Server) DeleteExecs(ctx context.Context, req *pb.ExecIds) (*pb.DeleteExecsResponse, error) {
	return mongodb.DeleteExecsFromDB(ctx, req)
}

// GRPC Handler for getting execs
func (s *Server) GetExecs(ctx context.Context, req *pb.GetExecsReq) (*pb.Execs, error) {
	return mongodb.GetExecsFromDB(ctx, req)
}

// rpc UpdateExecs(Execs) returns(Execs);
// GRPC Handler for updating execs
func (s *Server) UpdateExecs(ctx context.Context, req *pb.Execs) (*pb.Execs, error) {
	return mongodb.UpdateExecsInDB(ctx, req)
}
