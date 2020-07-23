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
	filter := bson.M{
		"user_id": userID,
		"deleted_at": bson.M{
			"$exists": false,
		},
	}

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
func (dao *ParcelDAO) Save(parcels []models.Parcel) []primitive.ObjectID {
	collection := Database.Collection("parcels")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	y := make([]interface{}, len(parcels))
	for i, v := range parcels {
		y[i] = v
	}

	res, err := collection.InsertMany(ctx, y)
	if err != nil {
		log.Fatal(err)
	}

	var ids []primitive.ObjectID

	for _, id := range res.InsertedIDs {
		ids = append(ids, id.(primitive.ObjectID))
	}

	fmt.Println("Inserted documents: ", ids)

	return ids
}

// Update func
func (dao *ParcelDAO) Update(parcel models.Parcel) primitive.ObjectID {
	collection := Database.Collection("parcels")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	filter := bson.M{"_id": parcel.Model.ID}

	update := bson.M{
		"$set": bson.M{
			"description": parcel.Description,
			"name":        parcel.Name,
			"timeline":    parcel.Timeline,
			"updated_at":  time.Now(),
		},
	}

	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
	}

	if res.MatchedCount == 0 {
		return primitive.NilObjectID
	}

	fmt.Println("Updated document: ", parcel.Model.ID)

	return parcel.Model.ID
}

// Delete func
func (dao *ParcelDAO) Delete(parcel models.Parcel) primitive.ObjectID {
	collection := Database.Collection("parcels")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	filter := bson.M{"_id": parcel.Model.ID}

	update := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
			"deleted_at": time.Now(),
		},
	}

	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
	}

	if res.MatchedCount == 0 {
		return primitive.NilObjectID
	}

	fmt.Println("Deleted document: ", parcel.Model.ID)

	return parcel.Model.ID
}
