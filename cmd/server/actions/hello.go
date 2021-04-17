package actions

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"local/action-transport-layer/pkg"
)

type HelloPayload struct {
	Name string
}

type HelloResponse struct {
	Result string `json:"result"`
	Error  error
}

func (hr HelloResponse) GetResult() interface{} {
	return hr.Result
}

func (hr HelloResponse) GetError() string {
	if hr.Error != nil {
		return hr.Error.Error()
	}
	return ""
}

func HelloHandler(cl *pkg.Client, payload map[string]interface{}) error {
	var params HelloPayload
	err := mapstructure.Decode(payload, &params)
	if err != nil {
		return err
	}

	cl.Send("hello", HelloResponse{Result: fmt.Sprintf("Hello %v", params.Name)})
	return nil
}
