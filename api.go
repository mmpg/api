package api

import (
	"log"
	"net/http"

	"github.com/rs/cors"

	"github.com/mmpg/api/endpoints"
	"github.com/mmpg/api/hub"
	"github.com/mmpg/api/notifier"
)

type maxBytesHandler struct {
	h http.Handler
	n int64
}

func (h *maxBytesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, h.n)
	h.h.ServeHTTP(w, r)
}

// Run the MMPG API server
func Run(uv endpoints.UserValidator) {
	go hub.Run()
	go notifier.Run()

	mux := http.NewServeMux()
	mux.HandleFunc("/events", endpoints.Events)
	mux.HandleFunc("/log", endpoints.Log)
	mux.HandleFunc("/auth", endpoints.Auth(uv))
	mux.HandleFunc("/player", endpoints.Player)
	mux.HandleFunc("/world", endpoints.World)

	// TODO: Configure CORS origins properly
	h := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"Accept", "Content-Type", "Authorization"},
	}).Handler(mux)

	log.Fatal(http.ListenAndServe(":8080", &maxBytesHandler{
		h: h,
		n: 1024 * 200,
	}))
}
