package view

import (
	"fmt"
	"strings"

	"myspoti/internal/player"
)

// Queue renders recent → current → upcoming as a tree.
func Queue(q player.QueueView) string {
	var b strings.Builder

	for _, t := range q.Recent {
		fmt.Fprintf(&b, "├─ %s\n", t.Label())
	}

	if q.Current == nil {
		b.WriteString("╰─▶ (nothing playing)")
	} else {
		fmt.Fprintf(&b, "╰─▶ %s", q.Current.Label())
	}

	if len(q.Upcoming) == 0 {
		return b.String()
	}

	b.WriteByte('\n')
	for i, t := range q.Upcoming {
		prefix := "   ├─ "
		if i == len(q.Upcoming)-1 {
			prefix = "   ╰─ "
		}
		b.WriteString(prefix)
		b.WriteString(t.Label())
		if i < len(q.Upcoming)-1 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}
