package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Loan struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"` // Reference to the User
	Amount    float64            `json:"amount" bson:"amount" binding:"required"`
	Interest  float64            `json:"interest" bson:"interest" binding:"required"`
	Status    string             `json:"status" bson:"status"` // e.g., "pending", "approved", "rejected", "closed"
	StartDate time.Time          `json:"start_date" bson:"start_date"`
	EndDate   time.Time          `json:"end_date" bson:"end_date"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// LoanUsecase defines the interface for loan-related business logic.
type LoanUsecase interface {
	// ApplyForLoan handles the loan application process.
	ApplyForLoan(ctx context.Context, loan *Loan) error

	// ViewLoanStatus retrieves the status of a specific loan.
	ViewLoanStatus(ctx context.Context, loanID primitive.ObjectID) (*Loan, error)

	// ViewUserLoans retrieves all loans associated with a specific user.
	ViewUserLoans(ctx context.Context, userID primitive.ObjectID) ([]Loan, error)

	// ViewAllLoans retrieves all loans, with optional filters (for admin use).
	ViewAllLoans(ctx context.Context) ([]Loan, error)
	ViewAllLoansWithPagination(ctx context.Context, status string, order string, page int64, limit int64) ([]Loan, PaginationMetadata, error)
	// ApproveRejectLoan handles the approval or rejection of a loan (for admin use).
	ApproveRejectLoan(ctx context.Context, loanID primitive.ObjectID, status string) error

	// DeleteLoan deletes a specific loan (for admin use).
	DeleteLoan(ctx context.Context, loanID primitive.ObjectID) error
}

// LoanRepository defines the interface for loan-related database operations.
type LoanRepository interface {
	// CreateLoan adds a new loan to the database.
	CreateLoan(ctx context.Context, loan *Loan) error

	// GetLoanByID retrieves a loan by its ID.
	GetLoanByID(ctx context.Context, loanID primitive.ObjectID) (*Loan, error)
	GetAllLoansWithPagination(ctx context.Context, status string, order string, page int64, limit int64) ([]Loan, PaginationMetadata, error)
	// GetLoansByUserID retrieves all loans associated with a specific user.
	GetLoansByUserID(ctx context.Context, userID primitive.ObjectID) ([]Loan, error)

	// GetAllLoans retrieves all loans, with optional filters (for admin use).
	GetAllLoans(ctx context.Context) ([]Loan, error)

	// UpdateLoanStatus updates the status of a loan (for admin use).
	UpdateLoanStatus(ctx context.Context, loanID primitive.ObjectID, status string) error

	// DeleteLoan deletes a loan by its ID (for admin use).
	DeleteLoan(ctx context.Context, loanID primitive.ObjectID) error
}
