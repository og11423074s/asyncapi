package apiserver

import (
	"context"
	"github.com/og11423074s/asyncapi/config"
	"github.com/og11423074s/asyncapi/store"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"
)

type ApiServer struct {
	Config *config.Config
	Logger *slog.Logger
	store  *store.Store
}

func New(conf *config.Config, logger *slog.Logger, store *store.Store) *ApiServer {
	return &ApiServer{
		Config: conf,
		Logger: logger,
		store:  store,
	}
}

func (s *ApiServer) ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func (s *ApiServer) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ping", s.ping)
	// Signup
	mux.HandleFunc("POST /auth/signup", s.signupHandler())

	middleware := NewLoggerMiddleware(s.Logger)

	server := &http.Server{
		Addr:    net.JoinHostPort(s.Config.ApiServerHost, s.Config.ApiServerPort),
		Handler: middleware(mux),
	}

	go func() {
		s.Logger.Info("apiserver running", "port", s.Config.ApiServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.Logger.Error("apiserver failed to listen and serve", "error", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			s.Logger.Error("apiserver failed to shutdown", "error", err)
		}

	}()

	wg.Wait()
	return nil
}
