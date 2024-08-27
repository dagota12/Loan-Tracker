package repository

import (
	"context"
	"log"

	"errors"

	"github.com/dagota12/Loan-Tracker/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidID    = errors.New("user not id invalid")
)

type userRepository struct {
	db    *mongo.Database
	users *mongo.Collection
}

func NewUserRepository(db *mongo.Database) domain.UserRepository {
	return &userRepository{
		db:    db,
		users: db.Collection("users"),
	}
}

// ActivateUser implements domain.UserRepository.
func (ur *userRepository) ActivateUser(ctx context.Context, userID string) error {
	ObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return ErrInvalidID
	}

	filter := bson.M{"_id": ObjID}
	res, err := ur.users.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"active": true}})
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return ErrUserNotFound
	}
	return nil
}

// Create implements domain.UserRepository.
func (ur *userRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	res, err := ur.users.InsertOne(ctx, user)
	if err != nil {
		return domain.User{}, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

// Delete implements domain.UserRepository.
func (ur *userRepository) Delete(ctx context.Context, userID string) error {

	res, err := ur.users.DeleteOne(ctx, userID)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrUserNotFound
	}
	return nil
}

// GetAll implements domain.UserRepository.
func (ur *userRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	//get all users
	cursor, err := ur.users.Find(ctx, bson.M{})
	//check if there is an error
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	users := make([]domain.User, 0)
	err = cursor.All(ctx, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetByEmail implements domain.UserRepository.
func (ur *userRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	//get user by email
	filter := bson.M{"email": email}
	user := domain.User{}
	err := ur.users.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {

			return domain.User{}, ErrUserNotFound
		}
		log.Println("[repo] user get by email", err)
		return domain.User{}, err
	}
	return user, nil
}

// GetByID implements domain.UserRepository.
func (ur *userRepository) GetByID(ctx context.Context, userID string) (domain.User, error) {
	ObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return domain.User{}, ErrInvalidID
	}

	filter := bson.M{"_id": ObjID}
	user := domain.User{}
	err = ur.users.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {

			return domain.User{}, ErrUserNotFound
		}
		log.Println("[repo] user get by email", err)
		return domain.User{}, err
	}
	return user, nil
}

// IsUserActive implements domain.UserRepository.
func (ur *userRepository) IsUserActive(ctx context.Context, userID string) (bool, error) {
	ObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, ErrInvalidID
	}

	filter := bson.M{"_id": ObjID}
	user := domain.User{}
	//check if user exists
	err = ur.users.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, ErrUserNotFound
		}
		log.Println("[repo] on is user active", err)
		return false, err
	}
	return user.Active, nil
}

func (ur *userRepository) IsOwner(ctx context.Context, userID string) (bool, error) {
	ObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, ErrInvalidID
	}

	filter := bson.M{"_id": ObjID}
	user := domain.User{}
	//check if user exists
	err = ur.users.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, ErrUserNotFound
		}
		log.Println("[repo] on is user active", err)
		return false, err
	}
	return user.Active, nil

}

// RefreshTokenExist implements domain.UserRepository.
func (ur *userRepository) RefreshTokenExist(ctx context.Context, userID string, refreshToken string) (bool, error) {
	ObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, ErrInvalidID
	}

	filter := bson.M{"_id": ObjID, "refresh_tokens": refreshToken}
	res, err := ur.users.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return res > 0, nil

}

// ResetUserPassword implements domain.UserRepository.
func (ur *userRepository) ResetUserPassword(ctx context.Context, userID string, resetPassword domain.ResetPasswordRequest) error {
	ObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return ErrInvalidID
	}

	filter := bson.M{"_id": ObjID}
	update := bson.M{"$set": bson.M{"password": resetPassword.NewPassword}}
	res, err := ur.users.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return ErrUserNotFound
	}
	return nil

}

// RevokeRefreshToken implements domain.UserRepository.
func (ur *userRepository) RevokeRefreshToken(ctx context.Context, userID string, refreshToken string) error {
	ObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return ErrInvalidID
	}

	filter := bson.M{"_id": ObjID}
	update := bson.M{"$pull": bson.M{"refresh_tokens": refreshToken}}
	res, err := ur.users.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return ErrUserNotFound
	}
	return nil
}

// Update updates the user information by user ID.
func (ur *userRepository) Update(ctx context.Context, userID string, user domain.UserUpdate) (domain.User, error) {
	// Convert string ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return domain.User{}, ErrInvalidID
	}

	// Define update fields using bson.D
	updateFields := bson.M{"$set": user}
	// Create update document
	update := bson.D{
		{Key: "$set", Value: updateFields},
	}

	// Perform the update operation
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedUser domain.User
	err = ur.users.FindOneAndUpdate(ctx, bson.M{"_id": objID}, update, opts).Decode(&updatedUser)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, ErrUserNotFound
		}
		return domain.User{}, err
	}

	return updatedUser, nil
}

// UpdateRefreshToken updates the refresh token of the user.
func (ur *userRepository) UpdateRefreshToken(ctx context.Context, userID string, refreshToken string) error {
	// Convert string ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return ErrInvalidID
	}

	// Define update fields
	update := bson.D{
		{Key: "$push", Value: bson.D{{Key: "refresh_tokens", Value: refreshToken}}},
	}

	// Perform the update operation
	result, err := ur.users.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return ErrUserNotFound
	}

	return nil
}

// UpdateUserPassword updates the user's password.
func (ur *userRepository) UpdateUserPassword(ctx context.Context, userID string, updatePassword domain.UpdatePassword) error {
	// Convert string ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return ErrInvalidID
	}

	// Define update fields
	update := bson.D{
		{Key: "$set", Value: bson.D{{Key: "password", Value: updatePassword.NewPassword}}},
	}

	// Perform the update operation
	result, err := ur.users.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return ErrUserNotFound
	}

	return nil
}
