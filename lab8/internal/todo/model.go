package todo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Todo represents a single task item stored in MongoDB.
type Todo struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title" json:"title"`
	Completed bool               `bson:"completed" json:"completed"`
	Notes     string             `bson:"notes,omitempty" json:"notes,omitempty"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// CreateTodoRequest defines the payload for creating a new todo item.
type CreateTodoRequest struct {
	Title string `json:"title"`
	Notes string `json:"notes"`
}

// UpdateTodoRequest defines the payload for updating an existing todo item.
type UpdateTodoRequest struct {
	Title     *string `json:"title"`
	Notes     *string `json:"notes"`
	Completed *bool   `json:"completed"`
}
