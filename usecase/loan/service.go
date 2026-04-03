package loan

import (
	"SRATest/domain"
	"context"
)

type LoanRepository interface {
	Save(ctx context.Context, loan *domain.Loan) error
	FindByID(ctx context.Context, id string) (*domain.Loan, error)
	Update(ctx context.Context, loan *domain.Loan) error
}

type LoanUsecase interface {
	CreateLoan(ctx context.Context, id string, principal int64) (*domain.Loan, error)
	MakePayment(ctx context.Context, id string, amount int64) error
	GetOutstanding(ctx context.Context, id string) (int64, error)
	IsDelinquent(ctx context.Context, id string) (bool, error)
}
