package client

import (
  "time"
  "net/http"

  "github.com/gorilla/websocket"
)

type Client struct {
  ws *websocket.Conn
  send chan []byte
}

const (
  writeWait = 10 * time.Second
  pongWait = 60 * time.Second
  pingPeriod = (pongWait * 9) / 10
  maxMessageSize = 1024 * 1024
)

var upgrader = websocket.Upgrader{
  ReadBufferSize:  maxMessageSize,
  WriteBufferSize: maxMessageSize,
}

func New(ws *websocket.Conn) *Client {
  return &Client{
    send: make(chan []byte, maxMessageSize),
    ws: ws,
  }
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*Client, error) {
  ws, err := upgrader.Upgrade(w, r, nil)

  if err != nil {
    return nil, err
  }

  return New(ws), nil
}

func (c *Client) Listen() {
  defer func() {
    c.ws.Close()
  }()

  go c.writeOutput()
  c.readInput()
}

func (c *Client) Send(m string) bool {
  select {
  case c.send <- []byte(m):
    return true

  // Channel is closed
  default:
    return false
  }
}

func (c *Client) Close() {
  close(c.send)
}

func (c *Client) readInput() {
  // Set read parameters
  c.ws.SetReadLimit(maxMessageSize)
  c.ws.SetPongHandler(func(string) error {
    c.ws.SetReadDeadline(time.Now().Add(pongWait))
    return nil
  })
  c.ws.SetReadDeadline(time.Now().Add(pongWait))

  for {
    // Discard input messages
    _, _, err := c.ws.ReadMessage()

    if err != nil {
      break
    }
  }
}

func (c *Client) writeOutput() {
  ticker := time.NewTicker(pingPeriod)

  defer func() {
    ticker.Stop()
  }()

  for {
    select {
    case m, ok := <-c.send:
      if !ok {
        c.write(websocket.CloseMessage, []byte{})
        return
      }
      if err := c.write(websocket.TextMessage, m); err != nil {
        return
      }
    case <-ticker.C:
      if err := c.write(websocket.PingMessage, []byte{}); err != nil {
        return
      }
    }
  }
}

func (c *Client) write(mt int, m []byte) error {
  c.ws.SetWriteDeadline(time.Now().Add(writeWait))
  return c.ws.WriteMessage(mt, m)
}
