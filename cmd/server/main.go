package main

import (
	"log"
	"net/http"

	forum "stepik.leoscode.http/internal/gen/api"
	"stepik.leoscode.http/internal/http/handlers"
	"stepik.leoscode.http/internal/repository/inMemoryRepo"
	"stepik.leoscode.http/internal/service"
)

func main() {

	repo := inMemoryRepo.NewRepoInMemory()
	service2 := service.NewService(repo)
	server := handlers.NewServer(service2)
	//handler := forum.Handler(server)
	handler := forum.HandlerWithOptions(server, forum.StdHTTPServerOptions{
		ErrorHandlerFunc: handlers.ApiErrorHandler,
	})
	if err := http.ListenAndServe("localhost:8080", handler); err != nil {
		log.Fatal(err)
	}

}
