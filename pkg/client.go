package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Payload interface {
	GetResult() interface{}
	GetError() string
}

type Client struct {
	ID                string
	Context           map[string]interface{}
	Connection        *websocket.Conn
	Endpoint          *ActionEndpoint
	ReadEvent         chan []byte
	WriteEvent        chan []byte
	MessageType       int
	UpgradeError      error
	AuthenticateError error
	MessageReadError  error
	MessageWriteError error
}

func NewClient(c *websocket.Conn, endpoint *ActionEndpoint) *Client {
	return &Client{
		ID:                uuid.NewString(),
		Context:           make(map[string]interface{}),
		Connection:        c,
		Endpoint:          endpoint,
		ReadEvent:         make(chan []byte),
		WriteEvent:        make(chan []byte),
		MessageType:       0,
		UpgradeError:      nil,
		AuthenticateError: nil,
		MessageReadError:  nil,
		MessageWriteError: nil,
	}
}

func (cl *Client) listen() {
	// start the Listeners
	cl.receiveData()
	cl.sendData()

	for {
		mt, message, err := cl.Connection.ReadMessage()
		if err != nil {
			cl.MessageReadError = err
			break
		}
		cl.MessageType = mt
		cl.parseMessage(message)

		if cl.MessageWriteError != nil {
			break
		}
	}
}

func (cl *Client) Send(typ string, payload Payload) {
	go func() {
		msg, err := json.Marshal(Action{
			Type: typ,
			Payload: map[string]interface{}{
				"Result": payload.GetResult(),
				"Error":  payload.GetError(),
			},
		})
		if err != nil {
			println(err.Error())
			return
		}
		cl.WriteEvent <- msg
	}()
}

func (cl *Client) parseMessage(message []byte) {
	go func() {
		cl.ReadEvent <- message
	}()
}

func (cl *Client) receiveData() {
	go func() {
		for {
			msg := <-cl.ReadEvent

			// parse the sent Data into a Action Struct
			var parsedAction Action
			err := json.Unmarshal(msg, &parsedAction)
			if err != nil {
				println(err.Error())
				continue
			}

			// select from Action List and execute the Action
			action := cl.Endpoint.config.Actions[parsedAction.Type]
			if action == nil {
				println(fmt.Sprintf("the Action from Type %v not found", parsedAction.Type))
				continue
			}
			err = action(cl, parsedAction.Payload)
			if err != nil {
				println(err.Error())
			}
		}
	}()
}

func (cl *Client) sendData() {
	go func() {
		for {
			msg := <-cl.WriteEvent
			cl.MessageWriteError = cl.Connection.WriteMessage(cl.MessageType, []byte(msg))
			if cl.MessageWriteError != nil {
				break
			}
		}
	}()
}
