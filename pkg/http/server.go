package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

type Server struct {
	Server *http.Server
	Router *Router
}

func NewServer(router *Router, port string) *Server {
	s := &Server{
		Router: router,
	}

	s.Server = &http.Server{
		Addr:    port,
		Handler: s.Router.Engine,
	}

	return s
}

func (s *Server) Init() error {
	go func() {
		if err := s.Server.ListenAndServe(); err != nil {
			logrus.Infof("Server is up and listening on port: %s", s.Server.Addr)
		}
	}()
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Server.Shutdown(ctx); err != nil {
		logrus.Errorf("Server forced to shutdown: %v", err)
		return err
	}

	logrus.Info("Server exiting ")
	return nil
}
