package main

import (
	"log"
	"net/http"
)

func main() {

	// Initialize the server and set up the routes
	mux := http.NewServeMux()

	mux.HandleFunc("/", handleHome)
  mux.HandleFunc("/menu", handleMenu)  
  mux.HandleFunc("/reviews", handleReviews)
  mux.HandleFunc("GET /review", handleReviewForm)
  mux.HandleFunc("POST /submitReview", handleReviewSubmission)




	// Serve static files
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Println("Starting the server at http://localhost:4001")

	// Start the server
  err := http.ListenAndServe(":4001", mux)
	if err != nil {
		log.Fatal("Error occurred while starting the server:", err)  // Log the error message
	}
}
