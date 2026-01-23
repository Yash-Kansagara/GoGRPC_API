package utils

import (
	pb "github.com/Yash-Kansagara/GoGRPC_API/proto/gen"
	"go.mongodb.org/mongo-driver/v2/bson"
	"google.golang.org/protobuf/proto"
)

// GetSortBsonDoc converts a slice of SortField to a bson.D for sorting
func GetSortBsonDoc(sortFields []*pb.SortField) bson.D {
	sort := bson.D{}
	for _, s := range sortFields {

		order := int32(1)
		if s.Order == pb.Order_DESC {
			order = int32(-1)
		}
		sort = append(sort, bson.E{Key: s.Field, Value: order})
	}
	return sort
}

// GetFilterBsonDoc converts a protobuf message to a bson.M for filtering
func GetFilterBsonDoc[T *proto.Message](filters T) bson.M {
	filter := bson.M{}
	if filters == nil {
		return filter
	}
	reqref := (*filters).ProtoReflect()
	fds := reqref.Descriptor().Fields()

	// add key value from protobuf message to filter
	for i := 0; i < fds.Len(); i++ {
		field := fds.Get(i)
		if reqref.Has(field) && reqref.Get(field).IsValid() {
			filter[string(field.Name())] = reqref.Get(field).Interface()
		}
	}
	return filter
}
