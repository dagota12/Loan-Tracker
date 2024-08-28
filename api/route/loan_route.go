package route

import (
	"log"
	"time"

	"github.com/dagota12/Loan-Tracker/api/controller"
	"github.com/dagota12/Loan-Tracker/api/middleware"
	"github.com/dagota12/Loan-Tracker/bootstrap"
	"github.com/dagota12/Loan-Tracker/repository"
	"github.com/dagota12/Loan-Tracker/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewLoanRouter(env *bootstrap.Env, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	// Initialize repository and usecase
	loanRepo := repository.NewLoanRepository(db)
	log.Println("[route] context timeout", env.ContextTimeout)
	loanUsecase := usecase.NewLoanUsecase(loanRepo, time.Duration(env.ContextTimeout))
	lc := controller.NewLoanController(loanUsecase, env)

	adminMiddleware := middleware.AdminMiddleware()

	// User routes
	loanGroup := group.Group("/loans")
	{
		loanGroup.GET("", lc.ViewUserLoans)      // get user loan
		loanGroup.POST("", lc.ApplyForLoan)      // Apply for a loan
		loanGroup.GET("/:id", lc.ViewLoanStatus) // View loan status by ID
	}

	// Admin routes
	adminLoanGroup := group.Group("/admin/loans")
	adminLoanGroup.Use(adminMiddleware)
	{
		adminLoanGroup.GET("", lc.ViewAllLoansWithPagination) // View all loans with pagination
		adminLoanGroup.PATCH("/:id", lc.ApproveRejectLoan)    // Approve or reject a loan
		adminLoanGroup.DELETE("/:id", lc.DeleteLoan)          // Delete a loan
	}
}
