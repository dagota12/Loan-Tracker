package usecase

import (
	"context"
	"time"

	"github.com/dagota12/Loan-Tracker/domain"
	"golang.org/x/crypto/bcrypt"
)

type resetPasswordUsecase struct {
	resetPasswordRepository domain.ResetPasswordRepository
	contextTimeout          time.Duration
}

func NewResetPasswordUsecase(resetPasswordRepository domain.ResetPasswordRepository, timeout time.Duration) domain.ResetPasswordUsecase {
	return &resetPasswordUsecase{
		resetPasswordRepository: resetPasswordRepository,
		contextTimeout:          timeout,
	}
}
func (r *resetPasswordUsecase) SaveOtp(c context.Context, otp *domain.OtpSave) error {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()
	err := r.resetPasswordRepository.SaveOtp(ctx, otp)
	return err
}

func (r *resetPasswordUsecase) GetUserByEmail(c context.Context, email string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()
	user, err := r.resetPasswordRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return *user, nil
}
func (r *resetPasswordUsecase) ResetPassword(c context.Context, userID string, resetPassword *domain.ResetPasswordRequest) error {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()
	//enctypt the user password

	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(resetPassword.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	resetPassword.NewPassword = string(bcryptPassword)

	err = r.resetPasswordRepository.ResetPassword(ctx, userID, resetPassword)
	return err
}

func (r *resetPasswordUsecase) GetOTPByEmail(c context.Context, email string) (*domain.OtpSave, error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()
	otp, err := r.resetPasswordRepository.GetOTPByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return otp, nil
}
func (r *resetPasswordUsecase) DeleteOtp(c context.Context, email string) error {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()
	err := r.resetPasswordRepository.DeleteOtp(ctx, email)
	return err
}
