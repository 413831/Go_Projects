package server

import (
	"fmt"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Method not allowed")
		return
	}
	fmt.Fprint(w, "Hello there %s", "visitor")
}

func getCountries(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%v", countries)
}

func addCountry(w http.ResponseWriter, r *http.Request) {

}
