package main

import (
	"api_rest_basic/server"
	"fmt"
)

func main() {

	srv := server.New(":8080")
	err := srv.ListenAndServe()

	if err != nil {
		fmt.Println("Error:" + err.Error())
	}
}
