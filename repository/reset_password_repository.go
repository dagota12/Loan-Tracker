package repository

import (
	"context"
	"errors"
	"log"

	"github.com/dagota12/Loan-Tracker/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrOtpNotFound = errors.New("otp not found")
)

type resetPasswordRepository struct {
	database        *mongo.Database
	usersCollection string
	resetCollection string
}

func NewResetPasswordRepository(db *mongo.Database, userCollection string, resetCollection string) domain.ResetPasswordRepository {
	return &resetPasswordRepository{
		database:        db,
		usersCollection: userCollection,
		resetCollection: resetCollection,
	}
}

func (rp *resetPasswordRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	collection := rp.database.Collection(rp.usersCollection)
	var user domain.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		log.Println("[repo] restePwd", err.Error())
		return nil, ErrUserNotFound
	}
	return &user, err
}
func (rp *resetPasswordRepository) ResetPassword(ctx context.Context, userID string, resetPassword *domain.ResetPasswordRequest) error {

	collection := rp.database.Collection(rp.usersCollection)
	ObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return ErrInvalidID
	}
	res, err := collection.UpdateOne(ctx, bson.M{"_id": ObjID}, bson.M{"$set": bson.M{"password": resetPassword.NewPassword}})
	if err != nil {
		return err
	}
	if res.MatchedCount < 1 {
		return ErrUserNotFound
	}
	return nil
}

func (rp *resetPasswordRepository) SaveOtp(ctx context.Context, otp *domain.OtpSave) error {
	collection := rp.database.Collection(rp.resetCollection)

	_, err := collection.InsertOne(ctx, otp)

	if err != nil {
		return err
	}

	return err
}

func (rp *resetPasswordRepository) GetOTPByEmail(ctx context.Context, email string) (*domain.OtpSave, error) {

	collection := rp.database.Collection(rp.resetCollection)
	var otp domain.OtpSave

	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&otp)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrOtpNotFound
	}

	if err != nil {
		return nil, err
	}

	return &otp, err
}

func (rp *resetPasswordRepository) DeleteOtp(ctx context.Context, email string) error {

	collection := rp.database.Collection(rp.resetCollection)

	_, err := collection.DeleteOne(ctx, bson.M{"email": email})

	if err != nil {
		return err
	}

	return err
}
