package repositories

import (
	"context"
	"errors"
	"time"

	"task_manager/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserRepository defines user data methods
type UserRepository interface {
	Create(ctx context.Context, u domain.User) (domain.User, error)
	FindByUsername(ctx context.Context, username string) (domain.User, error)
	PromoteToAdmin(ctx context.Context, username string) (domain.User, error)
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

func (r *mongoUserRepository) Create(ctx context.Context, u domain.User) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	res, err := r.coll.InsertOne(ctx, u)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.User{}, errors.New("username already exists")
		}
		return domain.User{}, err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		u.ID = oid
	}
	u.PasswordHash = "" // never return hash
	return u, nil
}

func (r *mongoUserRepository) FindByUsername(ctx context.Context, username string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	var u domain.User
	err := r.coll.FindOne(ctx, bson.M{"username": username}).Decode(&u)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.User{}, nil
		}
		return domain.User{}, err
	}
	return u, nil
}

func (r *mongoUserRepository) PromoteToAdmin(ctx context.Context, username string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	update := bson.M{"$set": bson.M{"role": "admin"}}
	var updated domain.User
	if err := r.coll.FindOneAndUpdate(ctx, bson.M{"username": username}, update, opts).Decode(&updated); err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.User{}, nil
		}
		return domain.User{}, err
	}
	updated.PasswordHash = ""
	return updated, nil
}

func (r *mongoUserRepository) Count(ctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	return r.coll.CountDocuments(ctx, bson.M{})
}
