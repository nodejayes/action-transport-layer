package main

import (
	"local/action-transport-layer/cmd/server/actions"
	"local/action-transport-layer/pkg"
)

func main() {
	e := pkg.NewActionEndpoint(pkg.ActionEndpointConfig{
		Address: "localhost:3001",
		Auth: func(cl *pkg.Client) error {
			return nil
		},
		Actions: map[string]func(cl *pkg.Client, payload map[string]interface{}) error{
			"hello": actions.HelloHandler,
		},
	})
	e.Start()
}
