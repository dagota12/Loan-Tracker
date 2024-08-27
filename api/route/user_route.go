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

func NewUsersRouter(env *bootstrap.Env, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo, env)
	userController := controller.NewUserController(userUsecase)

	group.GET("/users", userController.GetAllUsers)
	group.GET("/users/:id", userController.GetUser)
	group.DELETE("/users/:id", userController.DeleteUser)
}
