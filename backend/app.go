package main

import (
	"goth-complete-setup/internal/server"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {

	// Load all environment variables
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	port, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		log.Fatal("APP_PORT env variable not set correctly")
	}

	server := server.NewServer(port)
	log.Printf("Running server on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
