package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type URLShortener struct {
	store       map[string]string
	domainCount map[string]int
	mu          sync.RWMutex
}

func NewURLShorterner() *URLShortener {
	return &URLShortener{
		store:       make(map[string]string),
		domainCount: make(map[string]int),
	}
}

// Genartes the shortURL using some random Alpha-Numberic value
func generateShortURL() string {
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 7)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// Extracts the domain form the url
func extarctDomain(url string) string {

	//removes // from http://www.google.com/
	parts := strings.Split(url, "/")

	if len(parts) < 3 {
		return ""
	}
	domain := parts[2]

	//removes the 'www'
	domain = strings.TrimPrefix(domain, "www.")

	domainParts := strings.Split(domain, ".")

	if len(domainParts) >= 2 {
		return domainParts[len(domainParts)-2]
	}
	return domain
}

// Handles  the URL shortening request
func (us *URLShortener) shortenURL(w http.ResponseWriter, r *http.Request) {
	var request struct {
		URL string `json:"url"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	//check the existance of the URL
	us.mu.RLock()
	for short, original := range us.store {
		if original == request.URL {
			us.mu.RUnlock()
			response := map[string]string{"short_url": short}
			json.NewEncoder(w).Encode(response)
			return
		}
	}
	us.mu.RUnlock()

	//generates the new shortURL
	shortURL := generateShortURL()

	//Extracts the domain
	domain := extarctDomain(request.URL)
	us.mu.Lock()
	us.store[shortURL] = request.URL
	us.domainCount[domain]++
	us.mu.Unlock()

	response := map[string]string{"short_url": shortURL}
	json.NewEncoder(w).Encode(response)
}

// Function to redirect the shortenURL
func (us *URLShortener) redirectURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	shortURL := vars["shortURL"]

	us.mu.RLock()
	originalURL, exists := us.store[shortURL]
	us.mu.RUnlock()

	if exists {
		http.Redirect(w, r, originalURL, http.StatusFound)
	} else {
		http.Error(w, "URL not found", http.StatusNotFound)
	}
}

// Function to handle the top 3 domains
func (us *URLShortener) metrics(w http.ResponseWriter, r *http.Request) {
	type domainEntry struct {
		Domain string
		Count  int
	}

	us.mu.RLock()

	defer us.mu.RUnlock()

	var domains []domainEntry
	for domain, count := range us.domainCount {
		domains = append(domains, domainEntry{Domain: domain, Count: count})

	}

	for i := 0; i < len(domains); i++ {
		for j := i + 1; j < len(domains); j++ {
			if domains[i].Count < domains[j].Count {
				domains[i], domains[j] = domains[j], domains[i]
			}
		}
	}
	topDomains := domains

	if len(domains) > 3 {
		topDomains = domains[:3]
	}
	response := ""
	for _, entry := range topDomains {
		response += fmt.Sprintf("%s: %d, ", entry.Domain, entry.Count)
		fmt.Sprintln()
	}

	// Return the JSON response
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, response)
}

// function to display the list of shortURLs and original URLS in the memory
func (us *URLShortener) listURLs(w http.ResponseWriter, r *http.Request) {
	us.mu.RLock()
	defer us.mu.RUnlock()

	// Create a response map with short_url -> original_url mappings
	response := make(map[string]string)
	for shortURL, originalURL := range us.store {
		response[shortURL] = originalURL
	}

	// Return the JSON response
	json.NewEncoder(w).Encode(response)
}
func main() {
	urlShortener := NewURLShorterner()

	r := mux.NewRouter()

	r.HandleFunc("/shorten", urlShortener.shortenURL).Methods("POST")
	r.HandleFunc("/metrics", urlShortener.metrics).Methods("GET")
	r.HandleFunc("/{shortURL}", urlShortener.redirectURL).Methods("GET")
	r.HandleFunc("/metrics/list", urlShortener.listURLs).Methods("GET")

	log.Println("Starting Server on Port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
