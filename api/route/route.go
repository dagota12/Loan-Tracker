package route

import (
	"time"

	"github.com/dagota12/Loan-Tracker/bootstrap"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Setup(env *bootstrap.Env, timeout time.Duration, db *mongo.Database, gin *gin.Engine) {

	publicRouter := gin.Group("")

	// All Public APIs
	NewSignupRouter(env, timeout, db, publicRouter)
	NewLoginRouter(env, timeout, db, publicRouter)
	NewRefreshTokenRouter(env, timeout, db, publicRouter)
	NewPublicBlogsRouter(env, timeout, db, publicRouter)
	NewPublicResetPasswordRouter(env, timeout, db, publicRouter)

	// Static files
	// NewPublicFileRouter(env, publicRouter)

	protectedRouter := gin.Group("")
	protectedRouter.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret))

	// All Protected APIs
	NewProtectedBlogsRouter(env, timeout, db, protectedRouter)
	NewUsersRouter(env, timeout, db, protectedRouter, cloudinary)
}
