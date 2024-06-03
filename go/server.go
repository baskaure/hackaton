package hackaton

import (
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
			http.Redirect(w, r, "/home", http.StatusSeeOther)
		} else {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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
