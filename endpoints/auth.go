package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/mmpg/api/engine"
)

// UserValidator tells whether a given email and password are valid
type UserValidator func(string, string) bool

type authMessage struct {
	Email    string `json:email`
	Password string `json:password`
}

// Auth handles authentication
func Auth(uv UserValidator) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var m authMessage

		if err := decoder.Decode(&m); err != nil {
			w.WriteHeader(403)
			return
		}

		if !uv(m.Email, m.Password) {
			w.WriteHeader(403)
			return
		}

		res, _, _ := engine.PlayerExists(m.Email)

		if res != "TRUE" {
			w.WriteHeader(400)
			return
		}

		w.Write([]byte("{}"))
	}
}
