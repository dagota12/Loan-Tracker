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

func NewAuthRouter(env *bootstrap.Env, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	userRepo := repository.NewUserRepository(db)
	authUsecase := usecase.NewAuthUsease(userRepo)
	authController := controller.NewAuthController(authUsecase, env)

	group.POST("/users/login", authController.Login)
	group.POST("/users/token/refresh", authController.RefreshToken)
}
