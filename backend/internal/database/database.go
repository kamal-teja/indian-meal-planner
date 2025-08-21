package database

import (
	"context"
	"strings"
	"time"

	"meal-planner-backend/internal/config"
	"meal-planner-backend/pkg/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database wraps MongoDB client and provides database operations
type Database struct {
	client *mongo.Client
	db     *mongo.Database
	logger *logger.Logger
}

// Connect establishes a connection to MongoDB
func Connect(uri string, cfg config.DatabaseConfig, log *logger.Logger) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set up MongoDB client options
	clientOptions := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(uint64(cfg.MaxPoolSize)).
		SetMinPoolSize(uint64(cfg.MinPoolSize)).
		SetMaxConnIdleTime(cfg.MaxIdleTime).
		SetServerSelectionTimeout(5 * time.Second).
		SetSocketTimeout(45 * time.Second).
		SetConnectTimeout(10 * time.Second)

	// For MongoDB Atlas, ensure we're using the correct auth source
	if strings.Contains(uri, "mongodb+srv://") || strings.Contains(uri, "mongodb.net") {
		if !strings.Contains(uri, "authSource=") {
			separator := "?"
			if strings.Contains(uri, "?") {
				separator = "&"
			}
			uri = uri + separator + "authSource=admin"
			clientOptions = options.Client().ApplyURI(uri).
				SetMaxPoolSize(uint64(cfg.MaxPoolSize)).
				SetMinPoolSize(uint64(cfg.MinPoolSize)).
				SetMaxConnIdleTime(cfg.MaxIdleTime).
				SetServerSelectionTimeout(5 * time.Second).
				SetSocketTimeout(45 * time.Second).
				SetConnectTimeout(10 * time.Second)
		}
	}

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	// Extract database name from URI or use default
	dbName := extractDatabaseName(uri)
	if dbName == "" {
		dbName = "meal-planner"
	}

	db := client.Database(dbName)

	log.Info("MongoDB connected successfully",
		"host", uri,
		"database", dbName,
		"maxPoolSize", cfg.MaxPoolSize,
		"minPoolSize", cfg.MinPoolSize)

	return &Database{
		client: client,
		db:     db,
		logger: log,
	}, nil
}

// GetDB returns the MongoDB database instance
func (d *Database) GetDB() *mongo.Database {
	return d.db
}

// GetClient returns the MongoDB client instance
func (d *Database) GetClient() *mongo.Client {
	return d.client
}

// Disconnect closes the database connection
func (d *Database) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := d.client.Disconnect(ctx); err != nil {
		d.logger.Error("Failed to disconnect from MongoDB", "error", err)
		return err
	}

	d.logger.Info("Disconnected from MongoDB")
	return nil
}

// Ping verifies the database connection
func (d *Database) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return d.client.Ping(ctx, nil)
}

// extractDatabaseName extracts database name from MongoDB URI
func extractDatabaseName(uri string) string {
	// Parse URI to extract database name
	// Format: mongodb://[username:password@]host1[:port1][,...hostN[:portN]][/[database][?options]]
	// Or: mongodb+srv://[username:password@]host[/[database][?options]]

	// Find the start of the database name after the last '/'
	lastSlash := -1
	for i := len(uri) - 1; i >= 0; i-- {
		if uri[i] == '/' {
			// Skip if this is part of the protocol (://)
			if i >= 2 && uri[i-2:i] == ":/" {
				continue
			}
			lastSlash = i
			break
		}
	}

	if lastSlash == -1 || lastSlash == len(uri)-1 {
		return ""
	}

	// Extract database name (everything after '/' until '?' or end)
	dbName := uri[lastSlash+1:]
	if questionMark := strings.Index(dbName, "?"); questionMark != -1 {
		dbName = dbName[:questionMark]
	}

	return dbName
}
