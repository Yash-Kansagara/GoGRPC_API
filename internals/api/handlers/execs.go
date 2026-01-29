package handlers

import (
	"context"

	"github.com/Yash-Kansagara/GoGRPC_API/internals/repositories/mongodb"
	tokendb "github.com/Yash-Kansagara/GoGRPC_API/internals/repositories/token_memory_db"
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

// GRPC Handler for updating execs
func (s *Server) UpdateExecs(ctx context.Context, req *pb.Execs) (*pb.Execs, error) {
	return mongodb.UpdateExecsInDB(ctx, req)
}

// GRPC handler for Logging in exec users
func (s *Server) LoginExec(ctx context.Context, req *pb.ExecLoginReq) (*pb.ExecLoginRes, error) {
	return mongodb.LoginExec(ctx, req)
}

// GRPC handler for Updateing password
func (s *Server) UpdatePasswordExec(ctx context.Context, req *pb.ExecUpdatePasswordReq) (*pb.ExecUpdatePasswordRes, error) {
	return mongodb.UpdatePasswordExec(ctx, req)
}

// GRPC handler to deactivate exec users
func (s *Server) DeactivateUserExec(ctx context.Context, req *pb.DeactivateUserReq) (*pb.ConfirmationResp, error) {
	return mongodb.DeactivateUserExec(ctx, req)
}

// GRPC handler for forget password
func (s *Server) ForgetPasswordExec(ctx context.Context, req *pb.ExecForgetPasswordReq) (*pb.ConfirmationResp, error) {
	return mongodb.ForgetPasswordExec(ctx, req)
}

// GRPC handler for reset password
func (s *Server) ResetPasswordExec(ctx context.Context, req *pb.ExecResetPasswordReq) (*pb.ConfirmationResp, error) {
	return mongodb.ResetPasswordExec(ctx, req)
}

// GRPC handler for refresh token
func (s *Server) RefreshTokenExec(ctx context.Context, req *pb.ExecRefreshTokenReq) (*pb.ExecRefreshTokenRes, error) {
	return tokendb.RefreshTokenExec(ctx, req)
}

// GRPC handler for logging out exec users
func (s *Server) LogoutExec(ctx context.Context, req *pb.ExecLogoutReq) (*pb.ExecLogoutResp, error) {
	return tokendb.LogoutExec(ctx, req)
}
