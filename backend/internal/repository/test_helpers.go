package repository

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

// TestDB provides a mocked MongoDB database for testing
type TestDB struct {
	mt       *mtest.T
	database *mongo.Database
}

// NewTestDB creates a new test database instance
func NewTestDB(t *testing.T) *TestDB {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	return &TestDB{
		mt:       mt,
		database: mt.DB,
	}
}

// Run executes a test function with the mocked database
func (tdb *TestDB) Run(name string, testFunc func(*mtest.T)) {
	tdb.mt.Run(name, testFunc)
}

// GetDatabase returns the mocked database instance
func (tdb *TestDB) GetDatabase() *mongo.Database {
	return tdb.database
}

// Close cleans up the test database
func (tdb *TestDB) Close() {
	// Cleanup is handled by mtest
}

// Helper function to create a context for testing
func testContext() context.Context {
	return context.Background()
}