package main

import (
	"SRATest/infrastructure"
	server "SRATest/server/http"
	usecase "SRATest/usecase/loan"
	"log"
	"net/http"
)

func main() {
	// init repo
	repo := infrastructure.NewInMemoryLoanRepository()

	// init usecase
	uc := usecase.NewLoanUsecase(repo)

	// init handler
	handler := server.NewHandler(uc)

	// init router
	router := server.NewRouter(handler)

	log.Println("server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
