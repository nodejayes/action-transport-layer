package pkg

type Authenticator struct {
	CustomAuth func(cl *Client) error
}

func (a *Authenticator) Validate(client *Client) error {
	if a.CustomAuth != nil {
		customError := a.CustomAuth(client)
		if customError != nil {
			return customError
		}
	}
	return nil
}
