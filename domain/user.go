package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	FirstName   string             `json:"first_name" bson:"first_name"`
	LastName    string             `json:"last_name" bson:"last_name"`
	Emali       string             `json:"email" bson:"email"`
	Active      bool               `json:"active" bson:"active"`
	Password    string             `json:"-" bson:"password"`
	VerifyToken string             `json:"-" bson:"verify_token"`
	IsOwner     bool               `json:"is_owner" bson:"is_owner"`
	Tokens      []string           `json:"-" bson:"refresh_tokens"`
	Role        string             `joson:"role" bson:"role"`
	CreatedAt   primitive.DateTime `json:"created_at" bson:"created_at"`
	UpdatedAt   primitive.DateTime `json:"updated_ats" bson:"updated_at"`
	LastLogin   primitive.DateTime `json:"last_logins" bson:"last_login"`
}

type UserUpdate struct {
	FirstName string             `json:"first_name" bson:"first_name"`
	LastName  string             `json:"last_name" bson:"last_name"`
	UpdatedAt primitive.DateTime `json:"updated_ats" bson:"updated_at"`
}

// user repository
type UserRepository interface {
	GetAll(ctx context.Context) ([]User, error)
	GetByID(ctx context.Context, userID string) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	Create(ctx context.Context, user User) (User, error)
	Update(ctx context.Context, userID string, user UserUpdate) (User, error)
	Delete(ctx context.Context, userID string) error

	RevokeRefreshToken(ctx context.Context, userID, refreshToken string) error
	UpdateRefreshToken(ctx context.Context, userID string, refreshToken string) error
	RefreshTokenExist(ctx context.Context, userID, refreshToken string) (bool, error)

	IsUserActive(ctx context.Context, userID string) (bool, error)
	ActivateUser(ctx context.Context, userID string) error

	ResetUserPassword(ctx context.Context, userID string, resetPassword ResetPasswordRequest) error
	UpdateUserPassword(ctx context.Context, userID string, updatePassword UpdatePassword) error
}
