package api

import (
	"log"
	"net/http"

	"github.com/mmpg/api/client"
	"github.com/mmpg/api/engine"
	"github.com/mmpg/api/hub"
	"github.com/mmpg/api/notifier"
)

// Run the MMPG Api:
// 1. Starts the subscriber hub
// 2. Starts the event notifier
// 3. Starts the API server
func Run() {
	go hub.Run()
	go notifier.Run()

	http.HandleFunc("/events", serveEvents)
	http.HandleFunc("/test", serveTest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveEvents(w http.ResponseWriter, r *http.Request) {
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

func serveTest(w http.ResponseWriter, r *http.Request) {
	reply, err := engine.Test()

	if err != nil {
		log.Println(err)
		return
	}

	w.Write([]byte(reply))
}
