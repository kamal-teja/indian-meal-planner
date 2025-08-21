package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, using environment variables")
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI environment variable is required")
	}

	fmt.Println("🔌 Testing MongoDB Atlas connection...")
	fmt.Printf("📋 URI: %s\n", hidePassword(mongoURI))

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Test the connection
	fmt.Println("🏓 Pinging MongoDB...")
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("❌ Failed to ping MongoDB: %v", err)
	}

	fmt.Println("✅ MongoDB connection successful!")

	// Test database access
	db := client.Database("meal-planner")
	fmt.Println("📊 Testing database access...")

	// Try to list collections
	collections, err := db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		log.Fatalf("❌ Failed to list collections: %v", err)
	}

	fmt.Printf("✅ Successfully accessed database 'meal-planner'\n")
	fmt.Printf("📁 Found %d collections: %v\n", len(collections), collections)

	// Test creating a test document
	testCollection := db.Collection("connection_test")
	testDoc := bson.M{
		"test":      true,
		"timestamp": time.Now(),
		"message":   "Connection test successful",
	}

	fmt.Println("📝 Testing write operation...")
	result, err := testCollection.InsertOne(ctx, testDoc)
	if err != nil {
		log.Printf("⚠️  Warning: Failed to insert test document: %v", err)
	} else {
		fmt.Printf("✅ Write test successful! Inserted document ID: %v\n", result.InsertedID)

		// Clean up test document
		_, err = testCollection.DeleteOne(ctx, bson.M{"_id": result.InsertedID})
		if err != nil {
			log.Printf("⚠️  Warning: Failed to clean up test document: %v", err)
		} else {
			fmt.Println("🧹 Cleaned up test document")
		}
	}

	fmt.Println("🎉 All database tests passed! Your MongoDB Atlas connection is working correctly.")
}

// hidePassword masks the password in the URI for safe logging
func hidePassword(uri string) string {
	// Simple password masking for display
	if len(uri) > 20 {
		return uri[:15] + "***[HIDDEN]***" + uri[len(uri)-15:]
	}
	return "***[HIDDEN]***"
}
