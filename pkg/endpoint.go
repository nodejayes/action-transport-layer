package pkg

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type ActionEndpointConfig struct {
	Address string
	Auth    func(cl *Client) error
	Actions map[string]func(cl *Client, payload map[string]interface{}) error
}

type ActionEndpoint struct {
	Authenticator Authenticator
	Connections   *EndpointStore
	upgrader      websocket.Upgrader
	config        ActionEndpointConfig
}

func NewActionEndpoint(cfg ActionEndpointConfig) *ActionEndpoint {
	return &ActionEndpoint{
		Authenticator: Authenticator{},
		Connections:   NewEndpointStore(),
		config:        cfg,
		upgrader: websocket.Upgrader{
			EnableCompression: true,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			HandshakeTimeout: 30 * time.Second,
		},
	}
}

func (ae *ActionEndpoint) Start() {
	ae.Authenticator.CustomAuth = ae.config.Auth

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, upgradeError := ae.upgrader.Upgrade(w, r, nil)
		if upgradeError != nil {
			return
		}
		defer func() {
			_ = c.Close()
		}()
		cl := NewClient(c, ae)
		ae.Connections.add(cl)

		authenticateError := ae.Authenticator.Validate(cl)
		if authenticateError != nil {
			println(authenticateError.Error())
			// ends also the client connection
			ae.Connections.remove(cl)
			return
		}

		cl.listen()
		// remove the client from store while the client ends the listener or have some errors on read/write
		ae.Connections.remove(cl)
	})
	log.Fatal(http.ListenAndServe(ae.config.Address, nil))
}
