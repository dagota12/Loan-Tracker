package route

import (
	"time"

	"github.com/dagota12/Loan-Tracker/api/controller"
	"github.com/dagota12/Loan-Tracker/bootstrap"
	"github.com/dagota12/Loan-Tracker/repository"
	"github.com/dagota12/Loan-Tracker/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewPasswordRouter(env *bootstrap.Env, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	userRepo := repository.NewResetPasswordRepository(db, "users", "password-reset")
	userUsecase := usecase.NewResetPasswordUsecase(userRepo, timeout)
	userController := controller.NewResetPasswordController(env, userUsecase)

	group.POST("/users/reset-password", userController.ResetPassword)
	group.POST("/users/forgot-password", userController.ForgotPassword)
}
