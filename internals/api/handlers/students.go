package handlers

import (
	"context"

	"github.com/Yash-Kansagara/GoGRPC_API/internals/models"
	"github.com/Yash-Kansagara/GoGRPC_API/internals/repositories/mongodb"
	"github.com/Yash-Kansagara/GoGRPC_API/pkg/utils"
	pb "github.com/Yash-Kansagara/GoGRPC_API/proto/gen"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// GRPC Handler for adding students
func (s *Server) AddStudents(ctx context.Context, req *pb.Students) (*pb.Students, error) {

	// prepare bson data to insert
	newStudents := req.GetStudents()
	newTBson := make([]models.Student, len(newStudents))

	for ind, student := range newStudents {
		temp := &models.Student{}
		err := utils.CopyValues(student, temp)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error copying values")
		}
		newTBson[ind] = *temp
	}

	// insert into mongodb
	res, err := mongodb.StudentsCollection.InsertMany(ctx, newTBson)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error inserting students")
	}

	// populate ids in response (reusing request since mongodb does not return full docs)
	for i, iid := range res.InsertedIDs {
		req.Students[i].Id = iid.(bson.ObjectID).Hex()
	}

	return req, nil
}

// GRPC Handler for deleting students
func (s *Server) DeleteStudents(ctx context.Context, req *pb.StudentIds) (*pb.DeleteStudentResponse, error) {

	// prepare delete bson data
	objectIds := bson.A{}
	for _, id := range req.StudentId {
		if objid, err := bson.ObjectIDFromHex(id); err == nil {
			objectIds = append(objectIds, objid)
		}
	}

	// delete from mongodb
	res, err := mongodb.StudentsCollection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": objectIds}})
	if err != nil {
		return &pb.DeleteStudentResponse{
			Status:       "Fail",
			DeletedCount: int32(0),
		}, utils.ErrorHandler(err, "Error deleting students")
	}

	// return response
	return &pb.DeleteStudentResponse{
		Status:       "Success",
		DeletedCount: int32(res.DeletedCount),
	}, nil
}

// GRPC Handler for getting students
func (s *Server) GetStudents(ctx context.Context, req *pb.GetStudentsReq) (*pb.Students, error) {

	filter := bson.M{}

	reqref := req.Student.ProtoReflect()
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

	cursor, err := mongodb.StudentsCollection.Find(ctx, filter, opt)
	if err != nil {
		return &pb.Students{
			Students: nil,
		}, utils.ErrorHandler(err, "Error finding students")
	}

	students := &pb.Students{
		Students: make([]*pb.Student, 0),
	}

	temp := &models.Student{}

	for cursor.Next(ctx) {

		err := cursor.Decode(temp)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error parsing cursor data")
		}
		pbTemp := &pb.Student{}
		err = utils.CopyValues(temp, pbTemp)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error copying values")
		}
		students.Students = append(students.Students, pbTemp)
	}

	return students, nil
}

// GRPC Handler for updating students
func (s *Server) UpdateStudents(ctx context.Context, req *pb.Students) (*pb.Students, error) {

	var upddated []*pb.Student
	for _, student := range req.GetStudents() {
		temp := &models.Student{}
		err := utils.CopyValues(student, temp)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error copying values")
		}
		id, err := bson.ObjectIDFromHex(student.Id)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error converting id to object id")
		}

		temp.Id = ""
		res, err := mongodb.StudentsCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": temp})
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error Updating object")
		}
		if res.Acknowledged && res.MatchedCount == 1 {
			upddated = append(upddated, student)
		}
	}
	return &pb.Students{
		Students: upddated,
	}, nil
}
