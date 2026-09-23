package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type Container struct {
	ID       string
	Name     string
	Image    string
	ImageID  string
	State    string
	Ports    []Port
	Networks []string
	Volumes  []string
}

type Port struct {
	Public  int
	Private int
}

type apiContainer struct {
	ID      string   `json:"Id"`
	Names   []string `json:"Names"`
	Image   string   `json:"Image"`
	ImageID string   `json:"ImageID"`
	State   string   `json:"State"`
	Ports   []struct {
		PrivatePort int `json:"PrivatePort"`
		PublicPort  int `json:"PublicPort"`
	} `json:"Ports"`
	NetworkSettings struct {
		Networks map[string]struct{} `json:"Networks"`
	} `json:"NetworkSettings"`
	Mounts []struct {
		Type string `json:"Type"`
		Name string `json:"Name"`
	} `json:"Mounts"`
}

func (c *Client) Containers(ctx context.Context) ([]Container, error) {
	body, err := c.get(ctx, "/containers/json?all=1")
	if err != nil {
		return nil, err
	}
	defer body.Close()
	var raw []apiContainer
	if err := json.NewDecoder(body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decode containers: %w", err)
	}

	cs := make([]Container, 0, len(raw))
	for _, r := range raw {
		name := ""
		if len(r.Names) > 0 {
			name = strings.TrimPrefix(r.Names[0], "/")
		}
		c := Container{ID: r.ID, Name: name, Image: r.Image, ImageID: r.ImageID, State: r.State}
		for _, p := range r.Ports {
			if p.PublicPort != 0 {
				c.Ports = append(c.Ports, Port{Public: p.PublicPort, Private: p.PrivatePort})
			}
		}
		for n := range r.NetworkSettings.Networks {
			c.Networks = append(c.Networks, n)
		}
		sort.Strings(c.Networks)
		for _, m := range r.Mounts {
			if m.Type == "volume" && !isAnonymous(m.Name) {
				c.Volumes = append(c.Volumes, m.Name)
			}
		}
		sort.Strings(c.Volumes)
		cs = append(cs, c)
	}
	sort.Slice(cs, func(i, j int) bool { return cs[i].Name < cs[j].Name })
	return cs, nil
}

func isAnonymous(name string) bool {
	if len(name) != 64 {
		return false
	}
	for _, r := range name {
		if !(r >= '0' && r <= '9' || r >= 'a' && r <= 'f') {
			return false
		}
	}
	return true
}
