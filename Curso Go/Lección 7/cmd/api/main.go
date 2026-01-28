package main

import (
	handlers "learning-go/internal/handlers"
	repos "learning-go/internal/repos"
	"net/http"
)

func main() {
	/*r := internal.NewRouterHttp()

	fmt.Println("Listening on port 8080")

	if err := r.Run(":8080"); err != nil {
		fmt.Errorf("error running server %s", err)
	}*/

	repo := repos.NewMockRepo()

	handler := handlers.NewUserHandler(repo)
	http.HandleFunc("/users", handler.GetUser)

	http.ListenAndServe(":8080", nil)
}
