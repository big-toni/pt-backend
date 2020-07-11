package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Model definition same as gorm.Model, but including column and json tags
type Model struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updated_at"`
	DeletedAt *time.Time         `json:"deletedAt,omitempty" bson:"deleted_at,omitempty"`
}

// User model
type User struct {
	Model        `bson:",inline"`
	Email        string `bson:"email"`
	PasswordHash string `bson:"password_hash"`
}

// Token model
type Token struct {
	Model     `bson:",inline"`
	AppID     string             `json:"appId" bson:"app_id"`
	TokenHash string             `json:"tokenHash" bson:"token_hash"`
	Source    string             `json:"source" bson:"source"`
	UserID    primitive.ObjectID `json:"userId,omitempty" bson:"user_id,omitempty"`
}

// Parcel model
type Parcel struct {
	Model
	TrackingID string `json:"trackingId" bson:"tracking_id"`
}
