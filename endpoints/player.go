package endpoints

import (
	"log"
	"net/http"
)

// Max player file size: 200 KB
const maxPlayerFileSize = 1024 * 200

// Player endpoint allows deployment of new players
func Player(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}

	_, err := ValidateToken(r)

	if err != nil {
		w.WriteHeader(401)
		return
	}

	err = r.ParseMultipartForm(maxPlayerFileSize)

	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
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

	w.Write([]byte("OK!"))
}
