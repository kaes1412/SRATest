package http

import (
	usecase "SRATest/usecase/loan"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Handler struct
type Handler struct {
	usecase usecase.LoanUsecase
}

func NewHandler(u usecase.LoanUsecase) *Handler {
	return &Handler{usecase: u}
}

func (h *Handler) CreateLoan(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		ID        string `json:"id"`
		Principal int64  `json:"principal"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	loan, err := h.usecase.CreateLoan(ctx, req.ID, req.Principal)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	WriteJSON(w, http.StatusOK, "Loan created successfully", loan)
}

func (h *Handler) MakePayment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]

	var req struct {
		Amount int64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSON(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	if err := h.usecase.MakePayment(ctx, id, req.Amount); err != nil {
		WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	WriteJSON(w, http.StatusOK, "payment success", nil)
}

func (h *Handler) GetOutstanding(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]

	val, err := h.usecase.GetOutstanding(ctx, id)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	WriteJSON(w, http.StatusOK, "Outstanding fetched successfully", map[string]int64{
		"outstanding": val,
	})
}

func (h *Handler) IsDelinquent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]

	val, err := h.usecase.IsDelinquent(ctx, id)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	WriteJSON(w, http.StatusOK, "Delinquent status fetched successfully", map[string]bool{
		"delinquent": val,
	})
}
