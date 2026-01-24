package mongodb

import (
	"context"

	"github.com/Yash-Kansagara/GoGRPC_API/internals/models"
	"github.com/Yash-Kansagara/GoGRPC_API/pkg/utils"
	pb "github.com/Yash-Kansagara/GoGRPC_API/proto/gen"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// GRPC Handler for adding execs
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

// GRPC Handler for deleting execs
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

// GRPC Handler for getting execs
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

// rpc UpdateExecs(Execs) returns(Execs);
// GRPC Handler for updating execs
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
