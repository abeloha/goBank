package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading.env file: ", err)
	}

	dbname := os.Getenv("DB_NAME")
	dbuser := os.Getenv("DB_USER")
	dbpassword := os.Getenv("DB_PASSWORD")

	println("DB_NAME: ", dbname, "DB_USER: ", dbuser, "DB_PASSWORD", dbpassword)

	storage, err := NewPostgreStore(dbuser, dbname, dbpassword)
	if err != nil {
		panic(err)
	}

	if err := storage.Init(); err != nil {

		log.Fatal(err)
	}
	server := NewAPIServer(":8080", storage)

	server.Run()
}
