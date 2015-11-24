package engine

import (
	"errors"
	"strconv"
	"time"

	"github.com/pebbe/zmq4"
)

var (
	host = "127.0.0.1"
)

const (
	requestTimeout = 2 * time.Second
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

func request(t string, m string) (r string, err error) {
	sck, err := newSocket(zmq4.REQ, 5555)
	defer sck.Close()

	sck.SetLinger(0)

	if err != nil {
		return
	}

	if _, err = sck.Send(t+" "+m, 0); err != nil {
		return
	}

	poller := zmq4.NewPoller()
	poller.Add(sck, zmq4.POLLIN)

	sockets, err := poller.Poll(requestTimeout)

	if err != nil {
		return
	}

	if len(sockets) < 1 {
		err = errors.New("Unable to connect with engine")
		return
	}

	r, err = sck.Recv(0)
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

// Log for the given time
func Log(t string) (string, error) {
	return request("LOG", t)
}
