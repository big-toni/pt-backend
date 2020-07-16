package database

import (
	"context"
	"fmt"
	"log"
	"pt-server/database/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ParcelDAO persists parcel data in database
type ParcelDAO struct{}

// NewParcelDAO creates a new ParcelDAO
func NewParcelDAO() *ParcelDAO {
	return &ParcelDAO{}
}

// GetParcelsForUserID func
func (dao *ParcelDAO) GetParcelsForUserID(userID primitive.ObjectID) []*models.Parcel {
	var documents []*models.Parcel
	collection := Database.Collection("parcels")
	filter := bson.M{"user_id": userID}

	findOptions := options.Find()
	// Sort by `createdAt` field descending
	findOptions.SetSort(bson.D{primitive.E{Key: "created_at", Value: -1}})

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cur, err := collection.Find(ctx, filter, findOptions)

	if err != nil {
		log.Println(err)
		return nil
	}

	for cur.Next(context.TODO()) {
		var elem models.Parcel
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println("Parse error : ", err)
		}
		documents = append(documents, &elem)
	}

	return documents
}

// Save func
func (dao *ParcelDAO) Save(parcel models.Parcel) primitive.ObjectID {
	collection := Database.Collection("parcels")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := collection.InsertOne(ctx, parcel)
	if err != nil {
		log.Fatal(err)
	}

	id := res.InsertedID.(primitive.ObjectID)
	fmt.Println("Inserted a single document: ", id)

	return id
}
