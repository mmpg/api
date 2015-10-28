package hub

import (
	"github.com/mmpg/api/client"
)

type hub struct {
	clients    map[*client.Client]bool
	broadcast  chan string
	register   chan *client.Client
	unregister chan *client.Client
}

var h = hub{
	broadcast:  make(chan string),
	register:   make(chan *client.Client),
	unregister: make(chan *client.Client),
	clients:    make(map[*client.Client]bool),
}

func Run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = true
			// TODO: Send current game state
			break

		case c := <-h.unregister:
			_, ok := h.clients[c]
			if ok {
				remove(c)
			}
			break

		case m := <-h.broadcast:
			broadcastMessage(m)
			break
		}
	}
}

func Register(c *client.Client) {
	h.register <- c
}

func Unregister(c *client.Client) {
	h.unregister <- c
}

func Broadcast(m string) {
	h.broadcast <- m
}

func remove(c *client.Client) {
	c.Close()
	delete(h.clients, c)
}

func broadcastMessage(m string) {
	for c := range h.clients {
		if !c.Send(m) {
			remove(c)
		}
	}
}
