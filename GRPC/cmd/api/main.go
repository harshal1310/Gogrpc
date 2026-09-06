package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"mygrpc/pkg/client"
)

// CreateUserPayload represents the JSON payload for creating a user
type CreateUserPayload struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// GetUserPayload represents the JSON payload for getting a user
type GetUserPayload struct {
	Id string `json:"id"`
}

// ResponsePayload represents a unified response structure
type ResponsePayload struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// CreateUserHandler handles POST /users to create a new user
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode(ResponsePayload{
			Success: false,
			Error:   fmt.Sprintf("failed to read body: %v", err),
		})
		return
	}
	defer r.Body.Close()

	var payload CreateUserPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponsePayload{
			Success: false,
			Error:   fmt.Sprintf("invalid JSON: %v", err),
		})
		return
	}

	if payload.Name == "" || payload.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponsePayload{
			Success: false,
			Error:   "name and email are required",
		})
		return
	}

	// Call the gRPC client
	res, err := client.CreateUserRequest(payload.Name, payload.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponsePayload{
			Success: false,
			Error:   fmt.Sprintf("failed to create user: %v", err),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ResponsePayload{
		Success: true,
		Data: map[string]string{
			"id":    res.GetId(),
			"name":  res.GetName(),
			"email": res.GetEmail(),
		},
	})
}

// GetUserHandler handles GET /users?id=X or POST /users with JSON body
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var id string

	// For GET requests, read from query parameter
	if r.Method == http.MethodGet {
		id = r.URL.Query().Get("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ResponsePayload{
				Success: false,
				Error:   "id query parameter is required",
			})
			return
		}
	} else if r.Method == http.MethodPost {
		// For POST requests, read from JSON body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			json.NewEncoder(w).Encode(ResponsePayload{
				Success: false,
				Error:   fmt.Sprintf("failed to read body: %v", err),
			})
			return
		}
		defer r.Body.Close()

		var payload GetUserPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ResponsePayload{
				Success: false,
				Error:   fmt.Sprintf("invalid JSON: %v", err),
			})
			return
		}

		id = payload.Id
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ResponsePayload{
				Success: false,
				Error:   "id is required",
			})
			return
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Call the gRPC client
	res, err := client.GetUserRequest(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ResponsePayload{
			Success: false,
			Error:   fmt.Sprintf("failed to get user: %v", err),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponsePayload{
		Success: true,
		Data: map[string]string{
			"id":    res.GetId(),
			"name":  res.GetName(),
			"email": res.GetEmail(),
		},
	})
}

// HealthHandler handles GET /health
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func main() {
	http.HandleFunc("/health", HealthHandler)
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			// Check if body has id field (GET user) or name/email (CREATE user)
			body, _ := io.ReadAll(r.Body)
			defer r.Body.Close()

			var data map[string]interface{}
			json.Unmarshal(body, &data)

			// Reset body for handlers
			r.Body = io.NopCloser(bytes.NewReader(body))

			if _, hasId := data["id"]; hasId && len(data) == 1 {
				GetUserHandler(w, r)
			} else {
				CreateUserHandler(w, r)
			}
		} else if r.Method == http.MethodGet {
			GetUserHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	port := ":8080"
	fmt.Printf("REST API server starting on %s\n", port)
	fmt.Println("Endpoints:")
	fmt.Println("  POST /users - Create user (JSON: {\"name\": \"...\", \"email\": \"...\"})")
	fmt.Println("  GET /users?id=<id> - Get user by ID")
	fmt.Println("  POST /users - Get user by ID (JSON: {\"id\": \"...\"})")
	fmt.Println("  GET /health - Health check")

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
