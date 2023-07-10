package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

//go:embed favicon.ico
var faviconFile embed.FS

func handleRoot(w http.ResponseWriter, r *http.Request) {

	serverMetadata := map[string]interface{}{
		"ServerName": r.Host,
		"Method":     r.Method,
		"Path":       r.URL.Path,
	}

	headers := make(map[string]interface{})
	for key, values := range r.Header {
		if len(values) == 1 {
			headers[key] = values[0]
		} else {
			headers[key] = values
		}
	}

	queryParams := make(map[string]interface{})
	queryValues := r.URL.Query()
	for key, values := range queryValues {
		if len(values) == 1 {
			queryParams[key] = values[0]
		} else {
			queryParams[key] = values
		}
	}

	postParams := make(map[string]interface{})
	err := r.ParseForm()
	if err == nil {
		for key, values := range r.Form {
			if len(values) == 1 {
				postParams[key] = values[0]
			} else {
				postParams[key] = values
			}
		}
	}

	response := map[string]interface{}{
		"ServerMetadata": serverMetadata,
		"Headers":        headers,
		"QueryParams":    queryParams,
		"PostParams":     postParams,
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func handlerFavicon(w http.ResponseWriter, r *http.Request) {
	faviconData, err := faviconFile.ReadFile("favicon.ico")
	if err != nil {
		http.Error(w, "Favicon not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "image/x-icon")
	_, _ = w.Write(faviconData)
}

func main() {
	port := 8080
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/favicon.ico", handlerFavicon)
	fmt.Printf("Starting server on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
