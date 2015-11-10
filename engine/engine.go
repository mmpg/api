package engine

import (
	"strconv"

	"github.com/pebbe/zmq4"
)

var (
	host = "127.0.0.1"
)

type handler func(string)

func newSocket(t zmq4.Type, port int) (sck *zmq4.Socket, err error) {
	sck, err = zmq4.NewSocket(t)

	if err != nil {
		return
	}

	err = sck.Connect("tcp://" + host + ":" + strconv.Itoa(port))

	return
}

// Subscribe to the engine and pass new events to the given handler
func Subscribe(fn handler) error {
	sck, err := newSocket(zmq4.SUB, 5556)
	defer sck.Close()

	if err != nil {
		return err
	}

	sck.SetSubscribe("")

	for {
		s, err := sck.Recv(0)

		if err != nil {
			return err
		}

		fn(s)
	}
}

// Test engine connectivity
func Test() (s string, err error) {
	sck, err := newSocket(zmq4.REQ, 5555)
	defer sck.Close()

	if err != nil {
		return
	}

	if _, err = sck.Send("Testing connectivity", 0); err != nil {
		return
	}

	s, err = sck.Recv(0)

	return
}
