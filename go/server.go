package hackaton

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

func Login(w http.ResponseWriter, r *http.Request) {
	db := InitDB()
	defer db.Close()

	switch r.Method {
	case "GET":
		tmpl, err := template.ParseFiles("html/login.html")
		if err != nil {
			http.Error(w, "Error loading login page", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	case "POST":
		username := r.FormValue("username")
		password := r.FormValue("password")

		authenticated, err := AuthenticateUser(db, username, password)
		if err != nil {
			http.Error(w, "Failed to authenticate user", http.StatusInternalServerError)
			log.Printf("Error authenticating user: %v", err)
			return
		}

		if authenticated {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func GetHistory(w http.ResponseWriter, r *http.Request) {
	db := InitDB()
	defer db.Close()

	rows, err := db.Query("SELECT user_id, latitude, longitude, timestamp FROM HISTORY")
	if err != nil {
		http.Error(w, "Impossible de récupérer l'historique", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type HistoryItem struct {
		UserID    int     `json:"user_id"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Timestamp string  `json:"timestamp"`
	}
	var history []HistoryItem

	for rows.Next() {
		var item HistoryItem
		err := rows.Scan(&item.UserID, &item.Latitude, &item.Longitude, &item.Timestamp)
		if err != nil {
			http.Error(w, "Erreur lors de la lecture de l'historique", http.StatusInternalServerError)
			return
		}
		history = append(history, item)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Erreur lors de la lecture de l'historique", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(history)
	if err != nil {
		http.Error(w, "Erreur lors de l'encodage de l'historique en JSON", http.StatusInternalServerError)
		return
	}
}

func Server() {
	db := InitDB()
	defer db.Close()

	http.HandleFunc("/Signup", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			tmpl, err := template.ParseFiles("html/signup.html")
			if err != nil {
				http.Error(w, "Error loading signup page", http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, nil)
		case "POST":
			username := r.FormValue("newUsername")
			email := r.FormValue("newEmail")
			password := r.FormValue("newPassword")
			if err := InsertUser(db, username, email, password); err != nil {
				http.Error(w, "Failed to insert user", http.StatusInternalServerError)
				log.Printf("Error inserting user: %v", err)
				return
			}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/loginn", Loginn)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/update-coordinates", UpdateCoordinates)
	http.HandleFunc("/activitie", Activites)
	http.HandleFunc("/hotel", Hotel)
	http.HandleFunc("/history", GetHistory)

	publicDir := filepath.Join("public")

	fs := http.FileServer(http.Dir(publicDir))
	http.Handle("/", fs)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server running at http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
