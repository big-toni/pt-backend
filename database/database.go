package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database var
var Database *mongo.Database

// Connect func
func Connect() {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Set client options
	// clientOptions := options.Client().ApplyURI("mongodb://" + user + ":" + password + "@ds155663.mlab.com:55663/" + dbName + "?retryWrites=false")
	clientOptions := options.Client().ApplyURI("mongodb+srv://" + user + ":" + password + "@freecluster.vmlxh.mongodb.net/" + dbName + "?retryWrites=true&w=majority")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	Database = client.Database(dbName)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

}
