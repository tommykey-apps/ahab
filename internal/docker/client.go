package docker

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
)

const apiVersion = "v1.44"

type Client struct {
	http *http.Client
}

func New() *Client {
	sock := socketPath()
	return NewWithTransport(&http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "unix", sock)
		},
	})
}

func NewWithTransport(t http.RoundTripper) *Client {
	return &Client{http: &http.Client{Transport: t}}
}

func socketPath() string {
	if h := os.Getenv("DOCKER_HOST"); strings.HasPrefix(h, "unix://") {
		return strings.TrimPrefix(h, "unix://")
	}
	return "/var/run/docker.sock"
}

func (c *Client) get(ctx context.Context, path string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://docker/"+apiVersion+path, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("docker api %s: %w", path, err)
	}
	if res.StatusCode != http.StatusOK {
		res.Body.Close()
		return nil, fmt.Errorf("docker api %s: status %d", path, res.StatusCode)
	}
	return res.Body, nil
}
