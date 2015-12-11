package endpoints

import (
	"log"
	"net/http"

	"github.com/mmpg/api/engine"
)

func handleError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}

	log.Println(err)

	switch err {
	case engine.ErrConnectionFailed, engine.ErrInvalidBase64Encoding:
		w.WriteHeader(500)

	default:
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}

	return true
}
