package repository

import (
	"context"
	"time"

	"nourish-backend/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UndoRepository interface {
	Create(ctx context.Context, op *models.UndoOperation) error
	GetByToken(ctx context.Context, token string) (*models.UndoOperation, error)
	MarkUndone(ctx context.Context, token string) error
	CleanupExpired(ctx context.Context) error
}

type undoRepository struct {
	collection *mongo.Collection
}

func NewUndoRepository(db *mongo.Database) UndoRepository {
	col := db.Collection("undo_operations")

	// TTL index on expiresAt
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "expiresAt", Value: 1}},
		Options: nil,
	})

	return &undoRepository{collection: col}
}

func (r *undoRepository) Create(ctx context.Context, op *models.UndoOperation) error {
	op.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, op)
	return err
}

func (r *undoRepository) GetByToken(ctx context.Context, token string) (*models.UndoOperation, error) {
	var op models.UndoOperation
	err := r.collection.FindOne(ctx, bson.M{"token": token}).Decode(&op)
	if err != nil {
		return nil, err
	}
	return &op, nil
}

func (r *undoRepository) MarkUndone(ctx context.Context, token string) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"token": token}, bson.M{"$set": bson.M{"undone": true}})
	return err
}

func (r *undoRepository) CleanupExpired(ctx context.Context) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{"expiresAt": bson.M{"$lt": time.Now()}})
	return err
}
