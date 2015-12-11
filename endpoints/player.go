package endpoints

import (
	"log"
	"net/http"

	"github.com/mmpg/api/engine"
)

// Max player file size: 200 KB
const maxPlayerFileSize = 1024 * 200

// Player endpoint allows deployment of new players
func Player(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}

	token, err := ValidateToken(r)

	if err != nil {
		w.WriteHeader(401)
		return
	}

	err = r.ParseMultipartForm(maxPlayerFileSize)

	if err != nil {
		log.Println(err)
		w.WriteHeader(413)
		return
	}

	player := r.MultipartForm.File["player"]

	if len(player) != 1 {
		w.WriteHeader(400)
		return
	}

	content, err := player[0].Open()
	defer content.Close()

	if err != nil {
		w.WriteHeader(400)
		return
	}

	// TODO: Ensure that content represents a C++ file
	err = engine.DeployPlayer(token.Email, content)

	if handleError(w, err) {
		return
	}
}
