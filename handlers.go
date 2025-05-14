package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "Error loading home page", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func handleMenu(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Make GET request to data server
	resp, err := http.Get("http://localhost:4002/data")
	if err != nil {
		http.Error(w, "Error fetching menu data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading menu items", http.StatusInternalServerError)
		return
	}

	// Unmarshal JSON data
	var menuItems []MenuItem
	err = json.Unmarshal(body, &menuItems)
	if err != nil {
		http.Error(w, "Error decoding menu items", http.StatusInternalServerError)
		return
	}

	// Parse template and render menu items
	tmpl, err := template.ParseFiles("templates/menu.html")
	if err != nil {
		http.Error(w, "Error loading menu page", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, menuItems)
}

func handleReviewForm(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/review_form.html")
	if err != nil {
		http.Error(w, "Error loading review form", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func handleReviewSubmission(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	// Retrieve and convert form values
	name := r.FormValue("name")
	dish := r.FormValue("dish")
	rating := stringToInt(r.FormValue("rating"))
	comments := r.FormValue("comments")

	// Create review object
	review := Review{
		Name:     name,
		Dish:     dish,
		Rating:   rating,
		Comments: comments,
	}

	// Marshal review to JSON
	reviewData, err := json.Marshal(review)
	if err != nil {
		http.Error(w, "Error encoding review data", http.StatusInternalServerError)
		return
	}

	// POST the review to data server
	resp, err := http.Post("http://localhost:4002/addReview", "application/json", bytes.NewBuffer(reviewData))
	if err != nil {
		http.Error(w, "Error posting review", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Redirect to /reviews after successful submission
	http.Redirect(w, r, "/reviews", http.StatusSeeOther)
}


func handleReviews(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Step 1: Make GET request to the data server
	resp, err := http.Get("http://localhost:4002/reviews")
	if err != nil {
		http.Error(w, "Error fetching reviews", http.StatusInternalServerError)
		fmt.Println("Fetch error:", err)
		return
	}
	defer resp.Body.Close() // Step 4: Always close response body

	// Step 2: Check for non-200 status
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to fetch reviews", http.StatusInternalServerError)
		fmt.Println("Non-200 status code:", resp.Status)
		return
	}

	// Step 3: Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading reviews", http.StatusInternalServerError)
		fmt.Println("Read error:", err)
		return
	}

	// Step 5: Decode JSON into Review slice
	var reviews []Review
	err = json.Unmarshal(body, &reviews)
	if err != nil {
		http.Error(w, "Error decoding reviews", http.StatusInternalServerError)
		fmt.Println("Unmarshal error:", err)
		return
	}

	// Step 6: Parse and render the template
	tmpl, err := template.ParseFiles("templates/reviews.html")
	if err != nil {
		http.Error(w, "Error loading reviews page", http.StatusInternalServerError)
		fmt.Println("Template parse error:", err)
		return
	}
	tmpl.Execute(w, reviews)
}



func stringToInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}
