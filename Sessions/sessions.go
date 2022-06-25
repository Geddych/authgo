package sessions

import (
	"authgo/tokens"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.TODO()
var sColl *mongo.Collection
var client *mongo.Client

type Session struct {
	id       string
	token    string
	IssuedAt time.Time
	Active   int
	ClosedAt time.Time
}

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
	sColl = client.Database("auth").Collection("Session")

}

func CreateSession(token string, sId interface{}) {
	initBase()
	hashingToken := tokens.HashToken(token)
	sColl.InsertOne(ctx, bson.D{{"_id", sId}, {"token", hashingToken}, {"IssuedAt", time.Now()}, {"Active", 1}, {"ClosedAt", ""}})

}

func CloseSession(sId interface{}) {
	initBase()
	filter := bson.D{{"_id", sId}}
	update := bson.M{
		"$set": bson.M{"Active": 0, "ClosedAt": time.Now()},
	}
	sColl.FindOneAndUpdate(ctx, filter, update)

}

func GetTokenFromSessionByID(sId interface{}) interface{} {
	initBase()
	var result bson.D
	filter := bson.D{{"_id", sId}}
	err := sColl.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Print("No Documents finded!")
	}

	id := result.Map()["_id"]
	return id
}

func UpdateTokenSession(sId interface{}, newToken string) {
	initBase()
	newHashed := tokens.HashToken(newToken)
	filter := bson.D{{"_id", sId}}
	update := bson.M{
		"$set": bson.M{"token": newHashed},
	}
	sColl.FindOneAndUpdate(ctx, filter, update)
}

func CheckTokeninSession(tk string) bool {
	initBase()
	var match bool
	tid := tokens.GetIdFromToken(tk)
	sid := GetTokenFromSessionByID(tid)
	if sid != nil {
		match = true
	} else {
		match = false
	}
	fmt.Print(match)
	return match
}
