package domain

import (
	"context"
	"time"
)

// user reset and forgot password
type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required"`
	Code        string `json:"code" binding:"required"`
	NewPassword string `json:"password" bson:"password" binding:"required,min=4,max=30"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required"`
}

type UpdatePassword struct {
	OldPassword string `json:"old_password" bson:"old_password" binding:"required,min=4,max=30"`
	NewPassword string `json:"new_password" bson:"password" binding:"required,min=4,max=30"`
}
type OtpSave struct {
	Email     string    `json:"email" binding:"required"`
	Code      string    `json:"code" binding:"required"`
	ExpiresAt time.Time `json:"expiresat" sql:"expiresat"`
}
type ResetPasswordUsecase interface {
	GetUserByEmail(ctx context.Context, email string) (User, error)
	ResetPassword(ctx context.Context, userID string, resetPassword *ResetPasswordRequest) error
	SaveOtp(ctx context.Context, otp *OtpSave) error
	DeleteOtp(ctx context.Context, email string) error
	GetOTPByEmail(ctx context.Context, email string) (*OtpSave, error)
}

type ResetPasswordRepository interface {
	GetUserByEmail(context.Context, string) (*User, error)
	ResetPassword(c context.Context, userID string, resetPassword *ResetPasswordRequest) error
	SaveOtp(c context.Context, otp *OtpSave) error
	GetOTPByEmail(c context.Context, email string) (*OtpSave, error)
	DeleteOtp(c context.Context, email string) error
}

const (
	CollectionResetPassword = "resetPassword"
)
