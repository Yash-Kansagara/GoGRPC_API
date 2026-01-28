package interceptors

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Yash-Kansagara/GoGRPC_API/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var skipMap = map[string]bool{
	"/main.ExecService/LoginExec":          true,
	"/main.ExecService/LogoutExec":         true,
	"/main.ExecService/ResetPasswordExec":  true,
	"/main.ExecService/ForgetPasswordExec": true,
	"/main.ExecService/RefreshTokenExec":   true,
}

func AuthenticatorInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	if skipMap[info.FullMethod] {
		return handler(ctx, req)
	}
	metadata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, utils.ErrorHandler(nil, "Missing metadata")
	}
	authorizationMD := metadata.Get("authorization")
	if len(authorizationMD) == 0 || len(authorizationMD[0]) == 0 {
		return nil, utils.ErrorHandler(nil, "Missing authorization header")
	}
	token := authorizationMD[0]
	if !strings.Contains(token, "Bearer ") {
		return nil, utils.ErrorHandler(nil, "Invalid authorization header")
	}
	token = token[7:]
	fmt.Println("Token:", token)
	jwttoken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))

	if err != nil {
		return nil, utils.ErrorHandler(err, "Invalid token")
	}
	if !jwttoken.Valid {
		return nil, utils.ErrorHandler(err, "Invalid token")
	}
	claims, ok := jwttoken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, utils.ErrorHandler(err, "Invalid token")
	}
	if parsed, err := ParseUserClaims(ctx, &claims); !parsed {
		return nil, utils.ErrorHandler(err, "Invalid token")
	}
	return handler(ctx, req)
}

func ParseUserClaims(ctx context.Context, claims *jwt.MapClaims) (bool, error) {
	mapClaims := *claims
	username, ok := mapClaims["username"].(string)
	if !ok {
		return false, utils.ErrorHandler(nil, "Invalid token")
	}
	userid, ok := mapClaims["userid"].(string)
	if !ok {
		return false, utils.ErrorHandler(nil, "Invalid token")
	}
	role, ok := mapClaims["role"].(string)
	if !ok {
		return false, utils.ErrorHandler(nil, "Invalid token")
	}
	ctx = context.WithValue(ctx, "role", role)
	ctx = context.WithValue(ctx, "username", username)
	ctx = context.WithValue(ctx, "userid", userid)
	return true, nil
}
