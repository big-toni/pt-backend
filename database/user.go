package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"pt-backend/database/models"

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
	fmt.Printf("Inserted new User with ID: %+v\n", id)

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

	fmt.Printf("Found User with ID: %+v\n", user.ID)

	return &user
}

// GetUserByID func
func (dao *UserDAO) GetUserByID(ID primitive.ObjectID) *models.User {
	var user models.User
	filter := bson.M{"_id": ID}
	collection := Database.Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil
	}

	fmt.Printf("Found User with ID: %+v\n", user.ID)

	return &user
}

// Update func
func (dao *UserDAO) Update(user models.User) primitive.ObjectID {
	collection := Database.Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	filter := bson.M{"_id": user.Model.ID}

	update := bson.M{
		"$set": bson.M{
			"email":         user.Email,
			"password_hash": user.PasswordHash,
			"updated_at":    time.Now(),
		},
	}

	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
	}

	if res.MatchedCount == 0 {
		return primitive.NilObjectID
	}

	fmt.Printf("Updated User with ID: %+v\n", user.Model.ID)

	return user.Model.ID
}
