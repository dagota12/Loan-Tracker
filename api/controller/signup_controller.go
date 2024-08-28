package controller

import (
	"context"
	"fmt"
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

func (sc *SignupController) VerifyEmail(c *gin.Context) {
	Verificationtoken := c.Param("token")
	decodedToken, _ := b64.URLEncoding.DecodeString(Verificationtoken)

	valid, err := tokenutil.IsAuthorized(string(decodedToken), sc.Env.VerificationTokenSecret)

	fmt.Println(string(decodedToken))

	if !valid || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	claims, err := tokenutil.ExtractUserClaimsFromToken(string(decodedToken), sc.Env.VerificationTokenSecret)
	userID := claims["id"].(string)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user, err := sc.SignupUsecase.GetUserById(context.TODO(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user.Active {
		c.JSON(http.StatusConflict, gin.H{"error": "user already verified!"})
		return
	}

	err = sc.SignupUsecase.ActivateUser(context.TODO(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})

}
func (sc *SignupController) Signup(c *gin.Context) {
	var request domain.SignupRequest

	err := c.ShouldBindJSON(&request)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	user, err := sc.SignupUsecase.GetUserByEmail(c, request.Email)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		return
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	request.Password = string(encryptedPassword)

	IsOwner, err := sc.SignupUsecase.IsOwner(c, user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	}

	VerificationToken, err := sc.SignupUsecase.CreateVerificationToken(&NewUser, sc.Env.VerificationTokenSecret, sc.Env.VerificationTokenExpiryMin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user.VerifyToken = VerificationToken

	_, err = sc.SignupUsecase.Create(c, &user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//send email
	encodedToken := b64.URLEncoding.EncodeToString([]byte(VerificationToken))
	err = sc.SignupUsecase.SendVerificationEmail(user.Email, encodedToken, sc.Env)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, gin.H{"message": "email sent successfully, please verify your email"})

}
