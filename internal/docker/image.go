package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"time"
)

type Image struct {
	ID      string
	Tags    []string
	Size    int64
	Created time.Time
}

type apiImage struct {
	ID       string   `json:"Id"`
	RepoTags []string `json:"RepoTags"`
	Size     int64    `json:"Size"`
	Created  int64    `json:"Created"`
}

func (c *Client) Images(ctx context.Context) ([]Image, error) {
	body, err := c.get(ctx, "/images/json")
	if err != nil {
		return nil, err
	}
	defer body.Close()
	var raw []apiImage
	if err := json.NewDecoder(body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decode images: %w", err)
	}

	imgs := make([]Image, 0, len(raw))
	for _, r := range raw {
		img := Image{ID: r.ID, Size: r.Size, Created: time.Unix(r.Created, 0)}
		for _, t := range r.RepoTags {
			if t != "<none>:<none>" {
				img.Tags = append(img.Tags, t)
			}
		}
		sort.Strings(img.Tags)
		imgs = append(imgs, img)
	}
	sort.Slice(imgs, func(i, j int) bool {
		a, b := imgs[i], imgs[j]
		switch {
		case len(a.Tags) == 0 && len(b.Tags) > 0:
			return false
		case len(a.Tags) > 0 && len(b.Tags) == 0:
			return true
		case len(a.Tags) == 0:
			return a.ID < b.ID
		default:
			return a.Tags[0] < b.Tags[0]
		}
	})
	return imgs, nil
}

func (c *Client) RemoveImage(ctx context.Context, id string, force bool) error {
	q := url.Values{}
	if force {
		q.Set("force", "1")
	}
	return c.remove(ctx, "/images/"+url.PathEscape(id)+"?"+q.Encode())
}
