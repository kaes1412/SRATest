package loan_test

import (
	"SRATest/domain"
	"SRATest/usecase/loan"
	"context"
	"fmt"
	"testing"
)

// Mock repository
type mockRepo struct {
	store map[string]*domain.Loan
}

func (m *mockRepo) Save(ctx context.Context, loan *domain.Loan) error {
	if m.store == nil {
		m.store = make(map[string]*domain.Loan)
	}
	// copy loan supaya tidak sharing pointer
	cpy := *loan
	m.store[loan.ID] = &cpy
	return nil
}

func (m *mockRepo) FindByID(ctx context.Context, id string) (*domain.Loan, error) {
	if m.store == nil {
		return nil, fmt.Errorf("not found")
	}
	loan, ok := m.store[id]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	// return copy supaya test aman
	cpy := *loan
	return &cpy, nil
}

func (m *mockRepo) Update(ctx context.Context, loan *domain.Loan) error {
	if m.store == nil {
		return fmt.Errorf("not found")
	}
	// update dengan copy
	cpy := *loan
	m.store[loan.ID] = &cpy
	return nil
}

func TestCreateLoan(t *testing.T) {
	ctx := context.Background()
	repo := &mockRepo{}
	usecase := loan.NewLoanUsecase(repo)

	loanID := "loan1"
	principal := int64(5000000)
	l, err := usecase.CreateLoan(ctx, loanID, principal)
	if err != nil {
		t.Fatalf("CreateLoan failed: %v", err)
	}

	if l.ID != loanID {
		t.Errorf("Expected ID %s, got %s", loanID, l.ID)
	}
	if l.Principal != principal {
		t.Errorf("Expected Principal %d, got %d", principal, l.Principal)
	}
	expectedWeekly := int64(float64(principal)*1.1) / 50
	if l.WeeklyPayment != expectedWeekly {
		t.Errorf("Expected weekly payment %d, got %d", expectedWeekly, l.WeeklyPayment)
	}
}

func TestMakePaymentAndOutstanding(t *testing.T) {
	ctx := context.Background()
	repo := &mockRepo{}
	usecase := loan.NewLoanUsecase(repo)

	loanID := "loan2"
	principal := int64(5000000)
	usecase.CreateLoan(ctx, loanID, principal)

	weekly := int64(float64(principal)*1.1) / 50
	err := usecase.MakePayment(ctx, loanID, weekly*2)
	if err != nil {
		t.Fatalf("MakePayment failed: %v", err)
	}

	outstanding, _ := usecase.GetOutstanding(ctx, loanID)
	expected := int64(50-2) * weekly
	if outstanding != expected {
		t.Errorf("Expected outstanding %d, got %d", expected, outstanding)
	}
}

func TestIsDelinquent(t *testing.T) {
	ctx := context.Background()
	repo := &mockRepo{}
	usecase := loan.NewLoanUsecase(repo)

	// Buat loan 5 minggu untuk test lebih mudah
	loanID := "loan3"
	principal := int64(5000000)
	totalWeeks := 5
	weekly := int64(float64(principal)*1.1) / int64(totalWeeks)

	// Custom loan dengan 5 minggu
	payments := make([]domain.Payment, totalWeeks)
	for i := 0; i < totalWeeks; i++ {
		payments[i] = domain.Payment{Week: i + 1, Paid: false}
	}

	loan := &domain.Loan{
		ID:            loanID,
		Principal:     principal,
		InterestRate:  0.1,
		TotalWeeks:    totalWeeks,
		WeeklyPayment: weekly,
		Payments:      payments,
	}
	repo.Save(ctx, loan)

	usecase.MakePayment(ctx, loanID, weekly)
	delinquent, _ := usecase.IsDelinquent(ctx, loanID)
	if !delinquent {
		t.Errorf("Expected delinquent=true, got false")
	}

	usecase.MakePayment(ctx, loanID, weekly)
	delinquent, _ = usecase.IsDelinquent(ctx, loanID)
	if !delinquent {
		t.Errorf("Expected delinquent=true, got false")
	}

	usecase.MakePayment(ctx, loanID, weekly*2)
	delinquent, _ = usecase.IsDelinquent(ctx, loanID)
	if delinquent {
		t.Errorf("Expected delinquent=false, got true")
	}
}

func TestMakePaymentInvalidAmount(t *testing.T) {
	ctx := context.Background()
	repo := &mockRepo{}
	usecase := loan.NewLoanUsecase(repo)

	loanID := "loan4"
	principal := int64(5000000)
	usecase.CreateLoan(ctx, loanID, principal)

	weekly := int64(float64(principal)*1.1) / 50
	err := usecase.MakePayment(ctx, loanID, weekly+1)
	if err == nil {
		t.Errorf("Expected error for invalid payment amount, got nil")
	}
}
