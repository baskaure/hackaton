package hackaton

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func Server() {
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
