package http_test

import (
	"SRATest/domain"
	loanHandler "SRATest/server/http"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

// Mock usecase
type mockLoanUsecase struct {
	store map[string]*domain.Loan
}

func (m *mockLoanUsecase) CreateLoan(ctx context.Context, id string, principal int64) (*domain.Loan, error) {
	if m.store == nil {
		m.store = make(map[string]*domain.Loan)
	}
	l := &domain.Loan{
		ID:            id,
		Principal:     principal,
		InterestRate:  0.1,
		TotalWeeks:    50,
		WeeklyPayment: int64(float64(principal) * 1.1 / 50),
		Payments:      make([]domain.Payment, 50),
	}
	for i := 0; i < 50; i++ {
		l.Payments[i] = domain.Payment{Week: i + 1, Paid: false}
	}
	m.store[id] = l
	return l, nil
}

func (m *mockLoanUsecase) MakePayment(ctx context.Context, id string, amount int64) error {
	l, ok := m.store[id]
	if !ok {
		return nil
	}
	weekly := l.WeeklyPayment
	count := int(amount / weekly)
	for i := 0; i < len(l.Payments) && count > 0; i++ {
		if !l.Payments[i].Paid {
			l.Payments[i].Paid = true
			count--
		}
	}
	return nil
}

func (m *mockLoanUsecase) GetOutstanding(ctx context.Context, id string) (int64, error) {
	l, ok := m.store[id]
	if !ok {
		return 0, nil
	}
	unpaid := 0
	for _, p := range l.Payments {
		if !p.Paid {
			unpaid++
		}
	}
	return int64(unpaid) * l.WeeklyPayment, nil
}

func (m *mockLoanUsecase) IsDelinquent(ctx context.Context, id string) (bool, error) {
	l, ok := m.store[id]
	if !ok {
		return false, nil
	}
	missed := 0
	for _, p := range l.Payments {
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

func TestCreateLoanHandler(t *testing.T) {
	usecase := &mockLoanUsecase{}
	h := loanHandler.NewHandler(usecase)

	payload := []byte(`{"id":"loan1","principal":5000000}`)
	req := httptest.NewRequest(http.MethodPost, "/loan", bytes.NewBuffer(payload))
	w := httptest.NewRecorder()

	h.CreateLoan(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["message"] != "Loan created successfully" {
		t.Errorf("Unexpected message: %v", resp["message"])
	}
}

func TestMakePaymentHandler(t *testing.T) {
	usecase := &mockLoanUsecase{store: make(map[string]*domain.Loan)}
	h := loanHandler.NewHandler(usecase)
	usecase.CreateLoan(context.Background(), "loan2", 5000000)

	payload := []byte(`{"amount":110000}`)
	req := httptest.NewRequest(http.MethodPost, "/loan/loan2/pay", bytes.NewBuffer(payload))
	req = mux.SetURLVars(req, map[string]string{"id": "loan2"})
	w := httptest.NewRecorder()

	h.MakePayment(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
}

func TestGetOutstandingHandler(t *testing.T) {
	usecase := &mockLoanUsecase{store: make(map[string]*domain.Loan)}
	h := loanHandler.NewHandler(usecase)
	usecase.CreateLoan(context.Background(), "loan3", 5000000)

	req := httptest.NewRequest(http.MethodGet, "/loan/loan3/outstanding", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "loan3"})
	w := httptest.NewRecorder()

	h.GetOutstanding(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
}

func TestIsDelinquentHandler(t *testing.T) {
	usecase := &mockLoanUsecase{store: make(map[string]*domain.Loan)}
	h := loanHandler.NewHandler(usecase)
	usecase.CreateLoan(context.Background(), "loan4", 5000000)

	req := httptest.NewRequest(http.MethodGet, "/loan/loan4/delinquent", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "loan4"})
	w := httptest.NewRecorder()

	h.IsDelinquent(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
}
