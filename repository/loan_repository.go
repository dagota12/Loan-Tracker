package repository

import (
	"context"
	"errors"
	"time"

	"github.com/dagota12/Loan-Tracker/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type loanRepository struct {
	collection *mongo.Collection
}

func NewLoanRepository(db *mongo.Database) domain.LoanRepository {
	return &loanRepository{
		collection: db.Collection("loans"),
	}
}

func (lr *loanRepository) CreateLoan(ctx context.Context, loan *domain.Loan) error {
	loan.CreatedAt = time.Now()
	loan.UpdatedAt = time.Now()
	loan.Status = "pending" // default status

	_, err := lr.collection.InsertOne(ctx, loan)
	return err
}

// GetLoanByID retrieves a loan by its ID.
func (lr *loanRepository) GetLoanByID(ctx context.Context, loanID primitive.ObjectID) (*domain.Loan, error) {
	var loan domain.Loan
	err := lr.collection.FindOne(ctx, bson.M{"_id": loanID}).Decode(&loan)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("loan not found")
	}
	return &loan, err
}

// GetLoansByUserID retrieves all loans associated with a specific user.
func (lr *loanRepository) GetLoansByUserID(ctx context.Context, userID primitive.ObjectID) ([]domain.Loan, error) {
	var loans []domain.Loan
	cursor, err := lr.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var loan domain.Loan
		if err := cursor.Decode(&loan); err != nil {
			return nil, err
		}
		loans = append(loans, loan)
	}

	return loans, cursor.Err()
}

// GetAllLoansWithPagination retrieves all loans with pagination support (for admin use).
func (lr *loanRepository) GetAllLoansWithPagination(ctx context.Context, status string, order string, page int64, limit int64) ([]domain.Loan, domain.PaginationMetadata, error) {
	filter := bson.M{}
	if status != "" {
		filter["status"] = status
	}

	findOptions := options.Find()
	if order == "asc" {
		findOptions.SetSort(bson.M{"created_at": 1})
	} else {
		findOptions.SetSort(bson.M{"created_at": -1})
	}

	// Pagination options
	skip := (page - 1) * limit
	findOptions.SetSkip(skip)
	findOptions.SetLimit(limit)

	var loans []domain.Loan
	cursor, err := lr.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, domain.PaginationMetadata{}, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var loan domain.Loan
		if err := cursor.Decode(&loan); err != nil {
			return nil, domain.PaginationMetadata{}, err
		}
		loans = append(loans, loan)
	}

	// Get total count of documents that match the filter
	totalCount, err := lr.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, domain.PaginationMetadata{}, err
	}

	// Calculate total pages
	totalPages := totalCount / limit
	if totalCount%limit != 0 {
		totalPages++
	}

	// Create PaginationMetadata
	paginationMetadata := domain.PaginationMetadata{
		TotalCount:   totalCount,
		TotalPages:   totalPages,
		CurrentPage:  page,
		ItemsPerPage: limit,
	}

	return loans, paginationMetadata, cursor.Err()
}

func (lr *loanRepository) GetAllLoans(ctx context.Context) ([]domain.Loan, error) {
	filter := bson.M{}

	var loans []domain.Loan
	cursor, err := lr.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var loan domain.Loan
		if err := cursor.Decode(&loan); err != nil {
			return nil, err
		}
		loans = append(loans, loan)
	}

	return loans, cursor.Err()
}

// UpdateLoanStatus updates the status of a loan (for admin use).
func (lr *loanRepository) UpdateLoanStatus(ctx context.Context, loanID primitive.ObjectID, status string) error {
	update := bson.M{"$set": bson.M{"status": status, "updated_at": time.Now()}}
	result, err := lr.collection.UpdateOne(ctx, bson.M{"_id": loanID}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("loan not found")
	}
	return nil
}

// DeleteLoan deletes a loan by its ID (for admin use).
func (lr *loanRepository) DeleteLoan(ctx context.Context, loanID primitive.ObjectID) error {
	result, err := lr.collection.DeleteOne(ctx, bson.M{"_id": loanID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("loan not found")
	}
	return nil
}
