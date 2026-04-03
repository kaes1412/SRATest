package infrastructure

import (
	"SRATest/domain"
	"context"
	"fmt"
	"sync"
)

type InMemoryLoanRepository struct {
	data map[string]*domain.Loan
	mu   sync.RWMutex
}

func NewInMemoryLoanRepository() *InMemoryLoanRepository {
	return &InMemoryLoanRepository{
		data: make(map[string]*domain.Loan),
	}
}

func (r *InMemoryLoanRepository) Save(ctx context.Context, loan *domain.Loan) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[loan.ID] = loan
	return nil
}

func (r *InMemoryLoanRepository) FindByID(ctx context.Context, id string) (*domain.Loan, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	loan, ok := r.data[id]
	if !ok {
		return nil, fmt.Errorf("loan not found")
	}
	return loan, nil
}

func (r *InMemoryLoanRepository) Update(ctx context.Context, loan *domain.Loan) error {
	if _, ok := r.data[loan.ID]; !ok {
		return fmt.Errorf("loan not found")
	}
	r.data[loan.ID] = loan
	return nil
}
