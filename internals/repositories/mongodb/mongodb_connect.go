package mongodb

import (
	"context"
	"log"
	"os"

	"github.com/Yash-Kansagara/GoGRPC_API/pkg/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var MongoClient *mongo.Client
var TeachersCollection *mongo.Collection

func InitializeMongoDBClient() (*mongo.Client, error) {
	conenctionString := os.Getenv("MONGO_CONN_STR")
	if len(conenctionString) == 0 {
		return nil, utils.ErrorHandler(nil, "Failed to read mongodb connection string environment variable")
	}

	mongoClient, err := mongo.Connect(options.Client().ApplyURI(conenctionString).SetBSONOptions(&options.BSONOptions{ObjectIDAsHexString: true}))
	if err != nil {
		return nil, utils.ErrorHandler(nil, "Error connecting to mongodb")
	}

	err = mongoClient.Ping(context.Background(), nil)
	if err != nil {
		return nil, utils.ErrorHandler(nil, "Error pinging the mongodb server post connection")
	}

	log.Println("Connected to MongoDB!")

	TeachersCollection = mongoClient.Database("SchoolDB").Collection("teachers")
	return mongoClient, nil
}
