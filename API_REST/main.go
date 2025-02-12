package main

import (
	"fmt"
	"log"
	"net/http"
	"project/router"
)

func main() {
	r := router.SetupRouter()

	fmt.Println("Server running on port 8080")

	log.Fatal(http.ListenAndServe(":8080", r))
}
