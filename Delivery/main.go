package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"task_manager1/Delivery/controllers"
	"task_manager1/Delivery/routers"
	"task_manager1/Infrastructure/auth"
	"task_manager1/Infrastructure/security"
	"task_manager1/Repositories"
	"task_manager1/Usecases"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoClient(ctx context.Context, uri string) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(uri)
	opts.SetServerSelectionTimeout(10 * time.Second)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(ctx)
		return nil, err
	}
	return client, nil
}

func main() {
	_ = godotenv.Load()

	uri := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("MONGODB_DATABASE")
	taskColl := os.Getenv("TASKS_COLLECTION")
	userColl := os.Getenv("USERS_COLLECTION")
	jwtSecret := os.Getenv("JWT_SECRET")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if uri == "" || dbName == "" || taskColl == "" || userColl == "" || jwtSecret == "" {
		log.Fatal("MONGODB_URI, MONGODB_DATABASE, TASKS_COLLECTION, USERS_COLLECTION and JWT_SECRET must be set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	client, err := NewMongoClient(ctx, uri)
	if err != nil {
		log.Fatalf("mongo connect error: %v", err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database(dbName)
	taskCollection := db.Collection(taskColl)
	userCollection := db.Collection(userColl)

	// Wire Repositories
	userRepo := Repositories.NewMongoUserRepository(userCollection)
	taskRepo := Repositories.NewMongoTaskRepository(taskCollection)

	// Infrastructure services
	pwSvc := security.NewPasswordService()
	jwtSvc := auth.NewJWTService()
	authMw := auth.NewAuthMiddleware()

	// Usecases
	userUC := Usecases.NewUserUsecase(userRepo, pwSvc)
	taskUC := Usecases.NewTaskUsecase(taskRepo)

	// controller
	ctl := controllers.NewController(userUC, taskUC, jwtSvc)

	// router
	r := routers.SetupRouter(ctl, authMw)

	addr := fmt.Sprintf(":%s", port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server exit: %v", err)
	}
}
