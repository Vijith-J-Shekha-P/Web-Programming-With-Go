package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type MenuItem struct {
	Name  string `json:"name"`
	Price string `json:"price"`
}

type Review struct {
	Name     string `json:"name"`
	Dish     string `json:"dish"`
	Rating   int    `json:"rating"`
	Comments string `json:"comments"`
}

var menuItems = []MenuItem{
	{Name: "Spaghetti Carbonara", Price: "$12"},
	{Name: "Margherita Pizza", Price: "$10"},
	{Name: "Caesar Salad", Price: "$8"},
}

var reviews = []Review{
	{Name: "Abhisekh Kumar", Dish: "Spaghetti Carbonara", Rating: 5, Comments: "Absolutely delicious! The best Carbonara I've ever had."},
	{Name: "Ritu Gupta", Dish: "Margherita Pizza", Rating: 4, Comments: "Very tasty and fresh. Would have liked a bit more basil."},
}

func main() {
	http.HandleFunc("GET /data", handleData)
	http.HandleFunc("GET /reviews", handleGetReviews)
	http.HandleFunc("POST /addReview", handlePostReview)

	log.Println("Starting server to handle requests")
	err := http.ListenAndServe(":4002", nil)
	if err != nil {
		log.Fatal("Error starting request server: ", err)
	}
	log.Println("Request server stopped")
}

func handleData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(menuItems)
	if err != nil {
		http.Error(w, "Error encoding menu items", http.StatusInternalServerError)
	}
}

func handleGetReviews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(reviews)
	if err != nil {
		http.Error(w, "Error encoding reviews", http.StatusInternalServerError)
	}
}

func handlePostReview(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var review Review
	err = json.Unmarshal(body, &review)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}
	reviews = append(reviews, review)

	w.WriteHeader(http.StatusCreated)
}
