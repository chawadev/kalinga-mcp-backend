package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateIndexes(db *mongo.Database) error {
	ctx := context.Background()

	// Users collection indexes
	usersCollection := db.Collection("users")
	
	// Create unique index on email
	emailIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	
	_, err := usersCollection.Indexes().CreateOne(ctx, emailIndex)
	if err != nil {
		log.Printf("Warning: Error creating email index: %v", err)
		// Don't return error, just log it
	} else {
		log.Println("Database indexes created successfully")
	}

	return nil
}
