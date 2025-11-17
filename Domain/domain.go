package Domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// Task entity
type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	DueDate     string             `bson:"due_date,omitempty" json:"due_date,omitempty"`
	Status      string             `bson:"status,omitempty" json:"status,omitempty"`
}

// TaskResponse for API (ID as hex)
type TaskResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	DueDate     string `json:"due_date,omitempty"`
	Status      string `json:"status,omitempty"`
}

// User entity
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	Username     string             `bson:"username" json:"username"`
	PasswordHash string             `bson:"password_hash" json:"-"`
	Role         string             `bson:"role" json:"role"` // "admin" or "user"
}
