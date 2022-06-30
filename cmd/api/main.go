package main

import (
	"audioShareApp/pkg/api"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

// Add is our function that sums two integers
func Add(x, y int) (res int) {
	return x + y
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", api.ListAudios)
	r.Post("/audio", api.UploadAudio)

	fmt.Println("✅ Server up and running on port 8090")
	log.Fatalln(http.ListenAndServe(":8090", r))
}
