package endpoints

import (
	"log"
	"net/http"

	"github.com/mmpg/api/client"
	"github.com/mmpg/api/hub"
)

func Events(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	c, err := client.Upgrade(w, r)

	if err != nil {
		log.Println(err)
		return
	}

	hub.Register(c)

	defer func() {
		hub.Unregister(c)
	}()

	c.Listen()
}
