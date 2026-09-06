// ...existing code...
package main

import (
	"Users/harshalgosavi/workspace/GO/model"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type ApiResponse struct {
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	TimeStamp string      `json:"timestamp"`
}

func getData(w http.ResponseWriter, r *http.Request) {

	u := model.User{ID: 1, Name: "Alice", Email: "alice@example.com"}
	w.Header().Set("Content-Type", "application/json")
	respi := ApiResponse{
		Message:   "Data fetched successfully",
		Data:      u,
		TimeStamp: "2024-10-10T10:00:00Z",
	}
	_ = json.NewEncoder(w).Encode(respi)

}

func writeData(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	fmt.Println("body with string: " + string(body) + " body with string: " + string(body))
	defer r.Body.Close()
	var u model.User
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("Data written successfully"))
}


func main6() {
	http.HandleFunc("/", getData)
	http.HandleFunc("/write", writeData)

	log.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
