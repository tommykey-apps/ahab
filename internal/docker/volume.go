package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"time"
)

type Volume struct {
	Name       string
	Driver     string
	Mountpoint string
	Created    time.Time
	InUse      bool
}

type apiVolume struct {
	Name       string `json:"Name"`
	Driver     string `json:"Driver"`
	Mountpoint string `json:"Mountpoint"`
	CreatedAt  string `json:"CreatedAt"`
}

func (c *Client) Volumes(ctx context.Context) ([]Volume, error) {
	all, err := c.listVolumes(ctx, "")
	if err != nil {
		return nil, err
	}
	dangling, err := c.listVolumes(ctx, `{"dangling":["true"]}`)
	if err != nil {
		return nil, err
	}
	unused := map[string]bool{}
	for _, v := range dangling {
		unused[v.Name] = true
	}

	vols := make([]Volume, 0, len(all))
	for _, r := range all {
		t, _ := time.Parse(time.RFC3339, r.CreatedAt)
		vols = append(vols, Volume{
			Name:       r.Name,
			Driver:     r.Driver,
			Mountpoint: r.Mountpoint,
			Created:    t,
			InUse:      !unused[r.Name],
		})
	}
	sort.Slice(vols, func(i, j int) bool { return vols[i].Name < vols[j].Name })
	return vols, nil
}

func (c *Client) listVolumes(ctx context.Context, filters string) ([]apiVolume, error) {
	path := "/volumes"
	if filters != "" {
		path += "?" + url.Values{"filters": {filters}}.Encode()
	}
	body, err := c.get(ctx, path)
	if err != nil {
		return nil, err
	}
	defer body.Close()
	var raw struct {
		Volumes []apiVolume `json:"Volumes"`
	}
	if err := json.NewDecoder(body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decode volumes: %w", err)
	}
	return raw.Volumes, nil
}

func (c *Client) RemoveVolume(ctx context.Context, name string, force bool) error {
	q := url.Values{}
	if force {
		q.Set("force", "1")
	}
	return c.remove(ctx, "/volumes/"+url.PathEscape(name)+"?"+q.Encode())
}
