package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
)

type Customer struct {
	Name    string `json:"full_name"`
	City    string `json:"city"`
	Zipcode string `json:"zipcode"`
}

type Game struct {
	Name  string  `xml:"name"`
	Price float32 `xml:"price"`
}

func main() {
	// define routes
	http.HandleFunc("/greet", greet)
	http.HandleFunc("/customers", getAllCustomer)
	http.HandleFunc("/games", getAllGames)

	log.Println("starting server at port 8000")
	// starting server
	log.Fatal(http.ListenAndServe("localhost:8000", nil))

}

func greet(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "Hello World!")
}

func getAllCustomer(w http.ResponseWriter, _ *http.Request) {
	customers := []Customer{
		{Name: "Ashish", City: "New Delhi", Zipcode: "110075"},
		{Name: "Rob", City: "New Delhi", Zipcode: "110075"},
		{Name: "Thomas", City: "New Delhi", Zipcode: "110075"},
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
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
