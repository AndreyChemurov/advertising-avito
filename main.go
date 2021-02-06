package main

import (
	"advertising/internal/database"
	internalhttp "advertising/internal/http"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Создание таблиц и индексов
	if err := database.CreateTableAndIndecies(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Обработка запросов
	http.HandleFunc("/create", internalhttp.CreateAdv)
	http.HandleFunc("/getone", internalhttp.GetOneAdv)
	http.HandleFunc("/getall", internalhttp.GetAllAdv)

	http.HandleFunc("/", internalhttp.NotFound)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
