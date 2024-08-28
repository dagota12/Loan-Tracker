package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/dagota12/Loan-Tracker/bootstrap"
	"github.com/dagota12/Loan-Tracker/domain"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoanController struct {
	LoanUsecase domain.LoanUsecase
	Env         *bootstrap.Env
}

// NewLoanController creates a new instance of LoanController.
func NewLoanController(LoanUsecase domain.LoanUsecase, env *bootstrap.Env) *LoanController {
	return &LoanController{
		LoanUsecase: LoanUsecase,
		Env:         env,
	}
}

// ApplyForLoan handles the POST /loans endpoint to submit a loan application.
func (lc *LoanController) ApplyForLoan(c *gin.Context) {
	var loan domain.Loan
	userID := c.MustGet("x-user-id").(string)
	if err := c.ShouldBindJSON(&loan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var err error
	loan.UserID, err = primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
	}
	creatredLoan, err := lc.LoanUsecase.ApplyForLoan(c.Request.Context(), &loan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to apply for loan"})
		return
	}

	c.JSON(http.StatusCreated, creatredLoan)
}
func (lc *LoanController) ViewUserLoans(c *gin.Context) {
	userID := c.MustGet("x-user-id").(string)
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	loans, err := lc.LoanUsecase.ViewUserLoans(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user loans"})
		return
	}

	c.JSON(http.StatusOK, loans)
}

// ViewUserLoansWithPagination handles the GET /users/{id}/loans with pagination)
// ViewLoanStatus handles the GET /loans/{id} endpoint to retrieve a loan's status.
func (lc *LoanController) ViewLoanStatus(c *gin.Context) {
	loanID := c.Param("id")
	id, err := primitive.ObjectIDFromHex(loanID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
		return
	}

	loan, err := lc.LoanUsecase.ViewLoanStatus(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve loan status"})
		return
	}

	c.JSON(http.StatusOK, loan)
}

// ViewAllLoans handles the GET /admin/loans endpoint to retrieve all loan applications.
func (lc *LoanController) ViewAllLoans(c *gin.Context) {

	loans, err := lc.LoanUsecase.ViewAllLoans(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve loans"})
		return
	}

	c.JSON(http.StatusOK, loans)
}

// ViewAllLoansWithPagination handles the GET /admin/loans with pagination
func (lc *LoanController) ViewAllLoansWithPagination(c *gin.Context) {
	status := c.Query("status")
	order := c.Query("order")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	if page < 1 || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page or limit parameter"})
		return
	}
	loans, pagination, err := lc.LoanUsecase.ViewAllLoansWithPagination(c.Request.Context(), status, order, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve loans with pagination"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"loans": loans, "pagination": pagination})
}

// ApproveRejectLoan handles the PATCH /admin/loans/{id}/status endpoint to approve or reject a loan.
func (lc *LoanController) ApproveRejectLoan(c *gin.Context) {
	loanID := c.Param("id")
	status := c.Query("status")
	id, err := primitive.ObjectIDFromHex(loanID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
		return
	}
	log.Printf("[ctrl] on calling ApproveRejectLoan: %s %s", loanID, status)
	err = lc.LoanUsecase.ApproveRejectLoan(c.Request.Context(), id, status)
	if err != nil {
		log.Println("[ctrl] on ApproveRejectLoan: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update loan status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Loan status updated successfully"})
}

// DeleteLoan handles the DELETE /admin/loans/{id} endpoint to delete a loan.
func (lc *LoanController) DeleteLoan(c *gin.Context) {
	loanID := c.Param("id")
	id, err := primitive.ObjectIDFromHex(loanID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
		return
	}

	err = lc.LoanUsecase.DeleteLoan(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete loan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Loan deleted successfully"})
}
