package app

import (
	"banking/service"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type Customer struct {
	Name    string `json:"full_name"`
	City    string `json:"city"`
	Zipcode string `json:"zipcode"`
}

type CustomerHandlers struct {
	service service.CustomerService
}

type Game struct {
	Name  string  `xml:"name"`
	Price float32 `xml:"price"`
}

func greet(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "Hello World!")
}

func (ch *CustomerHandlers) getAllCustomer(w http.ResponseWriter, r *http.Request) {
	/*	customers := []Customer{
		{Name: "Ashish", City: "New Delhi", Zipcode: "110075"},
		{Name: "Rob", City: "New Delhi", Zipcode: "110075"},
		{Name: "Thomas", City: "New Delhi", Zipcode: "110075"},
	}*/

	customers, _ := ch.service.GetAllCustomer()

	if r.Header.Get("Content-Type") == "application/xml" {
		w.Header().Add("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(customers)
	} else {
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(customers)
	}
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprint(w, vars["customer_id"])
}

func createCustomer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Post request received")
}

func getAllGames(w http.ResponseWriter, r *http.Request) {
	games := []Game{
		{Name: "Final Fantasy X", Price: float32(32)},
		{Name: "Yakuza", Price: float32(15)},
		{Name: "Baldur's Gate 3", Price: float32(24)},
	}

	if r.Header.Get("COntent-Type") == "application/xml" {
		w.Header().Add("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(games)
	} else {
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(games)
	}
}
