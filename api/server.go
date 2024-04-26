package api

import (
	"context"

	"github.com/danyouknowme/assessment-tax/config"
	"github.com/danyouknowme/assessment-tax/db"
	"github.com/labstack/echo/v4"
)

type Server struct {
	config *config.Config
	store  db.Store
	router *echo.Echo
}

func NewServer(config *config.Config, store db.Store) *Server {
	server := &Server{
		config: config,
		store:  store,
	}

	server.setupRouter()

	return server
}

func (s *Server) setupRouter() {
	e := echo.New()

	e.POST("/tax/calculations", s.CalculateTax)
	e.POST("/admin/deductions/personal", s.basicAuth(s.SettingPersonalDeduction))

	s.router = e
}

func (s *Server) Start(address string) error {
	return s.router.Start(address)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.router.Shutdown(ctx)
}
