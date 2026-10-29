package errorhandler

import (
	"fmt"
	"sort"
	"strings"

	migrationsql "github.com/dipto0079/go-errorHandler/errorHandler/migrations/sql"
)

// upSQL returns the "Up" part of every embedded goose migration, in order.
func upSQL() ([]string, error) {
	entries, err := migrationsql.FS.ReadDir(".")
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		names = append(names, e.Name())
	}
	sort.Strings(names)

	var out []string
	for _, n := range names {
		b, err := migrationsql.FS.ReadFile(n)
		if err != nil {
			return nil, err
		}
		up := string(b)
		if i := strings.Index(up, "-- +goose Down"); i >= 0 {
			up = up[:i]
		}
		var lines []string
		for _, l := range strings.Split(up, "\n") {
			if strings.HasPrefix(strings.TrimSpace(l), "-- +goose") {
				continue
			}
			lines = append(lines, l)
		}
		stmt := strings.TrimSpace(strings.Join(lines, "\n"))
		if stmt == "" {
			return nil, fmt.Errorf("migration %s has no Up statement", n)
		}
		out = append(out, stmt)
	}
	return out, nil
}
