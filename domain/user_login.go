package domain

import (
	"context"
)

//user login

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type AuthUsecase interface {
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByID(ctx context.Context, userID string) (User, error)
	CreateAccessToken(user User, secret string, expiry int) (accessToken string, err error)
	CreateRefreshToken(user User, secret string, expiry int) (refreshToken string, err error)
	UpdateRefreshToken(ctx context.Context, userID string, refreshToken string) error
}
