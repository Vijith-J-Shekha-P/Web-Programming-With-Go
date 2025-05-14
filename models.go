package main

// MenuItem represents a dish on the restaurant menu
type MenuItem struct {
	Name  string `json:"name"`
	Price string `json:"price"`
}

// Review represents a review submitted by a user
type Review struct {
	Name     string `json:"name"`
	Dish     string `json:"dish"`
	Rating   int    `json:"rating"`
	Comments string `json:"comments"`
}
