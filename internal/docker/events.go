package docker

import (
	"context"
	"encoding/json"
	"net/url"
)

type Event struct {
	Action string `json:Action`
	Actor struct {
		ID string `json:ID`
	} `json:Actor`
}

func (c *Client) Events(ctx context.Context, out chan<- Event) error {
	q := url.Values{"filters": {`{"type":["container"]}`}}
	body, err := c.get(ctx, "/events?"+q.Encode())
	if err != nil {
		return err
	}
	defer body.Close()

	dec := json.NewDecoder(body)
	for {
		var ev Event
		if err := dec.Decode(&ev); err != nil {
			return err
		}
		select {
		case out <- ev:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
