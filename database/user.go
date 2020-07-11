package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"pt-server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//data access object

// UserDAO persists user data in database
type UserDAO struct{}

// NewUserDAO creates a new UserDAO
func NewUserDAO() *UserDAO {
	return &UserDAO{}
}

// Save func
func (dao *UserDAO) Save(user models.User) primitive.ObjectID {
	collection := Database.Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := collection.InsertOne(ctx, user)
	if err != nil {
		log.Fatal(err)
	}

	id := res.InsertedID.(primitive.ObjectID)
	fmt.Println("Inserted a single document: ", id)

	return id
}

// GetUserForEmail func
func (dao *UserDAO) GetUserForEmail(email string) *models.User {
	var user models.User
	filter := bson.M{"email": email}

	collection := Database.Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil
	}

	fmt.Printf("Found a single document: %+v\n", user)

	return &user
}

// GetUserByID func
func (dao *UserDAO) GetUserByID(ID primitive.ObjectID) models.User {
	var user models.User

	// userID, err := primitive.ObjectIDFromHex(ID)
	// if err != nil {
	// 	log.Println("Invalid ObjectID")
	// }

	filter := bson.M{"_id": ID}

	collection := Database.Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found a single document: %+v\n", user)

	return user
}
