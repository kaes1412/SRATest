package http

import (
	"github.com/gorilla/mux"
)

func NewRouter(h *Handler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/loan", h.CreateLoan).Methods("POST")
	r.HandleFunc("/loan/{id}/pay", h.MakePayment).Methods("POST")
	r.HandleFunc("/loan/{id}/outstanding", h.GetOutstanding).Methods("GET")
	r.HandleFunc("/loan/{id}/delinquent", h.IsDelinquent).Methods("GET")

	return r
}
