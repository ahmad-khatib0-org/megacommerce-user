package oauth

import (
	"context"
	"net/http"
	"time"

	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/models"
	"github.com/go-chi/chi/v5"
)

func (oa *OAuth) Shutdown() {
	ie := func(err error, msg string) {
		ie := &models.InternalError{Err: err, Msg: msg, Path: "users.oauth.Run"}
		oa.log.ErrorStruct(msg, ie)

		// Try to send error, but don't block if channel is full/closed
		select {
		case oa.errCh <- ie:
		default:
			ie.Msg = "error channel blocked during shutdown"
			oa.log.ErrorStruct(ie.Msg, ie)
		}
	}

	ctx, done := context.WithTimeout(context.Background(), time.Second*10)
	defer done()

	err := oa.server.Shutdown(ctx)
	if err != nil {
		ie(err, "an error occurred while shutting down the OAuth server")
	}
}

func (oa *OAuth) routes() *chi.Mux {
	mux := chi.NewRouter()

	mux.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusMethodNotAllowed)
		msg := `{"error": true, "message": "http method is not allowed"}`
		w.Write([]byte(msg))
	})

	mux.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		msg := `{"error": true, "message": "the requested url is not found"}`
		w.Write([]byte(msg))
	})

	mux.Get("/login", oa.Login)
	mux.Get("/error", oa.Error)
	mux.Get("/consent", oa.Consent)

	return mux
}
