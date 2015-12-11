package endpoints

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/mmpg/api/engine"
)

// Log endpoint
func Log(w http.ResponseWriter, r *http.Request) {
	t := r.URL.Query().Get("time")
	e := r.Header.Get("Accept-Encoding")

	if t == "" || !strings.Contains(e, "gzip") {
		w.WriteHeader(400)
		return
	}

	res, err := engine.Log(t)

	if err == engine.ErrConnectionFailed || err == engine.ErrInvalidBase64Encoding {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	b := bytes.NewBuffer(res)

	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Length", strconv.Itoa(b.Len()))
	io.Copy(w, b)
}
