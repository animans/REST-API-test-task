package main

import (
	"log"

	"github.com/animans/REST-API-test-task/http"
	"github.com/animans/REST-API-test-task/infastructure"
	"github.com/joho/godotenv"
)

// init ...
func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

// main ...
func main() {
	repo := infastructure.NewServiceRepoPG()
	err := repo.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()
	api := http.NewHandlers(repo)
	if err := api.Start(); err != nil {
		log.Fatal(err)
	}
}
