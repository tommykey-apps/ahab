package docker

import (
	"context"
	"encoding/binary"
	"io"
)

func (c *Client) Logs(ctx context.Context, id string, out chan<- string) error {
	body, err := c.get(ctx, "/containers/"+id+"/logs?follow=1&stdout=1&stderr=1&tail=100")
	if err != nil {
		return err
	}
	defer body.Close()

	hdr := make([]byte, 8)
	for {
		if _, err := io.ReadFull(body, hdr); err != nil {
			return err
		}
		payload := make([]byte, binary.BigEndian.Uint32(hdr[4:8]))
		if _, err := io.ReadFull(body, payload); err != nil {
			return err
		}
		select {
		case out <- string(payload):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
