package main

import (
	"log"
	"net/http"

	forum "stepik.leoscode.http/internal/gen/api"
	"stepik.leoscode.http/internal/http/handlers"
	"stepik.leoscode.http/internal/repository"
	"stepik.leoscode.http/internal/service"
)

func main() {

	repo := repository.NewRepoInMemory()
	service2 := service.NewService(repo)
	server := handlers.NewServer(service2)
	handler := forum.Handler(server)
	if err := http.ListenAndServe("localhost:8080", handler); err != nil {
		log.Fatal(err)
	}
}
