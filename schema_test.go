package errorhandler

import (
	"strings"
	"testing"
)

func TestUpSQL(t *testing.T) {
	got, err := upSQL()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) == 0 {
		t.Fatal("no migrations embedded")
	}
	for _, s := range got {
		if !strings.Contains(s, "errorhandler") {
			t.Errorf("unexpected statement: %q", s)
		}
		if strings.Contains(s, "goose") || strings.Contains(s, "DROP TABLE") {
			t.Errorf("Down or goose marker leaked: %q", s)
		}
	}
}
