package controller

import (
	"context"
	"net/http"

	"github.com/dagota12/Loan-Tracker/domain"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUsecase domain.UserUsecase
}

func NewUserController(usercase domain.UserUsecase) *UserController {
	return &UserController{
		userUsecase: usercase,
	}
}

func (uc *UserController) Register(ctx *gin.Context) {
	userdata := domain.UserForm{}
	if err := ctx.ShouldBindJSON(&userdata); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := domain.User{
		FirstName: userdata.FirstName,
		LastName:  userdata.LastName,
		Email:     userdata.Email,
		Password:  userdata.Password,
	}
	user, err := uc.userUsecase.Create(context.Background(), user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, user)
}

func (uc *UserController) GetUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing user ID"})
		return
	}

	user, err := uc.userUsecase.GetByID(context.Background(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (uc *UserController) GetAllUsers(ctx *gin.Context) {
	users, err := uc.userUsecase.GetAll(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (uc *UserController) UpdateUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing user ID"})
		return
	}

	var userUpdate domain.UserUpdate
	if err := ctx.ShouldBindJSON(&userUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := uc.userUsecase.Update(context.Background(), userID, userUpdate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (uc *UserController) DeleteUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing user ID"})
		return
	}

	err := uc.userUsecase.Delete(context.Background(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

func (uc *UserController) UpdatePassword(ctx *gin.Context) {
	userID := ctx.MustGet("x-user-id").(string)
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing user ID"})
		return
	}

	var passwordUpdate domain.UpdatePassword
	if err := ctx.ShouldBindJSON(&passwordUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := uc.userUsecase.UpdateUserPassword(context.Background(), userID, passwordUpdate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "password updated successfully"})
}
