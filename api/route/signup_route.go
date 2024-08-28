package route

import (
	"log"
	"time"

	"github.com/dagota12/Loan-Tracker/api/controller"
	"github.com/dagota12/Loan-Tracker/bootstrap"
	"github.com/dagota12/Loan-Tracker/repository"
	"github.com/dagota12/Loan-Tracker/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewSignupRouter(env *bootstrap.Env, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	userRepo := repository.NewUserRepository(db)
	log.Println("[route] context timeout", env.ContextTimeout)
	signupUsecase := usecase.NewSignupUsecase(userRepo, time.Duration(env.ContextTimeout))
	sc := controller.SignupController{
		SignupUsecase: signupUsecase,
		Env:           env,
	}

	group.POST("/users/register", sc.Signup)
	group.GET("/users/verify-email/:token", sc.VerifyEmail)
}
