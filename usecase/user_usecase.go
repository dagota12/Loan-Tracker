package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/dagota12/Loan-Tracker/bootstrap"
	"github.com/dagota12/Loan-Tracker/domain"
	"github.com/dagota12/Loan-Tracker/internal/security"
)

type userUsecase struct {
	UserRepo       domain.UserRepository
	contextTimeout time.Duration
	Env            *bootstrap.Env
}

func NewUserUsecase(repo domain.UserRepository, env *bootstrap.Env) domain.UserUsecase {
	return &userUsecase{
		UserRepo:       repo,
		contextTimeout: time.Duration(env.ContextTimeout) * time.Second,
		Env:            env,
	}
}

// Create implements domain.UserUsecase.
// Create creates a new user after checking for conflicting emails.
func (uc *userUsecase) Create(ctx context.Context, user domain.User) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	// Check if the email already exists
	existingUser, _ := uc.UserRepo.GetByEmail(ctx, user.Email)
	if existingUser.Email != "" {
		return domain.User{}, errors.New("email already in use")
	}

	// Hash the password before creating the user
	hashedPassword, err := security.HashPassword(user.Password)
	if err != nil {
		return domain.User{}, err
	}
	user.Password = hashedPassword

	createdUser, err := uc.UserRepo.Create(ctx, user)
	if err != nil {
		return domain.User{}, err
	}

	return createdUser, nil
}

// Delete implements domain.UserUsecase.
// Delete deletes a user by their ID.
func (uc *userUsecase) Delete(ctx context.Context, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	err := uc.UserRepo.Delete(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

// GetAll retrieves all users from the repository.
func (uc *userUsecase) GetAll(ctx context.Context) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	users, err := uc.UserRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetByEmail implements domain.UserUsecase.
// GetByEmail retrieves a user by their email.
func (uc *userUsecase) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	user, err := uc.UserRepo.GetByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

// GetByID implements domain.UserUsecase.
// GetByID retrieves a user by their ID.
func (uc *userUsecase) GetByID(ctx context.Context, userID string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	user, err := uc.UserRepo.GetByID(ctx, userID)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

// Update implements domain.UserUsecase.
// Update updates user details by user ID.
func (uc *userUsecase) Update(ctx context.Context, userID string, user domain.UserUpdate) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	updatedUser, err := uc.UserRepo.Update(ctx, userID, user)
	if err != nil {
		return domain.User{}, err
	}

	return updatedUser, nil
}

// ResetUserPassword implements domain.UserUsecase.
// ResetUserPassword resets the user's password using a reset token or temporary password.
func (uc *userUsecase) ResetUserPassword(ctx context.Context, userID string, resetPassword domain.ResetPasswordRequest) error {
	// Set up a context with a timeout
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	// Check if the user is active before proceeding with password reset
	active, err := uc.UserRepo.IsUserActive(ctx, userID)
	if err != nil {
		return err
	}
	if !active {
		return errors.New("user is not active or does not exist")
	}

	// Call the repository method to reset the password
	err = uc.UserRepo.ResetUserPassword(ctx, userID, resetPassword)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUserPassword updates the user's password.
func (uc *userUsecase) UpdateUserPassword(ctx context.Context, userID string, updatePassword domain.UpdatePassword) error {
	// Set up a context with a timeout
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	// Verify if the user is active
	active, err := uc.UserRepo.IsUserActive(ctx, userID)
	if err != nil {
		return err
	}
	if !active {
		return errors.New("user is not active or does not exist")
	}

	// Validate the current password
	user, err := uc.UserRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if security.CheckPasswordHash(updatePassword.OldPassword, user.Password) {
		return errors.New("current password is incorrect")
	}

	// Update the password using the repository method
	// Hash the new password before updating
	hashedPassword, err := security.HashPassword(updatePassword.NewPassword)
	if err != nil {
		return err
	}
	updatePassword.NewPassword = hashedPassword
	err = uc.UserRepo.UpdateUserPassword(ctx, userID, updatePassword)
	if err != nil {
		return err
	}

	return nil
}
