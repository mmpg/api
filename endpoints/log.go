package endpoints

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/mmpg/api/engine"
)

// Log endpoint
func Log(w http.ResponseWriter, r *http.Request) {
	t := r.URL.Query().Get("time")

	if t == "" {
		w.WriteHeader(400)
		return
	}

	res, connErr, err := engine.Log(t)

	if connErr != nil {
		log.Println(connErr)
		w.WriteHeader(500)
		return
	}

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	bytes := bytes.NewBufferString(res)
	w.Header().Set("Content-Length", strconv.Itoa(bytes.Len()))
	io.Copy(w, bytes)
}
