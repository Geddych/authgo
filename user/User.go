package user

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.TODO()
var userColl *mongo.Collection
var client *mongo.Client

func initBase() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	userColl = client.Database("auth").Collection("User")

}

func FindUserIDbyUsername(name string) interface{} {
	initBase()
	var result bson.D
	filter := bson.D{{"username", name}}
	err := userColl.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return mongo.ErrNoDocuments
	}
	id := result.Map()["_id"]
	return id
}
