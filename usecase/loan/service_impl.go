package loan

import (
	"SRATest/domain"
	"context"
	"fmt"
)

type loanUsecase struct {
	repo LoanRepository
}

func NewLoanUsecase(repo LoanRepository) LoanUsecase {
	return &loanUsecase{repo: repo}
}

func (u *loanUsecase) CreateLoan(ctx context.Context, id string, principal int64) (*domain.Loan, error) {
	total := int64(float64(principal) * 1.1)
	weekly := total / 50

	payments := make([]domain.Payment, 50)
	for i := 0; i < 50; i++ {
		payments[i] = domain.Payment{
			Week: i + 1,
			Paid: false,
		}
	}

	loan := &domain.Loan{
		ID:            id,
		Principal:     principal,
		InterestRate:  0.1,
		TotalWeeks:    50,
		WeeklyPayment: weekly,
		Payments:      payments,
	}

	return loan, u.repo.Save(ctx, loan)
}

func (u *loanUsecase) MakePayment(ctx context.Context, id string, amount int64) error {
	loan, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	weekly := loan.WeeklyPayment
	toPayWeeks := amount / weekly
	if amount%weekly != 0 {
		return fmt.Errorf("amount must be multiple of weekly payment: %d", weekly)
	}

	count := int(toPayWeeks)
	for i := 0; i < len(loan.Payments) && count > 0; i++ {
		if !loan.Payments[i].Paid {
			loan.Payments[i].Paid = true
			count--
		}
	}

	if count > 0 {
		return fmt.Errorf("amount exceeds remaining unpaid weeks")
	}

	return u.repo.Update(ctx, loan)
}

func (u *loanUsecase) GetOutstanding(ctx context.Context, id string) (int64, error) {
	loan, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return 0, err
	}

	unpaid := 0
	for _, p := range loan.Payments {
		if !p.Paid {
			unpaid++
		}
	}

	return int64(unpaid) * loan.WeeklyPayment, nil
}

func (u *loanUsecase) IsDelinquent(ctx context.Context, id string) (bool, error) {
	loan, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return false, err
	}

	missed := 0
	for _, p := range loan.Payments {
		if !p.Paid {
			missed++
			if missed >= 2 {
				return true, nil
			}
		} else {
			missed = 0
		}
	}

	return false, nil
}
