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

// UserRepository defines user data methods
type UserRepository interface {
	Create(ctx context.Context, u Domain.User) (Domain.User, error)
	FindByUsername(ctx context.Context, username string) (Domain.User, error)
	PromoteToAdmin(ctx context.Context, username string) (Domain.User, error)
	Count(ctx context.Context) (int64, error)
}

type mongoUserRepository struct {
	coll    *mongo.Collection
	timeout time.Duration
}

func NewMongoUserRepository(coll *mongo.Collection) UserRepository {
	// ensure unique username index
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, _ = coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	return &mongoUserRepository{coll: coll, timeout: 5 * time.Second}
}

func (r *mongoUserRepository) Create(ctx context.Context, u Domain.User) (Domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	res, err := r.coll.InsertOne(ctx, u)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return Domain.User{}, errors.New("username already exists")
		}
		return Domain.User{}, err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		u.ID = oid
	}
	u.PasswordHash = "" // never return hash
	return u, nil
}

func (r *mongoUserRepository) FindByUsername(ctx context.Context, username string) (Domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	var u Domain.User
	err := r.coll.FindOne(ctx, bson.M{"username": username}).Decode(&u)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Domain.User{}, nil
		}
		return Domain.User{}, err
	}
	return u, nil
}

func (r *mongoUserRepository) PromoteToAdmin(ctx context.Context, username string) (Domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	update := bson.M{"$set": bson.M{"role": "admin"}}
	var updated Domain.User
	if err := r.coll.FindOneAndUpdate(ctx, bson.M{"username": username}, update, opts).Decode(&updated); err != nil {
		if err == mongo.ErrNoDocuments {
			return Domain.User{}, nil
		}
		return Domain.User{}, err
	}
	updated.PasswordHash = ""
	return updated, nil
}

func (r *mongoUserRepository) Count(ctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	return r.coll.CountDocuments(ctx, bson.M{})
}
