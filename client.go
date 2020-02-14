package gotify

import "errors"

func NewClient(adaptersInit ...func() (adapter Adapter, err error)) (*Client, error) {
	client := &Client{}

	for _, init := range adaptersInit {
		// Init adapter
		adapter, err := init()
		if err != nil {
			return nil, err
		}

		client.adapters = append(client.adapters, adapter)
	}

	if len(client.adapters) == 0 {
		return nil, errors.New("newClient: adapter list is empty, please add one adapter")
	}

	return client, nil
}

type Client struct {
	adapters []Adapter
}

func (c Client) Send(subject string, message ...Line) (errs []error) {
	for _, adapter := range c.adapters {
		if err := adapter.Send(subject, message...); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}
