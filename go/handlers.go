package hackaton

import (
	"encoding/json"
	"log"
	"net/http"
	"text/template"
)

func Loginn(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("html/login.html")
	if err != nil {
		http.Error(w, "Error loading login page", http.StatusInternalServerError)
		log.Println("Error parsing template:", err)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		log.Println("Error executing template:", err)
		return
	}
}

func Historique(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("html/historique.html")
	if err != nil {
		http.Error(w, "Error loading login page", http.StatusInternalServerError)
		log.Println("Error parsing template:", err)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		log.Println("Error executing template:", err)
		return
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
