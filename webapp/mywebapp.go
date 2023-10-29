package webapp

import (
	"fmt"
	"net/http"
)

func StartServer() {
	// Create a new http.ServeMux to handle routes
	// mux is multiplexer, router, handler
	mux := http.NewServeMux()

	// Define your routes and handlers
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the homepage!")
	})

	mux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			fmt.Fprintf(w, "This is the About page. This response is for GET requests only.")
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Create an http.Server instance with custom settings
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Start the server
	fmt.Println("Server is running on port 8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Server error:", err)
	}
}
