package controller

import (
	"log"
	"net/http"

	"github.com/dagota12/Loan-Tracker/bootstrap"
	"github.com/dagota12/Loan-Tracker/domain"
	"github.com/dagota12/Loan-Tracker/internal/tokenutil"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	AuthUsecase domain.AuthUsecase
	Env         *bootstrap.Env
}

func NewAuthController(usecase domain.AuthUsecase, env *bootstrap.Env) *AuthController {
	return &AuthController{
		AuthUsecase: usecase,
		Env:         env,
	}
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

func (ac *AuthController) RefreshToken(c *gin.Context) {
	var request domain.RefreshTokenRequest

	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if valid, err := tokenutil.IsAuthorized(string(request.RefreshToken), ac.Env.RefreshTokenSecret); !valid || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	claims, err := tokenutil.ExtractUserClaimsFromToken(request.RefreshToken, ac.Env.RefreshTokenSecret)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	log.Printf("[ctrl] userclaims: %#v", "claims")

	userID := claims.ID
	user, err := ac.AuthUsecase.GetUserByID(c, userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := ac.AuthUsecase.CreateAccessToken(user, ac.Env.AccessTokenSecret, ac.Env.AccessTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	refreshToken, err := ac.AuthUsecase.CreateRefreshToken(user, ac.Env.RefreshTokenSecret, ac.Env.RefreshTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	refreshTokenResponse := domain.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, refreshTokenResponse)
}
