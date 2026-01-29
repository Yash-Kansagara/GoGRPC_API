package tokendb

import (
	"context"
	"strings"

	"github.com/Yash-Kansagara/GoGRPC_API/pkg/utils"
	pb "github.com/Yash-Kansagara/GoGRPC_API/proto/gen"
)

func RefreshTokenExec(ctx context.Context, req *pb.ExecRefreshTokenReq) (*pb.ExecRefreshTokenRes, error) {
	current := req.RefreshToken

	// parse token, this checks expiray and signing method
	refreshTokenClaims, _, err := utils.ParseRefreshToken(current)

	if err != nil {
		return nil, utils.ErrorHandler(err, "Invalid Refresh token")
	}

	// validate if refresh token exist in cache, i.e. not expired
	signature := strings.Split(current, ".")[2]
	_, ok := GetToken(signature)
	if !ok {
		return nil, utils.ErrorHandler(nil, "Invalid Refresh token")
	}

	// generate new access token
	newAccessToken, err := utils.GenerateAccessToken(refreshTokenClaims.Username, refreshTokenClaims.UserId, refreshTokenClaims.Role)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error generating new access token")
	}

	RemoveToken(signature)
	newRefreshToken, _, err := utils.GenerateRefreshToken(refreshTokenClaims.Username, refreshTokenClaims.UserId, refreshTokenClaims.Role)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error generating new refresh token")
	}

	err = AddToken(newRefreshToken)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error generating new refresh token")
	}

	return &pb.ExecRefreshTokenRes{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func LogoutExec(ctx context.Context, req *pb.ExecLogoutReq) (*pb.ExecLogoutResp, error) {
	RemoveToken(strings.Split(req.Token, ".")[2])
	return &pb.ExecLogoutResp{
		LoggedOut: true,
	}, nil
}
