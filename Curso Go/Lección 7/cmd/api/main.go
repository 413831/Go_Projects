package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", Hello)

	http.HandleFunc("/health", CheckHealth)

	http.HandleFunc("/user", UserHandler)

	http.ListenAndServe(":8080", nil)

	fmt.Println("Listening on port 8080")

}
