package gotify

type Client struct {
	adapters []adapter
}

func (c *Client) Send(subject string, message []Line) error {
	for _, ad := range c.adapters {
		err := ad.send(subject, message...)
		if err != nil {
			return err
		}
	}

	return nil
}

type Options interface {
	init() (adapter, error)
}

func NewClient(opts ...Options) (*Client, error) {
	client := &Client{}
	for _, o := range opts {

		// Init adapter
		ad, err := o.init()
		if err != nil {
			return nil, err
		}

		client.adapters = append(client.adapters, ad)
	}

	return client, nil
}
