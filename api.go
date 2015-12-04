package api

import (
	"log"
	"net/http"

	"github.com/rs/cors"

	"github.com/mmpg/api/endpoints"
	"github.com/mmpg/api/hub"
	"github.com/mmpg/api/notifier"
)

// Run the MMPG API server
func Run(uv endpoints.UserValidator) {
	go hub.Run()
	go notifier.Run()

	mux := http.NewServeMux()
	mux.HandleFunc("/events", endpoints.Events)
	mux.HandleFunc("/log", endpoints.Log)
	mux.HandleFunc("/auth", endpoints.Auth(uv))
	mux.HandleFunc("/player", endpoints.Player)

	// TODO: Configure CORS origins properly
	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"Accept", "Content-Type", "Authorization"},
	}).Handler(mux)

	log.Fatal(http.ListenAndServe(":8080", handler))
}
