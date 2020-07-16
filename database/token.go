package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"pt-server/database/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TokenDAO persists token data in database
type TokenDAO struct{}

// NewTokenDAO creates a new TokenDAO
func NewTokenDAO() *TokenDAO {
	return &TokenDAO{}
}

// Save func
func (dao *TokenDAO) Save(token models.Token) string {
	collection := Database.Collection("tokens")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// res, err := collection.InsertOne(ctx, bson.M{"name": "pi", "value": 3.14159})
	res, err := collection.InsertOne(ctx, token)

	if err != nil {
		log.Fatal(err)
	}
	id := fmt.Sprintf("%v", res.InsertedID)
	fmt.Println("Inserted a single document: ", id)

	return id
}

// GetUserID func
func (dao *TokenDAO) GetUserID(tokenHash string) *primitive.ObjectID {
	var document models.Token
	collection := Database.Collection("tokens")
	filter := bson.M{"token_hash": tokenHash}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, filter).Decode(&document)
	if err != nil {
		log.Println(err)
		return nil
	}

	return &document.UserID
}
