package domain

import (
	"context"

	"github.com/dagota12/Loan-Tracker/bootstrap"
)

type SignupRequest struct {
	FirstName string `json:"first_name" bson:"first_name" binding:"required,min=3,max=30"`
	LastName  string `json:"last_name" bson:"last_name" binding:"required,min=3,max=30"`
	Email     string `json:"email" bson:"email" binding:"required,email"`
	Password  string `json:"password" bson:"password" binding:"required,min=4,max=30"`
}
type SignupResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SignupUsecase interface {
	Create(ctx context.Context, user *User) (User, error)
	ActivateUser(c context.Context, userID string) error
	IsOwner(ctx context.Context, userID string) (bool, error)
	GetUserById(c context.Context, userId string) (*User, error)
	GetUserByEmail(c context.Context, email string) (User, error)
	CreateVerificationToken(user *User, secret string, expiry int) (accessToken string, err error)
	CreateAccessToken(user *User, secret string, expiry int) (accessToken string, err error)
	CreateRefreshToken(user *User, secret string, expiry int) (refreshToken string, err error)
	SendVerificationEmail(recipientEmail string, encodedToken string, env *bootstrap.Env) (err error)
}
