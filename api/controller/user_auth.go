package controller

import (
	"net/http"

	"github.com/dagota12/Loan-Tracker/bootstrap"
	"github.com/dagota12/Loan-Tracker/domain"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	AuthUsecase domain.AuthUsecase
	Env         *bootstrap.Env
}

func (ac *AuthController) Login(ctx *gin.Context) {
	var request domain.LoginRequest

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Request body cannot be empty"})
		return
	}

	user, err := ac.AuthUsecase.GetUserByEmail(ctx, request.Email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if !user.Active {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user is not active"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)) != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credential"})
		return
	}

	accessToken, err := ac.AuthUsecase.CreateAccessToken(user, ac.Env.AccessTokenSecret, ac.Env.AccessTokenExpiryHour)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	refreshToken, err := ac.AuthUsecase.CreateRefreshToken(user, ac.Env.RefreshTokenSecret, ac.Env.RefreshTokenExpiryHour)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = ac.AuthUsecase.UpdateRefreshToken(ctx.Request.Context(), user.ID.Hex(), refreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// return response
	loginResponse := domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	ctx.JSON(http.StatusOK, loginResponse)
}
