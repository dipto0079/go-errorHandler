package errorhandler

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

// CaptureOption customises a single Capture call.
type CaptureOption func(*Error)

// Code sets the error code. By default it is taken from an ErrorCode() string
// method on the error, or "ERROR".
func Code(code string) CaptureOption { return func(e *Error) { e.ErrorCode = code } }

// Meta adds a key/value to the record's metadata (user id, request id, region...).
// Do not put secrets or personal data in it.
func Meta(key, value string) CaptureOption {
	return func(e *Error) {
		if e.Metadata == nil {
			e.Metadata = Metadata{}
		}
		e.Metadata[key] = value
	}
}

// WithCode attaches an error code to err. Capture uses it as the record's
// code. The original error stays available through errors.Is/As.
func WithCode(err error, code string) error {
	if err == nil {
		return nil
	}
	return &codedError{err: err, code: code}
}

type codedError struct {
	err  error
	code string
}

func (c *codedError) Error() string     { return c.err.Error() }
func (c *codedError) Unwrap() error     { return c.err }
func (c *codedError) ErrorCode() string { return c.code }

// By sets who created the record, overriding the handler's service name.
func By(who string) CaptureOption { return func(e *Error) { e.CreatedBy = who } }

// Capture records err with its stack trace. Identical errors (same code,
// environment and message) that are still open are counted on one record
// instead of creating a new one. It returns the record id.
//
// The first occurrence also triggers the configured notifiers.
func (h *Handler) Capture(ctx context.Context, err error, opts ...CaptureOption) (string, error) {
	if err == nil {
		return "", errors.New("errorhandler: Capture called with nil error")
	}
	return h.capture(ctx, err, string(debug.Stack()), opts)
}

func (h *Handler) capture(ctx context.Context, err error, stack string, opts []CaptureOption) (string, error) {
	e := Error{
		ErrorCode:    "ERROR",
		ErrorDetails: err.Error(),
		EnvType:      h.env,
		CreatedBy:    h.service,
		CreatedAt:    time.Now(),
		StackTrace:   stack,
	}
	var coded interface{ ErrorCode() string }
	if errors.As(err, &coded) {
		e.ErrorCode = coded.ErrorCode()
	}
	for _, o := range opts {
		o(&e)
	}
	e.Fingerprint = fingerprint(e)

	id, count, serr := h.st.SaveError(ctx, e)
	if serr != nil {
		return "", serr
	}
	if count == 1 {
		e.ID = id
		e.OccurrenceCount = 1
		h.notify(e)
	}
	return id, nil
}

func fingerprint(e Error) string {
	sum := sha256.Sum256([]byte(e.EnvType + "\x00" + e.ErrorCode + "\x00" + e.ErrorDetails))
	return hex.EncodeToString(sum[:16])
}

// Middleware recovers panics in next, records them with Capture and answers
// 500. It works with net/http and routers such as gorilla/mux and chi.
func (h *Handler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rec := recover()
			if rec == nil {
				return
			}
			if rec == http.ErrAbortHandler {
				panic(rec)
			}
			err := fmt.Errorf("%s %s: panic: %v", r.Method, r.URL.Path, rec)
			// The request context is cancelled when the client disconnects.
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_, _ = h.capture(ctx, err, string(debug.Stack()), []CaptureOption{Code("PANIC"), Meta("method", r.Method), Meta("path", r.URL.Path)})
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}()
		next.ServeHTTP(w, r)
	})
}
