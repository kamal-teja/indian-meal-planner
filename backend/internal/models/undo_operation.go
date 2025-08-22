package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UndoOperation represents a soft-delete operation that can be undone with a token
type UndoOperation struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Token     string               `bson:"token" json:"token"`
	UserID    primitive.ObjectID   `bson:"userId" json:"userId"`
	MealIDs   []primitive.ObjectID `bson:"mealIds" json:"mealIds"`
	CreatedAt time.Time            `bson:"createdAt" json:"createdAt"`
	ExpiresAt time.Time            `bson:"expiresAt" json:"expiresAt"`
	Undone    bool                 `bson:"undone" json:"undone"`
}
