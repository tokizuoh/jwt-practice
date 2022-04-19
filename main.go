package main

import (
	"encoding/json"
	"net/http"
)

type Model struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		stubModel := &Model{
			Id:   1,
			Name: "Oka",
		}
		if err := json.NewEncoder(w).Encode(stubModel); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})
	http.ListenAndServe(":8080", nil)
}
