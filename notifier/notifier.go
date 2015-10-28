package notifier

import (
	"time"

	"github.com/mmpg/api/hub"
)

// Run the event notifier
func Run() {
	for {
		hub.Broadcast("This is a test!")
		time.Sleep(5 * time.Second)
	}
}
