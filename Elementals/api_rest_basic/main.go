package main

import (
	"api_rest_basic/server"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()

	sc := make(chan os.Signal, 1)

	signal.Notify(sc, os.Interrupt, syscall.SIGTERM)

	srv := server.New(":8080")

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()
	log.Println("server started")

	<-sc

	srv.Shutdown(ctx)
	log.Println("server stopped")

}
