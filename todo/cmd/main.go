package main

import (
	"log"
	"os"
	"strconv"
	"todo/handler"
	"todo/storage"
)

func main() {
	port, _ :=strconv.Atoi(os.Getenv("DB_PORT"))
	var temp = storage.Options{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}
	db, err := storage.New(temp)
	if err != nil {
		log.Fatalf("storage failed: %s", err.Error())
	}
	log.Println("storage: OK")
	defer db.Close()

	h := handler.New(db)
	h.Run()
}
