package internal

import (
	internal "learning-go/internal/handlers"
	"net/http"
)

type RouterHttp struct {
	handler *internal.Handler
}

func NewRouterHttp() *RouterHttp {
	handler := internal.NewHandler()

	http.HandleFunc("/", handler.Hello)

	http.HandleFunc("/health", handler.CheckHealth)

	http.HandleFunc("/user", handler.UserHandler)

	return &RouterHttp{handler: handler}
}

func (r RouterHttp) Run(port string) error {
	return http.ListenAndServe(port, nil)
}
