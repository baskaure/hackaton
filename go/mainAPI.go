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
	apiKey    string = "A3FGwGtaXvep9RKix1vLhjrNxbzWZ1In"
	apiSecret string = "IQ0twTTC58iIgUYL"
	token     string = ""
	latitude  float64
	longitude float64
	radius    int = 10
)

func updateRadiusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		Radius int `json:"radius"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	radius = data.Radius
	fmt.Printf("Updated radius: %d KM\n", radius)
	w.WriteHeader(http.StatusOK)
}

func UpdateCoordinates(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	cookie, errcoo := r.Cookie("user_id")
	if errcoo != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	var val int
	val, _ = strconv.Atoi(cookie.Value)

	var coords struct {
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

	err = InsertHistory(db, val, coords.Latitude, coords.Longitude)
	if err != nil {
		http.Error(w, "Failed to insert history", http.StatusInternalServerError)
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

	url := fmt.Sprintf("https://test.api.amadeus.com/v1/shopping/activities?latitude=" + strconv.FormatFloat(latitude, 'f', -1, 64) + "&longitude=" + strconv.FormatFloat(longitude, 'f', -1, 64) + "&radius=" + strconv.Itoa(radius))
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

	url := fmt.Sprintf("https://test.api.amadeus.com/v1/reference-data/locations/hotels/by-geocode?latitude=" + strconv.FormatFloat(latitude, 'f', -1, 64) + "&longitude=" + strconv.FormatFloat(longitude, 'f', -1, 64) + "&radius=" + strconv.Itoa(radius) + "&radiusUnit=KM&hotelSource=ALL")
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
