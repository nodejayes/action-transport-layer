package pkg

type EndpointStore struct {
	endpoints map[string]*Client
}

func NewEndpointStore() *EndpointStore {
	return &EndpointStore{
		endpoints: make(map[string]*Client),
	}
}

func (es *EndpointStore) add(client *Client) {
	es.endpoints[client.ID] = client
}

func (es *EndpointStore) remove(client *Client) {
	if es.endpoints[client.ID] != nil {
		_ = client.Connection.Close()
		delete(es.endpoints, client.ID)
	}
}

func (es *EndpointStore) GetClients(filter func(client *Client) bool) []*Client {
	var selected []*Client
	for _, cl := range es.endpoints {
		if filter(cl) {
			selected = append(selected, cl)
		}
	}
	if selected == nil {
		return []*Client{}
	}
	return selected
}
