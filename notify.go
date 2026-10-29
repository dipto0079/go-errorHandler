package errorhandler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"strings"
	"time"
)

// Notifier is called, in its own goroutine, the first time an error is seen.
type Notifier func(ctx context.Context, e Error) error

// Option configures a Handler.
type Option func(*Handler)

// WithEnv sets the environment name stored with captured errors (default "production").
func WithEnv(env string) Option { return func(h *Handler) { h.env = env } }

// WithService sets the creator name stored with captured errors (default "app").
func WithService(name string) Option { return func(h *Handler) { h.service = name } }

// WithNotifier adds a notifier. Several can be added.
func WithNotifier(n Notifier) Option { return func(h *Handler) { h.notifier = append(h.notifier, n) } }

func (h *Handler) notify(e Error) {
	for _, n := range h.notifier {
		n := n
		go func() {
			defer func() { _ = recover() }() // a faulty notifier must not crash the app
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			_ = n(ctx, e)
		}()
	}
}

func summary(e Error) string {
	return fmt.Sprintf("[%s] %s: %s", e.EnvType, e.ErrorCode, e.ErrorDetails)
}

// headerSafe removes line breaks so error text cannot inject mail headers.
func headerSafe(s string) string {
	return strings.NewReplacer("\r", " ", "\n", " ").Replace(s)
}

// WebhookNotifier POSTs JSON to url. The body has "text" (Slack) and
// "content" (Discord) set to a one-line summary, plus the error fields.
func WebhookNotifier(url string) Notifier {
	return func(ctx context.Context, e Error) error {
		body, err := json.Marshal(map[string]interface{}{
			"text":          summary(e),
			"content":       summary(e),
			"id":            e.ID,
			"error_code":    e.ErrorCode,
			"error_details": e.ErrorDetails,
			"env_type":      e.EnvType,
			"created_by":    e.CreatedBy,
			"stack_trace":   e.StackTrace,
		})
		if err != nil {
			return err
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 300 {
			return fmt.Errorf("errorhandler: webhook returned %s", resp.Status)
		}
		return nil
	}
}

// SMTPNotifier sends a plain-text email. addr is "host:port"; auth may be nil.
func SMTPNotifier(addr string, auth smtp.Auth, from string, to ...string) Notifier {
	return func(ctx context.Context, e Error) error {
		msg := fmt.Sprintf("From: %s\r\nSubject: %s\r\n\r\n%s\r\n\r\n%s\r\n",
			headerSafe(from), headerSafe(summary(e)), e.ErrorDetails, e.StackTrace)
		return smtp.SendMail(addr, auth, from, to, []byte(msg))
	}
}
