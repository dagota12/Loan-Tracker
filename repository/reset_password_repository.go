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

func (rp *resetPasswordRepository) GetUserByEmail(c context.Context, email string) (*domain.User, error) {
	collection := rp.database.Collection(rp.usersCollection)
	var user domain.User
	err := collection.FindOne(c, bson.M{"email": email}).Decode(&user)
	if err != nil {
		log.Println("[repo] restePwd", err.Error())
		return nil, ErrUserNotFound
	}
	return &user, err
}
func (rp *resetPasswordRepository) ResetPassword(c context.Context, userID string, resetPassword *domain.ResetPasswordRequest) error {

	collection := rp.database.Collection(rp.usersCollection)
	ObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return ErrInvalidID
	}
	res, err := collection.UpdateOne(c, bson.M{"_id": ObjID}, bson.M{"$set": bson.M{"password": resetPassword.NewPassword}})
	if err != nil {
		return err
	}
	if res.MatchedCount < 1 {
		return ErrUserNotFound
	}
	return nil
}

func (rp *resetPasswordRepository) SaveOtp(c context.Context, otp *domain.OtpSave) error {
	collection := rp.database.Collection(rp.resetCollection)

	_, err := collection.InsertOne(c, otp)

	if err != nil {
		return err
	}

	return err
}

func (rp *resetPasswordRepository) GetOTPByEmail(c context.Context, email string) (*domain.OtpSave, error) {

	collection := rp.database.Collection(rp.resetCollection)
	var otp domain.OtpSave

	err := collection.FindOne(c, bson.M{"email": email}).Decode(&otp)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &otp, err
}

func (rp *resetPasswordRepository) DeleteOtp(c context.Context, email string) error {

	collection := rp.database.Collection(rp.resetCollection)

	_, err := collection.DeleteOne(c, bson.M{"email": email})

	if err != nil {
		return err
	}

	return err
}
