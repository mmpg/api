package endpoints

import (
	"log"
	"net/http"

	"github.com/mmpg/api/engine"
)

func Log(w http.ResponseWriter, r *http.Request) {
	reply, err := engine.Test()

	if err != nil {
		log.Println(err)
		return
	}

	w.Write([]byte(reply))
}
