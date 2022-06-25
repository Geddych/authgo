package utils

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var ctx = context.TODO()

func CreateUser(uname string, uColl *mongo.Collection) {
	insResult, err := uColl.InsertOne(ctx, bson.D{{"_id", primitive.NewObjectID()}, {"username", uname}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(insResult)
}

func HashString(str string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(bytes), err
}
func CheckHash(str, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(str))
	return err == nil
}
