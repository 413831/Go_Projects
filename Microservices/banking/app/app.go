package app

import (
	"banking/domain"
	"banking/service"
	mux "github.com/gorilla/mux"
	"log"
	"net/http"
)

func Start() {
	router := mux.NewRouter()

	// wiring
	ch := CustomerHandlers{service: service.NewCustomerService(domain.NewCustomerRepositoryStub())}

	// define routes
	router.HandleFunc("/customers", ch.getAllCustomer).Methods(http.MethodGet)

	log.Println("starting server at port 8000")
	// starting server
	log.Fatal(http.ListenAndServe("localhost:8000", router))
}
