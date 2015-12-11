package endpoints

import (
	"net/http"

	"github.com/mmpg/api/client"
	"github.com/mmpg/api/hub"
)

// Events upgrades the http request to a subscription using websockets
func Events(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	c, err := client.Upgrade(w, r)

	if handleError(w, err) {
		return
	}

	hub.Register(c)

	defer func() {
		hub.Unregister(c)
	}()

	c.Listen()
}
