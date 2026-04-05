package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type URL struct {
	ShortCode   string     `bson:"shortCode"`
	OriginalURL string     `bson:"originalUrl"`
	CreatedAt   time.Time  `bson:"createdAt"`
	ExpiresAt   *time.Time `bson:"expiresAt,omitempty"`
}

var collection *mongo.Collection

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generateShortCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func parseExpiry(expiry *string) *time.Time {
	if expiry == nil {
		return nil
	}
	t, err := time.Parse(time.RFC3339, *expiry)
	if err != nil {
		return nil
	}
	return &t
}

func shortenURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		URL    string  `json:"url"`
		Expiry *string `json:"expiry,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	if req.URL == "" {
		http.Error(w, "URL is required", 400)
		return
	}

	shortCode := generateShortCode()
	expiresAt := parseExpiry(req.Expiry)

	url := URL{
		ShortCode:   shortCode,
		OriginalURL: req.URL,
		CreatedAt:   time.Now(),
		ExpiresAt:   expiresAt,
	}

	_, err := collection.InsertOne(context.TODO(), url)
	if err != nil {
		log.Printf("Failed to insert URL: %v", err)
		http.Error(w, "Failed to save URL", 500)
		return
	}

	shortURL := "http://localhost:8080/" + shortCode
	json.NewEncoder(w).Encode(map[string]string{"shortUrl": shortURL})
}

func redirectURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["shortCode"]

	var url URL
	err := collection.FindOne(context.TODO(), bson.M{"shortCode": shortCode}).Decode(&url)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.NotFound(w, r)
			return
		}
		log.Printf("Failed to find URL: %v", err)
		http.Error(w, "Internal server error", 500)
		return
	}

	if url.ExpiresAt != nil && time.Now().After(*url.ExpiresAt) {
		http.Error(w, "URL has expired", 410)
		return
	}

	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("urlshortener")
	collection = db.Collection("urls")

	r := mux.NewRouter()
	r.HandleFunc("/", serveIndex).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.HandleFunc("/shorten", shortenURL).Methods("POST")
	r.HandleFunc("/{shortCode}", redirectURL).Methods("GET")

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
