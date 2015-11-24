package endpoints

import (
	"encoding/json"
	"net/http"
)

type UserValidator func(string, string) bool

type authMessage struct {
	Email    string `json:email`
	Password string `json:password`
}

func Auth(uv UserValidator) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var m authMessage

		if err := decoder.Decode(&m); err != nil {
			w.WriteHeader(403)
			return
		}

		if uv(m.Email, m.Password) {
			w.Write([]byte("{}"))
		} else {
			w.WriteHeader(403)
		}
	}
}
