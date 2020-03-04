package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Auth type
type Auth struct {
	Email                    string    `bson:"email" json:"email"`
	FirstName                string    `bson:"first_name" json:"first_name"`
	LastName                 string    `bson:"last_name" json:"last_name"`
	Name                     string    `bson:"name" json:"name"`
	Picture                  string    `bson:"picture" json:"picture"`
	Password                 string    `bson:"password" json:"password"`
	PasswordHashMethod       string    `bson:"password_hash_method" json:"password_hash_method"`
	PasswordResetID          string    `bson:"password_reset_id" json:"password_reset_id"`
	PasswordResetToken       string    `bson:"password_reset_token" json:"password_reset_token"`
	PasswordResetExpires     time.Time `bson:"password_reset_expires" json:"password_reset_expires"`
	PasswordModified         time.Time `bson:"password_modified" json:"password_modified"`
	Verified                 bool      `bson:"verified" json:"verified"`
	VerifiedDate             time.Time `bson:"verified_date" json:"verified_date"`
	VerificationToken        string    `bson:"verification_token" json:"verification_token"`
	VerificationTokenExpires time.Time `bson:"verification_token_expires" json:"verification_token_expires"`
}

type Session struct {
	IP        string    `bson:"ip" json:"ip"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
	Location  string    `bson:"location" json:"location"`
	Lat       int32     `bson:"lat" json:"lat"`
	Long      int32     `bson:"long" json:"long"`
}

// User type
type User struct {
	ID        primitive.ObjectID `bson:"id,omitempty" json:"id"`
	Object    string             `bson:"object" json:"object" validate:" eq=user,required" `
	Auth      Auth               `bson:"auth" json:"auth"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	LoggedIn  time.Time          `bson:"logged_in" json:"logged_in"`
	Sessions  Session            `bson:"sessions" json:"sessions"`
}
