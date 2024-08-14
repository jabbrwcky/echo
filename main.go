package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"net/http"
	"os"
	"os/signal"

	"github.com/rs/zerolog/log"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-echo-code", "200")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/", http.HandlerFunc(echo))

	srv := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Error().Err(err).Msg("HTTP server shutdown")
		}
		log.Info().Msg("HTTP server shutdown finished")
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Err(err).Msg("error starting api server")
	}
}

func echo(w http.ResponseWriter, req *http.Request) {
	req.ParseMultipartForm(1204 * 1024)
	b, err := io.ReadAll(req.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read body")
	}

	resp := map[string]interface{}{
		"content_length":    req.ContentLength,
		"host":              req.Host,
		"method":            req.Method,
		"proto":             req.Proto,
		"remote_addr":       req.RemoteAddr,
		"transfer_encoding": req.TransferEncoding,
		"header":            req.Header,
		"form":              req.Form,
		"request":           req.RequestURI,
		"trailer":           req.Trailer,
		"user_agent":        req.UserAgent(),
		"referer":           req.Referer(),
		"cookies":           req.Cookies(),
		"body":              b,
	}

	enc := json.NewEncoder(w)
	enc.Encode(resp)
}
