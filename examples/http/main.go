// Command http shows errorhandler recording errors and panics in a web server.
//
//	DATABASE_URL="user=postgres password=secret host=localhost dbname=mydb sslmode=disable" go run ./examples/http
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	errorhandler "github.com/dipto0079/go-errorHandler"
)

func main() {
	ctx := context.Background()

	opts := []errorhandler.Option{errorhandler.WithEnv("development"), errorhandler.WithService("example")}
	if url := os.Getenv("ERROR_WEBHOOK_URL"); url != "" {
		opts = append(opts, errorhandler.WithNotifier(errorhandler.WebhookNotifier(url)))
	}
	h, err := errorhandler.New(os.Getenv("DATABASE_URL"), opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer h.Close()
	if err := h.Migrate(ctx); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/fail", func(w http.ResponseWriter, r *http.Request) {
		if _, err := h.Capture(r.Context(), errors.New("something failed")); err != nil {
			log.Print(err)
		}
		http.Error(w, "failed", http.StatusInternalServerError)
	})
	mux.HandleFunc("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("oh no")
	})
	mux.HandleFunc("/errors", func(w http.ResponseWriter, r *http.Request) {
		list, err := h.List(r.Context(), errorhandler.ListOptions{Status: errorhandler.StatusOpen, Limit: 20})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, e := range list {
			fmt.Fprintf(w, "#%s x%d [%s] %s: %s\n", e.ID, e.OccurrenceCount, e.EnvType, e.ErrorCode, e.ErrorDetails)
		}
	})

	log.Println("listening on :8080  (try /fail, /panic, /errors)")
	log.Fatal(http.ListenAndServe(":8080", h.Middleware(mux)))
}
