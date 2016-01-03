package endpoints

import (
	"net/http"

	"github.com/mmpg/api/engine"
)

// World endpoint
func World(w http.ResponseWriter, r *http.Request) {
	res, err := engine.World()

	if handleError(w, err) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(res))
}
