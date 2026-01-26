package mongodb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Yash-Kansagara/GoGRPC_API/internals/models"
	"github.com/Yash-Kansagara/GoGRPC_API/pkg/utils"
	pb "github.com/Yash-Kansagara/GoGRPC_API/proto/gen"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func AddExecsToDB(ctx context.Context, req *pb.Execs) (*pb.Execs, error) {

	// prepare bson data to insert
	newExecs := req.GetExecs()
	newTBson := make([]models.Exec, len(newExecs))

	for ind, exec := range newExecs {
		temp := &models.Exec{}
		err := utils.CopyValues(exec, temp)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error copying values")
		}
		temp.Password = utils.HashPassword(temp.Password)
		now, err := time.Now().UTC().MarshalText()
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error getting current time")
		}
		temp.UserCreatedAt = string(now)

		newTBson[ind] = *temp
	}

	// insert into mongodb
	res, err := ExecsCollection.InsertMany(ctx, newTBson)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error inserting execs")
	}

	// populate ids in response (reusing request since mongodb does not return full docs)
	for i, iid := range res.InsertedIDs {
		req.Execs[i].Id = iid.(bson.ObjectID).Hex()
	}

	return req, nil
}

func DeleteExecsFromDB(ctx context.Context, req *pb.ExecIds) (*pb.DeleteExecsResponse, error) {

	// prepare delete bson data
	objectIds := bson.A{}
	for _, id := range req.ExecId {
		if objid, err := bson.ObjectIDFromHex(id); err == nil {
			objectIds = append(objectIds, objid)
		}
	}

	// delete from mongodb
	res, err := ExecsCollection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": objectIds}})
	if err != nil {
		return &pb.DeleteExecsResponse{
			Status:       "Fail",
			DeletedCount: int32(0),
		}, utils.ErrorHandler(err, "Error deleting execs")
	}

	// return response
	return &pb.DeleteExecsResponse{
		Status:       "Success",
		DeletedCount: int32(res.DeletedCount),
	}, nil
}

func GetExecsFromDB(ctx context.Context, req *pb.GetExecsReq) (*pb.Execs, error) {

	filter := bson.M{}

	reqref := req.Exec.ProtoReflect()
	fds := reqref.Descriptor().Fields()

	// add key value from protobuf message to filter
	for i := 0; i < fds.Len(); i++ {
		field := fds.Get(i)
		if reqref.Has(field) {
			filter[string(field.Name())] = reqref.Get(field).Interface()
		}
	}

	// sort
	sort := utils.GetSortBsonDoc(req.SortBy)

	opt := options.Find().SetSort(sort)

	cursor, err := ExecsCollection.Find(ctx, filter, opt)
	if err != nil {
		return &pb.Execs{
			Execs: nil,
		}, utils.ErrorHandler(err, "Error finding execs")
	}

	execs := &pb.Execs{
		Execs: make([]*pb.Exec, 0),
	}

	temp := &models.Exec{}

	for cursor.Next(ctx) {

		err := cursor.Decode(temp)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error parsing cursor data")
		}
		pbTemp := &pb.Exec{}
		err = utils.CopyValues(temp, pbTemp)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error copying values")
		}
		execs.Execs = append(execs.Execs, pbTemp)
	}

	return execs, nil
}

func UpdateExecsInDB(ctx context.Context, req *pb.Execs) (*pb.Execs, error) {

	var upddated []*pb.Exec
	for _, exec := range req.GetExecs() {
		temp := &models.Exec{}
		err := utils.CopyValues(exec, temp)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error copying values")
		}
		id, err := bson.ObjectIDFromHex(exec.Id)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error converting id to object id")
		}

		temp.Id = ""
		res, err := ExecsCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": temp})
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error Updating object")
		}
		if res.Acknowledged && res.MatchedCount == 1 {
			upddated = append(upddated, exec)
		}
	}
	return &pb.Execs{
		Execs: upddated,
	}, nil
}

func LoginExec(ctx context.Context, req *pb.ExecLoginReq) (*pb.ExecLoginRes, error) {

	exec, err := GetExecByUsername(ctx, req.Username)
	if err != nil {
		return &pb.ExecLoginRes{
			Token:   "",
			Success: false,
			Message: "Failed to find exec with username",
		}, utils.ErrorHandler(err, "Failed to find exec with username")
	}

	if int32(exec.InactiveStatus) > 0 {
		return &pb.ExecLoginRes{
			Token:   "",
			Success: false,
			Message: "Exec account is inactive",
		}, utils.ErrorHandler(err, "Exec account is inactive")
	}
	isValid := utils.VerifyPassword(req.Password, exec.Password)
	if !isValid {
		return &pb.ExecLoginRes{
			Token:   "",
			Success: false,
			Message: "Invalid credentials",
		}, utils.ErrorHandler(err, "Invalid credentials")
	}
	tokenStr, err := utils.GenerateJWT(exec.Username, exec.Id, exec.Email)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error generating token")
	}
	return &pb.ExecLoginRes{
		Token:   tokenStr,
		Success: isValid,
		Message: "Login successful",
	}, nil
}

func UpdatePasswordExec(ctx context.Context, req *pb.ExecUpdatePasswordReq) (*pb.ExecUpdatePasswordRes, error) {
	exec, err := GetExecByUsername(ctx, req.Username)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Failed to find exec with username")
	}

	isValid := utils.VerifyPassword(req.CurrentPassword, exec.Password)
	if !isValid {
		return nil, utils.ErrorHandler(err, "Invalid credentials")
	}

	newPasswordHash := utils.HashPassword(req.NewPassword)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error processing new password")
	}
	setDocument := bson.M{"$set": bson.M{"passwoed": newPasswordHash}}

	execIdObj, err := bson.ObjectIDFromHex(exec.Id)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error processing exec id")
	}
	res, err := ExecsCollection.UpdateByID(ctx, execIdObj, setDocument)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error updating exec password")
	}
	if res.MatchedCount == 1 {
		token, err := utils.GenerateJWT(exec.Username, exec.Id, exec.Role)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Failed to create new token for user")
		}
		return &pb.ExecUpdatePasswordRes{
			Success: true,
			Token:   token,
		}, nil
	}
	return &pb.ExecUpdatePasswordRes{
		Success: false,
	}, nil
}

func DeactivateUserExec(ctx context.Context, req *pb.DeactivateUserReq) (*pb.ConfirmationResp, error) {
	exec, err := GetExecByUsername(ctx, req.Username)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Failed to find exec user with username")
	}

	if exec.InactiveStatus == int32(pb.InActiveStatus_INACTIVE) {
		return &pb.ConfirmationResp{
			Success: false,
			Message: "User is already inactive",
		}, nil
	}

	setDocument := bson.M{"$set": bson.M{"inactive_status": int32(pb.InActiveStatus_INACTIVE)}}

	id, err := bson.ObjectIDFromHex(exec.Id)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error processing exec id")
	}
	res, err := ExecsCollection.UpdateByID(ctx, id, setDocument)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error updating exec user")
	}
	if res.MatchedCount == 1 {
		return &pb.ConfirmationResp{
			Success: true,
			Message: "User deactivated successfully",
		}, nil
	}
	return &pb.ConfirmationResp{
		Success: false,
		Message: "User not deactivated",
	}, nil
}

func GetExecByUsername(ctx context.Context, username string) (*models.Exec, error) {
	res := ExecsCollection.FindOne(ctx, bson.M{"username": username})
	if err := res.Err(); err != nil {
		return nil, err
	}

	exec := models.Exec{}
	err := res.Decode(&exec)
	if err != nil {
		return nil, err
	}
	return &exec, nil
}

func GetExecById(ctx context.Context, id string) (*models.Exec, error) {
	idObj, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error processing exec id")
	}
	res := ExecsCollection.FindOne(ctx, bson.M{"_id": idObj})
	if err := res.Err(); err != nil {
		return nil, err
	}

	exec := models.Exec{}
	err = res.Decode(&exec)
	if err != nil {
		return nil, err
	}
	return &exec, nil
}

func GetExecByEmail(ctx context.Context, email string) (*models.Exec, error) {
	res := ExecsCollection.FindOne(ctx, bson.M{"email": email})
	if err := res.Err(); err != nil {
		return nil, err
	}

	exec := models.Exec{}
	err := res.Decode(&exec)
	if err != nil {
		return nil, err
	}
	return &exec, nil
}

func ForgetPasswordExec(ctx context.Context, req *pb.ExecForgetPasswordReq) (*pb.ConfirmationResp, error) {
	exec, err := GetExecByEmail(ctx, req.Email)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Failed to find exec user with email")
	}
	if exec.InactiveStatus == int32(pb.InActiveStatus_INACTIVE) {
		return &pb.ConfirmationResp{
			Success: false,
			Message: "User is inactive",
		}, nil
	}

	tokenStr := utils.GetRandomHash() + "-" + exec.Id
	expiray, _ := time.Now().Add(15 * time.Minute).MarshalText()
	mailBody := fmt.Sprintf("Hello %s,\n\n\tPlease click on the link to reset your password: \n%s", exec.Username, "http://localhost:8080/reset-password?token="+tokenStr)
	setDocument := bson.M{"$set": bson.M{"password_reset_token": tokenStr, "password_reset_token_expires": string(expiray)}}
	execIdObj, err := bson.ObjectIDFromHex(exec.Id)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error processing exec id")
	}
	res, err := ExecsCollection.UpdateByID(ctx, execIdObj, setDocument)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error updating exec user")
	}
	if res.MatchedCount != 1 {
		return &pb.ConfirmationResp{
			Success: false,
			Message: "User not found",
		}, utils.ErrorHandler(err, "Failed to send mail")
	}

	err = utils.SendMail(exec.Email, "Forget Password", mailBody)
	if err != nil {
		return &pb.ConfirmationResp{
			Success: false,
			Message: "Failed to send mail",
		}, utils.ErrorHandler(err, "Failed to send mail")
	}
	return &pb.ConfirmationResp{
		Success: true,
		Message: "Mail sent successfully",
	}, nil
}

func ResetPasswordExec(ctx context.Context, req *pb.ExecResetPasswordReq) (*pb.ConfirmationResp, error) {

	if req.NewPassword != req.ConfirmPassword {
		return &pb.ConfirmationResp{
			Success: false,
			Message: "Passwords do not match",
		}, nil
	}

	splits := strings.Split(req.ResetCode, "-")
	if len(splits) != 2 {
		return &pb.ConfirmationResp{
			Success: false,
			Message: "Invalid reset code",
		}, nil
	}
	id := splits[1]
	requestToken := splits[0]

	exec, err := GetExecById(ctx, id)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Failed to find exec user with id")
	}

	storedToken := strings.Split(exec.PasswordResetToken, "-")[0]
	if storedToken != requestToken {
		return &pb.ConfirmationResp{
			Success: false,
			Message: "Invalid reset code",
		}, nil
	}
	if exec.PasswordResetTokenExpires < time.Now().String() {
		return &pb.ConfirmationResp{
			Success: false,
			Message: "Reset code expired",
		}, nil
	}

	exec.Password = utils.HashPassword(req.NewPassword)
	exec.PasswordResetToken = ""
	exec.PasswordResetTokenExpires = ""
	now, err := time.Now().MarshalText()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error processing time")
	}
	exec.PasswordChangedAt = string(now)

	idObj, err := bson.ObjectIDFromHex(exec.Id)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error processing exec id")
	}
	exec.Id = ""
	res, err := ExecsCollection.UpdateByID(ctx, idObj, bson.M{"$set": exec})
	if err != nil {
		return &pb.ConfirmationResp{
			Success: false,
			Message: "Password not reset",
		}, utils.ErrorHandler(err, "Error updating exec user")
	}
	if res.MatchedCount == 1 {
		return &pb.ConfirmationResp{
			Success: true,
			Message: "Password reset successfully",
		}, nil
	}
	return &pb.ConfirmationResp{
		Success: false,
		Message: "Password not reset",
	}, nil
}
