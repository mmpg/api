package engine

import (
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/pebbe/zmq4"
)

var (
	host = "127.0.0.1"
	// ErrConnectionFailed represents the error that happens when connection
	// with the engine cannot be established
	ErrConnectionFailed = errors.New("engine: connection failed")

	// ErrInvalidBase64Encoding is thrown when the engine returns an invalid
	// base64 encoding
	ErrInvalidBase64Encoding = errors.New("engine: invalid base64 encoding")
)

const (
	requestTimeout = 5 * time.Second
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
		err = ErrConnectionFailed
		return
	}

	if _, err = sck.Send(t+" "+m, 0); err != nil {
		err = ErrConnectionFailed
		return
	}

	poller := zmq4.NewPoller()
	poller.Add(sck, zmq4.POLLIN)

	sockets, err := poller.Poll(requestTimeout)

	if err != nil {
		err = ErrConnectionFailed
		return
	}

	if len(sockets) < 1 {
		err = ErrConnectionFailed
		return
	}

	r, err = sck.Recv(0)

	if err != nil {
		err = ErrConnectionFailed
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
func Log(t string) ([]byte, error) {
	res, err := request("LOG", t)

	if err != nil {
		return nil, err
	}

	d, err := base64.StdEncoding.DecodeString(res)

	if err != nil {
		return nil, ErrInvalidBase64Encoding
	}

	return d, nil
}

// PlayerExists checks if a player exists in the engine
func PlayerExists(email string) (string, error) {
	return request("PLAYER_EXISTS", email)
}

// DeployPlayer compiles, installs and restarts the given player for the given
// email
func DeployPlayer(email string, p io.Reader) (err error) {
	c, err := ioutil.ReadAll(p)

	if err != nil {
		return
	}

	_, err = request("DEPLOY_PLAYER", email+" "+base64.StdEncoding.EncodeToString(c))
	return
}

// World returns the structure of the current game world
func World() (string, error) {
	return request("WORLD", "")
}
