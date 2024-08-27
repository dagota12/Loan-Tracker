package usecase

import (
	"context"
	"time"

	"github.com/dagota12/Loan-Tracker/bootstrap"
	"github.com/dagota12/Loan-Tracker/domain"
	"github.com/dagota12/Loan-Tracker/internal/emailutil"
	"github.com/dagota12/Loan-Tracker/internal/tokenutil"
)

type signupUsecase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

func NewSignupUsecase(userRepository domain.UserRepository, timeout time.Duration) domain.SignupUsecase {
	return &signupUsecase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (su *signupUsecase) GetUserById(ctx context.Context, userId string) (*domain.User, error) {
	user, err := su.userRepository.GetByID(ctx, userId)
	return &user, err
}

func (su *signupUsecase) ActivateUser(ctx context.Context, userID string) error {
	err := su.userRepository.ActivateUser(ctx, userID)
	return err
}

func (su *signupUsecase) Create(ctx context.Context, user *domain.User) (domain.User, error) {
	return su.userRepository.Create(ctx, *user)
}

func (su *signupUsecase) IsOwner(ctx context.Context, userID string) (bool, error) {
	result, err := su.userRepository.IsOwner(ctx, userID)
	return result, err
}
func (su *signupUsecase) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := su.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (su *signupUsecase) CreateAccessToken(user *domain.User, secret string, expiry int) (accessToken string, err error) {
	return tokenutil.CreateAccessToken(*user, secret, expiry)
}

func (su *signupUsecase) CreateRefreshToken(user *domain.User, secret string, expiry int) (refreshToken string, err error) {
	return tokenutil.CreateRefreshToken(user, secret, expiry)
}

func (su *signupUsecase) CreateVerificationToken(user *domain.User, secret string, expiry int) (refreshToken string, err error) {
	return tokenutil.CreateVerificationToken(user, secret, expiry)
}

func (su *signupUsecase) SendVerificationEmail(recipientEmail string, encodedToken string, env *bootstrap.Env) (err error) {
	return emailutil.SendVerificationEmail(recipientEmail, encodedToken, env)
}

// CanBeOwner implements domain.SignupUsecase.
func (su *signupUsecase) CanBeOwner(ctx context.Context) (bool, error) {
	users, _ := su.userRepository.GetAll(ctx)
	return len(users) == 0, nil
}
