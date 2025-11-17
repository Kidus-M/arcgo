package Repositories

import (
	"context"
	"errors"
	"time"

	"task_manager1/Domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskRepository interface {
	FindAll(ctx context.Context) ([]Domain.Task, error)
	FindByID(ctx context.Context, hexID string) (Domain.Task, error)
	Create(ctx context.Context, t Domain.Task) (Domain.Task, error)
	Update(ctx context.Context, hexID string, t Domain.Task) (Domain.Task, error)
	Delete(ctx context.Context, hexID string) (bool, error)
}

type mongoTaskRepository struct {
	coll    *mongo.Collection
	timeout time.Duration
}

func NewMongoTaskRepository(coll *mongo.Collection) TaskRepository {
	return &mongoTaskRepository{coll: coll, timeout: 5 * time.Second}
}

func (r *mongoTaskRepository) FindAll(ctx context.Context) ([]Domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	cur, err := r.coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var tasks []Domain.Task
	for cur.Next(ctx) {
		var t Domain.Task
		if err := cur.Decode(&t); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, cur.Err()
}

func (r *mongoTaskRepository) FindByID(ctx context.Context, hexID string) (Domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	oid, err := primitive.ObjectIDFromHex(hexID)
	if err != nil {
		return Domain.Task{}, errors.New("invalid id")
	}
	var t Domain.Task
	if err := r.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&t); err != nil {
		if err == mongo.ErrNoDocuments {
			return Domain.Task{}, nil
		}
		return Domain.Task{}, err
	}
	return t, nil
}

func (r *mongoTaskRepository) Create(ctx context.Context, t Domain.Task) (Domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	res, err := r.coll.InsertOne(ctx, t)
	if err != nil {
		return Domain.Task{}, err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		t.ID = oid
	}
	return t, nil
}

func (r *mongoTaskRepository) Update(ctx context.Context, hexID string, updated Domain.Task) (Domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	oid, err := primitive.ObjectIDFromHex(hexID)
	if err != nil {
		return Domain.Task{}, errors.New("invalid id")
	}
	updateDoc := bson.M{}
	if updated.Title != "" {
		updateDoc["title"] = updated.Title
	}
	if updated.Description != "" {
		updateDoc["description"] = updated.Description
	}
	if updated.DueDate != "" {
		updateDoc["due_date"] = updated.DueDate
	}
	if updated.Status != "" {
		updateDoc["status"] = updated.Status
	}
	if len(updateDoc) == 0 {
		return Domain.Task{}, errors.New("no fields to update")
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var result Domain.Task
	if err := r.coll.FindOneAndUpdate(ctx, bson.M{"_id": oid}, bson.M{"$set": updateDoc}, opts).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return Domain.Task{}, nil
		}
		return Domain.Task{}, err
	}
	return result, nil
}

func (r *mongoTaskRepository) Delete(ctx context.Context, hexID string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	oid, err := primitive.ObjectIDFromHex(hexID)
	if err != nil {
		return false, errors.New("invalid id")
	}
	res, err := r.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return false, err
	}
	return res.DeletedCount > 0, nil
}
