package errorhandler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

// newTestHandler creates an isolated database on the server named by
// DATABASE_CONNECTION (a lib/pq string without dbname) and migrates it.
func newTestHandler(t *testing.T, opts ...Option) *Handler {
	t.Helper()
	conn := os.Getenv("DATABASE_CONNECTION")
	if conn == "" {
		t.Skip("DATABASE_CONNECTION not set")
	}
	name := fmt.Sprintf("eh_test_%d", rand.Int63())
	admin, err := sql.Open("postgres", conn)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := admin.Exec(`CREATE DATABASE "` + name + `"`); err != nil {
		t.Fatal(err)
	}
	h, err := New(conn+" dbname="+name, opts...)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		h.Close()
		admin.Exec(`DROP DATABASE "` + name + `" WITH (FORCE)`)
		admin.Close()
	})
	ctx := context.Background()
	if err := h.Migrate(ctx); err != nil {
		t.Fatal(err)
	}
	if err := h.Migrate(ctx); err != nil {
		t.Fatalf("second Migrate must be a no-op: %v", err)
	}
	return h
}

func TestCRUD(t *testing.T) {
	h := newTestHandler(t)
	ctx := context.Background()

	if _, err := h.Get(ctx, "999"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get missing = %v, want ErrNotFound", err)
	}
	if err := h.Delete(ctx, "999", "test"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Delete missing = %v, want ErrNotFound", err)
	}

	id, err := h.Create(ctx, Error{ErrorCode: "E1", ErrorDetails: "boom", EnvType: "test", CreatedBy: "t"})
	if err != nil {
		t.Fatal(err)
	}
	got, err := h.Get(ctx, id)
	if err != nil || got.ErrorCode != "E1" || got.Status != StatusOpen || got.CreatedAt.IsZero() {
		t.Fatalf("Get = %+v, %v", got, err)
	}

	if err := h.Delete(ctx, id, "test"); err != nil {
		t.Fatal(err)
	}
	if _, err := h.Get(ctx, id); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get deleted = %v, want ErrNotFound", err)
	}
	if l, _ := h.List(ctx); len(l) != 0 {
		t.Fatalf("List shows deleted record: %+v", l)
	}
	l, err := h.List(ctx, ListOptions{IncludeDeleted: true})
	if err != nil || len(l) != 1 || !l[0].DeletedAt.Valid || l[0].DeleteByEnvType != "test" {
		t.Fatalf("List IncludeDeleted = %+v, %v", l, err)
	}
}

func TestListFilters(t *testing.T) {
	h := newTestHandler(t)
	ctx := context.Background()
	for i := 0; i < 5; i++ {
		env := "prod"
		if i%2 == 1 {
			env = "dev"
		}
		if _, err := h.Create(ctx, Error{ErrorCode: fmt.Sprintf("E%d", i), ErrorDetails: "x", EnvType: env, CreatedBy: "t"}); err != nil {
			t.Fatal(err)
		}
	}
	check := func(name string, o ListOptions, want int) {
		t.Helper()
		l, err := h.List(ctx, o)
		if err != nil || len(l) != want {
			t.Errorf("%s: got %d records, err %v, want %d", name, len(l), err, want)
		}
	}
	check("all", ListOptions{}, 5)
	check("env", ListOptions{Env: "dev"}, 2)
	check("code", ListOptions{Code: "E3"}, 1)
	check("limit", ListOptions{Limit: 2}, 2)
	check("offset", ListOptions{Limit: 10, Offset: 3}, 2)
	check("since future", ListOptions{Since: time.Now().Add(time.Hour)}, 0)
	check("since past", ListOptions{Since: time.Now().Add(-time.Hour)}, 5)

	l, _ := h.List(ctx, ListOptions{Limit: 1})
	if l[0].ErrorCode != "E4" {
		t.Errorf("newest first: got %s", l[0].ErrorCode)
	}
}

func TestCaptureDedupAndResolve(t *testing.T) {
	notified := make(chan Error, 4)
	h := newTestHandler(t, WithEnv("test"), WithService("svc"),
		WithNotifier(func(_ context.Context, e Error) error { notified <- e; return nil }))
	ctx := context.Background()

	if _, err := h.Capture(ctx, nil); err == nil {
		t.Fatal("Capture(nil) should fail")
	}

	id1, err := h.Capture(ctx, errors.New("db down"))
	if err != nil {
		t.Fatal(err)
	}
	id2, err := h.Capture(ctx, errors.New("db down"))
	if err != nil {
		t.Fatal(err)
	}
	if id1 != id2 {
		t.Fatalf("same error got two rows: %s vs %s", id1, id2)
	}
	got, _ := h.Get(ctx, id1)
	if got.OccurrenceCount != 2 || got.EnvType != "test" || got.CreatedBy != "svc" || got.ErrorCode != "ERROR" {
		t.Fatalf("record = %+v", got)
	}
	if !strings.Contains(got.StackTrace, "goroutine") {
		t.Errorf("no stack trace stored: %q", got.StackTrace)
	}

	select {
	case e := <-notified:
		if e.ID != id1 {
			t.Errorf("notified with id %s, want %s", e.ID, id1)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("no notification for first occurrence")
	}
	select {
	case <-notified:
		t.Fatal("repeat occurrence must not notify")
	case <-time.After(300 * time.Millisecond):
	}

	if err := h.Resolve(ctx, id1); err != nil {
		t.Fatal(err)
	}
	if got, _ := h.Get(ctx, id1); got.Status != StatusResolved || !got.ResolvedAt.Valid {
		t.Fatalf("not resolved: %+v", got)
	}
	id3, _ := h.Capture(ctx, errors.New("db down"), Code("DB"))
	id4, _ := h.Capture(ctx, errors.New("db down"))
	if id3 == id1 || id3 == id4 {
		t.Errorf("resolved error must reopen as new row (and a different code is a different error): %s %s %s", id1, id3, id4)
	}
	if open, _ := h.List(ctx, ListOptions{Status: StatusOpen}); len(open) != 2 {
		t.Errorf("open = %d, want 2", len(open))
	}
	if err := h.Resolve(ctx, "999"); !errors.Is(err, ErrNotFound) {
		t.Errorf("Resolve missing = %v", err)
	}
}

func TestMiddlewareRecordsPanic(t *testing.T) {
	h := newTestHandler(t)
	srv := httptest.NewServer(h.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/boom" {
			panic("kaboom")
		}
		fmt.Fprint(w, "ok")
	})))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/ok")
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("/ok: %v %v", resp, err)
	}
	resp, err = http.Get(srv.URL + "/boom")
	if err != nil || resp.StatusCode != 500 {
		t.Fatalf("/boom: %v %v", resp, err)
	}
	l, _ := h.List(context.Background(), ListOptions{Code: "PANIC"})
	if len(l) != 1 || !strings.Contains(l[0].ErrorDetails, "kaboom") || !strings.Contains(l[0].ErrorDetails, "/boom") {
		t.Fatalf("panic not recorded: %+v", l)
	}
}

func TestWebhookNotifier(t *testing.T) {
	got := make(chan string, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b := make([]byte, 4096)
		n, _ := r.Body.Read(b)
		got <- string(b[:n])
	}))
	defer srv.Close()

	err := WebhookNotifier(srv.URL)(context.Background(), Error{EnvType: "prod", ErrorCode: "E9", ErrorDetails: "bad"})
	if err != nil {
		t.Fatal(err)
	}
	if body := <-got; !strings.Contains(body, `"text":"[prod] E9: bad"`) {
		t.Errorf("body = %s", body)
	}
}

func TestRegressions(t *testing.T) {
	h := newTestHandler(t, WithNotifier(func(context.Context, Error) error { panic("bad notifier") }))
	ctx := context.Background()

	for name, err := range map[string]error{
		"Get":     func() error { _, err := h.Get(ctx, "abc"); return err }(),
		"Delete":  h.Delete(ctx, "abc", "t"),
		"Resolve": h.Resolve(ctx, "abc"),
	} {
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("%s with non-numeric id = %v, want ErrNotFound", name, err)
		}
	}

	before := time.Now()
	id, err := h.Create(ctx, Error{ErrorCode: "T", ErrorDetails: "t", EnvType: "e", CreatedBy: "c"})
	if err != nil {
		t.Fatal(err)
	}
	got, _ := h.Get(ctx, id)
	if d := got.CreatedAt.Sub(before); d < -time.Minute || d > time.Minute {
		t.Errorf("CreatedAt off by %s (timezone bug?)", d)
	}

	if err := h.Delete(ctx, id, "a"); err != nil {
		t.Fatal(err)
	}
	if err := h.Delete(ctx, id, "b"); !errors.Is(err, ErrNotFound) {
		t.Errorf("second Delete = %v, want ErrNotFound", err)
	}

	// A panicking notifier must not take the process down.
	if _, err := h.Capture(ctx, errors.New("trigger notifier")); err != nil {
		t.Fatal(err)
	}
	time.Sleep(200 * time.Millisecond)
}

func TestHeaderSafe(t *testing.T) {
	if got := headerSafe("a\r\nBcc: x@y.z"); strings.ContainsAny(got, "\r\n") {
		t.Errorf("headerSafe left a line break: %q", got)
	}
}
