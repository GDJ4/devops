package todo

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrNotFound = errors.New("todo not found")

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(client *mongo.Client, dbName string) *Repository {
	return &Repository{
		collection: client.Database(dbName).Collection("todos"),
	}
}

func (r *Repository) List(ctx context.Context) ([]Todo, error) {
	cursor, err := r.collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"createdAt": -1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var todos []Todo
	for cursor.Next(ctx) {
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return todos, nil
}

func (r *Repository) Get(ctx context.Context, id primitive.ObjectID) (Todo, error) {
	var todo Todo
	if err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&todo); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Todo{}, ErrNotFound
		}
		return Todo{}, err
	}
	return todo, nil
}

func (r *Repository) Create(ctx context.Context, req CreateTodoRequest) (Todo, error) {
	now := time.Now().UTC()
	todo := Todo{
		Title:     req.Title,
		Notes:     req.Notes,
		Completed: false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result, err := r.collection.InsertOne(ctx, todo)
	if err != nil {
		return Todo{}, err
	}
	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return Todo{}, errors.New("failed to read inserted id")
	}
	todo.ID = id
	return todo, nil
}

func (r *Repository) Update(ctx context.Context, id primitive.ObjectID, req UpdateTodoRequest) (Todo, error) {
	updates := bson.M{}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Notes != nil {
		updates["notes"] = *req.Notes
	}
	if req.Completed != nil {
		updates["completed"] = *req.Completed
	}
	if len(updates) == 0 {
		return r.Get(ctx, id)
	}
	updates["updatedAt"] = time.Now().UTC()

	var todo Todo
	err := r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": updates},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&todo)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Todo{}, ErrNotFound
		}
		return Todo{}, err
	}
	return todo, nil
}

func (r *Repository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return ErrNotFound
	}
	return nil
}
