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

// GRPC Handler for adding teachers
func (s *Server) AddTeachers(ctx context.Context, req *pb.Teachers) (*pb.Teachers, error) {

	// prepare bson data to insert
	newTeachers := req.GetTeachers()
	newTBson := make([]models.Teacher, len(newTeachers))

	for ind, teacher := range newTeachers {
		temp := &models.Teacher{}
		err := utils.CopyValues(teacher, temp)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error copying values")
		}
		newTBson[ind] = *temp
	}

	// insert into mongodb
	res, err := mongodb.TeachersCollection.InsertMany(ctx, newTBson)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error inserting teachers")
	}

	// populate ids in response( reusing request since mongodb does not return full docs)
	for i, iid := range res.InsertedIDs {
		req.Teachers[i].Id = iid.(bson.ObjectID).Hex()
	}

	return req, nil
}

// GRPC Handler for deleting teachers
func (s *Server) DeleteTeachers(ctx context.Context, req *pb.TeacherIds) (*pb.DeleteTeacherResponse, error) {

	// prepare delete bson data
	objectIds := bson.A{}
	for _, id := range req.TeacherIds {
		if objid, err := bson.ObjectIDFromHex(id); err == nil {
			objectIds = append(objectIds, objid)
		}
	}

	// delete from mongodb
	res, err := mongodb.TeachersCollection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": objectIds}})
	if err != nil {
		return &pb.DeleteTeacherResponse{
			Status:       "Fail",
			DeletedCount: int32(0),
		}, utils.ErrorHandler(err, "Error deleting teachers")
	}

	// return response
	return &pb.DeleteTeacherResponse{
		Status:       "Success",
		DeletedCount: int32(res.DeletedCount),
	}, nil
}

// GRPC Handler for getting teachers
func (s *Server) GetTeachers(ctx context.Context, req *pb.GetTeachersReq) (*pb.Teachers, error) {

	filter := bson.M{}
	// t := req.GetTeacher()
	// telem := reflect.ValueOf(t).Elem()

	reqref := req.Teacher.ProtoReflect()
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

	cursor, err := mongodb.TeachersCollection.Find(ctx, filter, opt)
	if err != nil {
		return &pb.Teachers{
			Teachers: nil,
		}, utils.ErrorHandler(err, "Error finding teachers")
	}

	teachers := &pb.Teachers{
		Teachers: make([]*pb.Teacher, 0),
	}

	temp := &models.Teacher{}

	for cursor.Next(ctx) {

		// decoder := bson.
		err := cursor.Decode(temp)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error parsing cursor data")
		}
		pbTemp := &pb.Teacher{}
		err = utils.CopyValues(temp, pbTemp)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error copying values")
		}
		teachers.Teachers = append(teachers.Teachers, pbTemp)
	}

	return teachers, nil
}

// GRPC Handler for updating teachers
func (s *Server) UpdateTeachers(ctx context.Context, req *pb.Teachers) (*pb.Teachers, error) {

	var upddated []*pb.Teacher
	for _, teacher := range req.GetTeachers() {
		temp := &models.Teacher{}
		err := utils.CopyValues(teacher, temp)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error copying values")
		}
		id, err := bson.ObjectIDFromHex(teacher.Id)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error converting id to object id")
		}

		temp.Id = ""
		res, err := mongodb.TeachersCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": temp})
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error Updating object")
		}
		if res.Acknowledged && res.MatchedCount == 1 {
			upddated = append(upddated, teacher)
		}
	}
	return &pb.Teachers{
		Teachers: upddated,
	}, nil
}

// GRPC Handler for getting students by class teacher
func (s *Server) GetStudentsByClassTeacher(ctx context.Context, req *pb.TeacherId) (*pb.Students, error) {

	teacher, err := GetTeacherById(ctx, req.TeacherId)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error finding teacher with given teacher id")
	}

	if len(teacher.Class) == 0 {
		return &pb.Students{
			Students: nil,
		}, utils.ErrorHandler(nil, "Teacher does not have class assigned")
	}

	studentsfilter := bson.M{"class": teacher.Class}

	cursor, err := mongodb.StudentsCollection.Find(ctx, studentsfilter)
	if err != nil {
		return &pb.Students{
			Students: nil,
		}, utils.ErrorHandler(err, "Error finding students with given class")
	}
	defer cursor.Close(ctx)

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
	if err := cursor.Err(); err != nil {
		return &pb.Students{
			Students: nil,
		}, utils.ErrorHandler(err, "Cursor Error finding students with given class")
	}

	return students, nil
}

// GRPC Handler for getting student count by teacher
func (s *Server) GetStudentCountByTeacher(ctx context.Context, req *pb.TeacherId) (*pb.StudentCount, error) {
	teacher, err := GetTeacherById(ctx, req.TeacherId)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error finding teacher with given teacher id")
	}
	if len(teacher.Class) == 0 {
		return &pb.StudentCount{
			Count: 0,
		}, utils.ErrorHandler(nil, "Teacher does not have class assigned")
	}

	filter := bson.M{"class": teacher.Class}
	count, err := mongodb.StudentsCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Failed to fetch count")
	}
	return &pb.StudentCount{
		Count: int32(count),
	}, nil
}

// helper function to get teacher by id
func GetTeacherById(ctx context.Context, teacherId string) (*models.Teacher, error) {
	teacherObjectId, err := bson.ObjectIDFromHex(teacherId)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Invalid teacher id")
	}
	filter := bson.M{"_id": teacherObjectId}
	res := mongodb.TeachersCollection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		return nil, utils.ErrorHandler(err, "Error finding teacher with given teacher id")
	}
	teacher := models.Teacher{}
	if err := res.Decode(&teacher); err != nil {
		return nil, utils.ErrorHandler(err, "Error decoding teacher")
	}
	return &teacher, nil
}
