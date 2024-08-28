package route

import (
	"time"

	"github.com/dagota12/Loan-Tracker/api/middleware"
	"github.com/dagota12/Loan-Tracker/bootstrap"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Setup(env *bootstrap.Env, timeout time.Duration, db *mongo.Database, gin *gin.Engine) {

	publicRouter := gin.Group("")

	// All Public APIs
	NewSignupRouter(env, timeout, db, publicRouter)
	NewAuthRouter(env, timeout, db, publicRouter)
	NewPasswordRouter(env, timeout, db, publicRouter)

	protectedRouter := gin.Group("")
	protectedRouter.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret))

	NewUsersRouter(env, timeout, db, protectedRouter)
	NewLoanRouter(env, timeout, db, protectedRouter)
}
