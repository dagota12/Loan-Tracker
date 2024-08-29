package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/dagota12/Loan-Tracker/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// loanUsecase struct implementing the LoanUsecase interface.
type loanUsecase struct {
	loanRepo domain.LoanRepository
	timeout  time.Duration
}

// NewLoanUsecase creates a new instance of loanUsecase.
func NewLoanUsecase(loanRepo domain.LoanRepository, timeout time.Duration) domain.LoanUsecase {
	return &loanUsecase{
		loanRepo: loanRepo,
		timeout:  timeout,
	}
}

func (lu *loanUsecase) ApplyForLoan(ctx context.Context, loan *domain.Loan) (*domain.Loan, error) {
	// Set initial loan status and timestamps
	loan.Status = "pending"
	loan.CreatedAt = time.Now()
	loan.UpdatedAt = time.Now()

	// Call repository to create a loan
	return lu.loanRepo.CreateLoan(ctx, loan)
}

func (lu *loanUsecase) ViewLoanStatus(ctx context.Context, loanID primitive.ObjectID) (*domain.Loan, error) {
	// Retrieve the loan by ID from the repository
	return lu.loanRepo.GetLoanByID(ctx, loanID)
}

func (lu *loanUsecase) ViewUserLoans(ctx context.Context, userID primitive.ObjectID) ([]domain.Loan, error) {
	// Retrieve all loans for a specific user
	return lu.loanRepo.GetLoansByUserID(ctx, userID)
}

func (lu *loanUsecase) ViewAllLoansWithPagination(ctx context.Context, status string, order string, page int64, limit int64) ([]domain.Loan, domain.PaginationMetadata, error) {
	// Retrieve all loans with pagination, status, and order filters
	return lu.loanRepo.GetAllLoansWithPagination(ctx, status, order, page, limit)
}

func (lu *loanUsecase) ApproveRejectLoan(ctx context.Context, loanID primitive.ObjectID, status string) error {
	// Validate status
	if status != "approved" && status != "rejected" {
		return errors.New("invalid status")
	}

	// Retrieve loan to update its status
	loan, err := lu.loanRepo.GetLoanByID(ctx, loanID)
	if err != nil {
		return err
	}

	// Update the loan status
	loan.Status = status
	loan.UpdatedAt = time.Now()

	// Save the updated loan status back to the repository
	return lu.loanRepo.UpdateLoanStatus(ctx, loanID, status)
}

func (lu *loanUsecase) DeleteLoan(ctx context.Context, loanID primitive.ObjectID) error {
	// Delete the loan by ID from the repository
	return lu.loanRepo.DeleteLoan(ctx, loanID)
}
func (lu *loanUsecase) ViewAllLoans(ctx context.Context) ([]domain.Loan, error) {
	// Retrieve all loans from the repository
	return lu.loanRepo.GetAllLoans(ctx)
}
