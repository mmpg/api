package notifier

import (
	"log"
	"time"

	"github.com/mmpg/api/engine"
	"github.com/mmpg/api/hub"
)

// Run the event notifier
func Run() {
	for {
		err := engine.Subscribe(broadcast)

		if err != nil {
			log.Println(err)
		}

		time.Sleep(5 * time.Second)
	}
}

func broadcast(event string) {
	hub.Broadcast(event)
}
