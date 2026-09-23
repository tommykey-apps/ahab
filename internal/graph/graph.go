package graph

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/tommykey-apps/ahab/internal/docker"
)

func nodeID(prefix, name string) string {
	var b strings.Builder
	b.WriteString(prefix)
	b.WriteByte('_')
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}
	return b.String()
}

func Render(cs []docker.Container) string {
	var b strings.Builder
	b.WriteString("graph LR\n")

	networks := map[string]bool{}
	volumes := map[string]bool{}
	hasHostPort := false
	for _, c := range cs {
		for _, n := range c.Networks {
			networks[n] = true
		}
		for _, v := range c.Volumes {
			volumes[v] = true
		}
		if len(c.Ports) > 0 {
			hasHostPort = true
		}
	}

	if hasHostPort {
		b.WriteString("  host{{\"host\"}}\n")
	}
	for _, n := range sorted(networks) {
		fmt.Fprintf(&b, "  %s((%q))\n", nodeID("n", n), n)
	}
	for _, v := range sorted(volumes) {
		fmt.Fprintf(&b, "  %s[(%q)]\n", nodeID("v", v), v)
	}

	ordered := slices.Clone(cs)
	slices.SortFunc(ordered, func(a, b docker.Container) int { return strings.Compare(a.Name, b.Name) })
	for _, c := range ordered {
		id := nodeID("c", c.Name)
		fmt.Fprintf(&b, "  %s[%q]\n", id, c.Name)
		for _, n := range sortedSlice(c.Networks) {
			fmt.Fprintf(&b, "  %s --> %s\n", id, nodeID("n", n))
		}
		for _, p := range c.Ports {
			fmt.Fprintf(&b, "  host -->|\":%d→%d\"| %s\n", p.Public, p.Private, id)
		}
		for _, v := range sortedSlice(c.Volumes) {
			fmt.Fprintf(&b, "  %s -.-> %s\n", id, nodeID("v", v))
		}
	}
	return b.String()
}

func sortedSlice(in []string) []string {
	out := slices.Clone(in)
	sort.Strings(out)
	return out
}

func sorted(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
