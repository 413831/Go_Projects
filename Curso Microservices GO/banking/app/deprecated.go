package app

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func OldStart() {
	router := mux.NewRouter()

	// define routes
	router.HandleFunc("/greet", greet).Methods(http.MethodGet)
	//router.HandleFunc("/customers", getAllCustomers).Methods(http.MethodGet)
	router.HandleFunc("/customers/{customer_id:[0-9]+}", getCustomer).Methods(http.MethodGet)
	router.HandleFunc("/customers", createCustomer).Methods(http.MethodPost)

	router.HandleFunc("/games", getAllGames).Methods(http.MethodGet)

	log.Println("starting server at port 8000")
	// starting server
	log.Fatal(http.ListenAndServe("localhost:8000", router))
}
