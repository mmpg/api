package engine

import (
	"errors"
	"strconv"
	"strings"
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

func request(t string, m string) (r string, connErr error, err error) {
	sck, err := newSocket(zmq4.REQ, 5555)
	defer sck.Close()

	sck.SetLinger(0)

	if connErr != nil {
		return
	}

	if _, connErr = sck.Send(t+" "+m, 0); connErr != nil {
		return
	}

	poller := zmq4.NewPoller()
	poller.Add(sck, zmq4.POLLIN)

	sockets, connErr := poller.Poll(requestTimeout)

	if connErr != nil {
		return
	}

	if len(sockets) < 1 {
		connErr = errors.New("Unable to connect with engine")
		return
	}

	r, connErr = sck.Recv(0)

	if connErr != nil {
		return
	}

	if strings.Contains(r, "ERROR") {
		parts := strings.Split(r, "|")

		if len(parts) > 1 {
			err = errors.New(parts[1])
		} else {
			err = errors.New("Unknown")
		}

		return
	}

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
func Log(t string) (string, error, error) {
	return request("LOG", t)
}
