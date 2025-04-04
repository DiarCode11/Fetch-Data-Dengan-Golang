package main

import (
	"fmt"
	"os"
	"encoding/json"
	"net/http"
	"context"
	"time"
	"github.com/joho/godotenv"
	"test-golang/config"
)

type Identity struct {
	Id string `json:id`
	Name string `json:"name"`
	Email string `json:"email"`
}

type Response struct {
	Message string `json:"message"`
	Identity []Identity `json:"identity"`
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	db_name := os.Getenv("DB_NAME")
	fmt.Println("DB_NAME:", db_name)

	var user Identity
	// Ubah request menjadi struct user
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// 🔥 Cetak JSON yang diterima dari user
	fmt.Printf("Received JSON: %+v\n", user)

	// Kirimkan response sukses
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"status": "success",
		"name": user.Name,
		"email": user.Email,
	}
	json.NewEncoder(w).Encode(response)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			rows, err := config.DB.Query(ctx, "SELECT * FROM users")
			if err != nil {
				http.Error(w, "Failed to query DB", http.StatusInternalServerError)
				fmt.Println("DB error:", err)
				return
			}
			defer rows.Close()

			var users []Identity

			for rows.Next() {
				var user Identity
				err := rows.Scan(&user.Id, &user.Name, &user.Email)
				if err != nil {
					http.Error(w, "Failed to read DB result", http.StatusInternalServerError)
					return
				}
				users = append(users, user)
				fmt.Print(user)
			}
			if len(users) == 0 {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}

			response := Response{
				Message: "Hello world",
				Identity: users,
			}
		
			json.NewEncoder(w).Encode(response)
		default:
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func main() {
	config.ConnectDB()
	defer config.CloseDB()

	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/login", postHandler)

	fmt.Println("Server running on port 8080 ...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server failed:", err)
	}
}