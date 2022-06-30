package main

import (
	"audioShareApp/pkg/api"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
)

// Add is our function that sums two integers
func Add(x, y int) (res int) {
	return x + y
}

func main() {
	port := os.Getenv("PORT")
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Hello world"))
	})
	r.Get("/audio", api.ListAudios)
	r.Post("/audio", api.UploadAudio)

	fmt.Println("✅ Server up and running on port: " + port)
	log.Fatalln(http.ListenAndServe(":"+port, r))
}
