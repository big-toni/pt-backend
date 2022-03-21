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

// Token model
type Token struct {
	Model     `bson:",inline"`
	AppID     string             `json:"appId" bson:"app_id"`
	TokenHash string             `json:"tokenHash" bson:"token_hash"`
	Source    string             `json:"source" bson:"source"`
	UserID    primitive.ObjectID `json:"userId,omitempty" bson:"user_id,omitempty"`
}

// User model
type User struct {
	Model        `bson:",inline"`
	Email        string `json:"email" bson:"email"`
	PasswordHash string `json:"-" bson:"password_hash"`
	// PrivateParcels []PrivateParcel `json:"privateParcels" bson:"private_parcels"`
}

// Parcel model
type Parcel struct {
	Model          `bson:",inline"`
	Description    string             `json:"description" bson:"description"`
	Name           string             `json:"name" bson:"name"`
	Timeline       *[]timelineEntry   `json:"timeline"`
	TrackingNumber string             `json:"trackingNumber" bson:"tracking_number"`
	UserID         primitive.ObjectID `json:"userId,omitempty" bson:"user_id,omitempty"`
}

type timelineEntry struct {
	Description string    `json:"description"`
	ID          string    `json:"id"`
	Location    *address  `json:"location"`
	Status      string    `json:"status"`
	Time        time.Time `json:"time"`
}

type address struct {
	City    string `json:"city"`
	Country string `json:"country"`
	Zip     string `json:"zip"`
}
