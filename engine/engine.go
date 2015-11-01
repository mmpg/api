package engine

import (
	"log"

	"github.com/pebbe/zmq4"
)

var (
	uri = "tcp://127.0.0.1:5555"
)

// Test engine connectivity
func Test() (s string, err error) {
	sck, err := zmq4.NewSocket(zmq4.REQ)
	defer sck.Close()

	if err != nil {
		log.Println("1")
		return
	}

	if err = sck.Connect(uri); err != nil {
		return
	}

	if _, err = sck.Send("Testing connectivity", 0); err != nil {
		return
	}

	s, err = sck.Recv(0)

	return
}
