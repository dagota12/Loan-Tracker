package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `json:"_id"  bson:"_id,omitempty"`
	FirstName   string             `json:"first_name" bson:"first_name" binding:"required,min=3,max=30"`
	LastName    string             `json:"last_name" bson:"last_name" binding:"max=30"`
	Email       string             `json:"email" bson:"email" binding:"required,email"`
	Active      bool               `json:"active" bson:"active"`
	Password    string             `json:"-" bson:"password"`
	VerifyToken string             `json:"-" bson:"verify_token"`
	IsOwner     bool               `json:"is_owner" bson:"is_owner"`
	Tokens      []string           `json:"-" bson:"refresh_tokens"`
	Role        string             `joson:"role" bson:"role"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
	LastLogin   time.Time          `json:"last_login" bson:"last_login"`
}

type UserUpdate struct {
	FirstName string    `json:"first_name" bson:"first_name"`
	LastName  string    `json:"last_name" bson:"last_name"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
type UserForm struct {
	FirstName string    `json:"first_name" bson:"first_name" binding:"required,min=3,max=30"`
	LastName  string    `json:"last_name" bson:"last_name" binding:"max=30"`
	Email     string    `json:"email" bson:"email" binding:"required,email"`
	Password  string    `json:"password" bson:"password" binding:"required"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

// user repository
type UserRepository interface {
	GetAll(ctx context.Context) ([]User, error)
	GetByID(ctx context.Context, userID string) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	Create(ctx context.Context, user User) (User, error)
	Update(ctx context.Context, userID string, user UserUpdate) (User, error)
	Delete(ctx context.Context, userID string) error

	IsOwner(ctx context.Context, userID string) (bool, error)

	RevokeRefreshToken(ctx context.Context, userID, refreshToken string) error
	UpdateRefreshToken(ctx context.Context, userID string, refreshToken string) error
	RefreshTokenExist(ctx context.Context, userID, refreshToken string) (bool, error)

	IsUserActive(ctx context.Context, userID string) (bool, error)
	ActivateUser(ctx context.Context, userID string) error

	ResetUserPassword(ctx context.Context, userID string, resetPassword ResetPasswordRequest) error
	UpdateUserPassword(ctx context.Context, userID string, updatePassword UpdatePassword) error
}

type UserUsecase interface {
	GetAll(ctx context.Context) ([]User, error)
	GetByID(ctx context.Context, userID string) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	Create(ctx context.Context, user User) (User, error)
	Update(ctx context.Context, userID string, user UserUpdate) (User, error)
	Delete(ctx context.Context, userID string) error
	ResetUserPassword(ctx context.Context, userID string, resetPassword ResetPasswordRequest) error
	UpdateUserPassword(ctx context.Context, userID string, updatePassword UpdatePassword) error
}
