package hackaton

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	apiKey    string = "JAA5uTwfT7ThDYxrXFJDVeE79WptQ6FS"
	apiSecret string = "GdnrCuu8RL9PQ2Vu"
	token     string = ""
	latitude  float64
	longitude float64
)

func UpdateCoordinates(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var coords struct {
		UserId    int     `json:"user_id"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	err := json.NewDecoder(r.Body).Decode(&coords)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	latitude = coords.Latitude
	longitude = coords.Longitude

	db := InitDB()
	defer db.Close()
	
	err = InsertHistory(db, coords.UserId, latitude, longitude)
	if err != nil {
		http.Error(w, "Error saving coordinates to history", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}


func GetToken() error {
	tokenURL := "https://test.api.amadeus.com/v1/security/oauth2/token"
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", apiKey)
	data.Set("client_secret", apiSecret)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var jsonResponse map[string]interface{}
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		return err
	}

	token = jsonResponse["access_token"].(string)
	return nil
}

func Activites(w http.ResponseWriter, r *http.Request) {
	if token == "" {
		if err := GetToken(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	url := fmt.Sprintf("https://test.api.amadeus.com/v1/shopping/activities?latitude=" + strconv.FormatFloat(latitude, 'f', -1, 64) + "&longitude=" + strconv.FormatFloat(longitude, 'f', -1, 64) + "&radius=10")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var jsonResponse map[string]interface{}
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jsonResponse)
}

func Hotel(w http.ResponseWriter, r *http.Request) {
	if token == "" {
		if err := GetToken(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	url := fmt.Sprintf("https://test.api.amadeus.com/v1/reference-data/locations/hotels/by-geocode?latitude=" + strconv.FormatFloat(latitude, 'f', -1, 64) + "&longitude=" + strconv.FormatFloat(longitude, 'f', -1, 64) + "&radius=15&radiusUnit=KM&hotelSource=ALL")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var jsonResponse map[string]interface{}
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jsonResponse)
}

func init() {
	go func() {
		for {
			GetToken()
			time.Sleep(30 * time.Minute)
		}
	}()
}
