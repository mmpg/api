package endpoints

import (
	"log"
	"net/http"
	"strings"

	"github.com/mmpg/api/engine"
)

// Log endpoint
func Log(w http.ResponseWriter, r *http.Request) {
	t := r.URL.Query().Get("time")

	if t == "" {
		w.WriteHeader(400)
		return
	}

	res, err := engine.Log(t)

	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	if strings.Contains(res, "ERROR") {
		parts := strings.Split(res, " ")
		w.WriteHeader(400)

		if len(parts) > 1 {
			w.Write([]byte(parts[1]))
		}

		return
	}

	w.Write([]byte(res))
}
