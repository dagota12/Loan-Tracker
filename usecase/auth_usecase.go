package usecase

import (
	"context"

	"github.com/dagota12/Loan-Tracker/domain"
	"github.com/dagota12/Loan-Tracker/internal/tokenutil"
)

type authUsecase struct {
	userRepo domain.UserRepository
}

// GetUserByID implements domain.AuthUsecase.
func (au *authUsecase) GetUserByID(ctx context.Context, userID string) (domain.User, error) {
	return au.userRepo.GetByID(ctx, userID)
}

// CreateAccessToken implements domain.AuthUsecase.
func (au *authUsecase) CreateAccessToken(user domain.User, secret string, expiry int) (accessToken string, err error) {
	return tokenutil.CreateAccessToken(user, secret, expiry)
}

// CreateRefreshToken implements domain.AuthUsecase.
func (au *authUsecase) CreateRefreshToken(user domain.User, secret string, expiry int) (refreshToken string, err error) {
	return tokenutil.CreateRefreshToken(&user, secret, expiry)
}

// GetUserByEmail implements domain.AuthUsecase.
func (au *authUsecase) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	return au.userRepo.GetByEmail(ctx, email)
}

// UpdateRefreshToken implements domain.AuthUsecase.
func (au *authUsecase) UpdateRefreshToken(ctx context.Context, userID string, refreshToken string) error {
	return au.userRepo.UpdateRefreshToken(ctx, userID, refreshToken)
}

func NewAuthUsease(userRepo domain.UserRepository) domain.AuthUsecase {
	return &authUsecase{
		userRepo: userRepo,
	}
}
