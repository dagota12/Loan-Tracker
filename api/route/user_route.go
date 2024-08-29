package route

import (
	"time"

	"github.com/dagota12/Loan-Tracker/api/controller"
	"github.com/dagota12/Loan-Tracker/api/middleware"
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

	group.GET("/users/profile", userController.GetUserProfile)

	admin := group.Group("")
	admin.Use(middleware.AdminMiddleware())
	admin.GET("/admin/users/:id", userController.GetUser)
	admin.GET("/admin/users", userController.GetAllUsers)
	admin.DELETE("admin/users/:id", userController.DeleteUser)

	//protected routes
	protected := group.Group("")
	protected.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret))
	protected.POST("/users/update-password", userController.UpdatePassword)
}
