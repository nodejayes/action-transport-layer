package pkg

type Action struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}
