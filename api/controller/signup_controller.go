package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	b64 "encoding/base64"

	"github.com/dagota12/Loan-Tracker/bootstrap"
	"github.com/dagota12/Loan-Tracker/domain"
	"github.com/dagota12/Loan-Tracker/internal/tokenutil"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type SignupController struct {
	SignupUsecase domain.SignupUsecase
	Env           *bootstrap.Env
}

func (sc *SignupController) VerifyEmail(ctx *gin.Context) {
	Verificationtoken := ctx.Param("token")
	decodedToken, _ := b64.URLEncoding.DecodeString(Verificationtoken)

	valid, err := tokenutil.IsAuthorized(string(decodedToken), sc.Env.VerificationTokenSecret)

	fmt.Println(string(decodedToken))

	if !valid || err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	claims, err := tokenutil.ExtractUserClaimsFromToken(string(decodedToken), sc.Env.VerificationTokenSecret)
	userID := claims["id"].(string)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user, err := sc.SignupUsecase.GetUserById(context.TODO(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user.Active {
		ctx.JSON(http.StatusConflict, gin.H{"error": "user already verified!"})
		return
	}

	err = sc.SignupUsecase.ActivateUser(context.TODO(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})

}
func (sc *SignupController) Signup(ctx *gin.Context) {
	var request domain.SignupRequest

	err := ctx.ShouldBindJSON(&request)
	if err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	_, err = sc.SignupUsecase.GetUserByEmail(context.TODO(), request.Email)
	if err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		return
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	request.Password = string(encryptedPassword)

	//if there are no users in db this user owner
	IsOwner, err := sc.SignupUsecase.CanBeOwner(context.TODO())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var role string
	if IsOwner {
		role = "admin"
	} else {
		role = "user"
	}

	NewUser := domain.User{
		ID:        primitive.NewObjectID(),
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Password:  request.Password,
		Active:    false,
		IsOwner:   IsOwner,
		Role:      role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Tokens:    []string{},
	}

	VerificationToken, err := sc.SignupUsecase.CreateVerificationToken(&NewUser, sc.Env.VerificationTokenSecret, sc.Env.VerificationTokenExpiryMin)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	NewUser.VerifyToken = VerificationToken

	_, err = sc.SignupUsecase.Create(context.TODO(), &NewUser)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//send email
	encodedToken := b64.URLEncoding.EncodeToString([]byte(VerificationToken))
	//prinnt the user;
	log.Printf("user: %#v", NewUser)
	err = sc.SignupUsecase.SendVerificationEmail(NewUser.Email, encodedToken, sc.Env)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "email sent successfully, please verify your email"})

}
