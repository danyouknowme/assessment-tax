package api

import (
	"context"
	"database/sql"

	"github.com/danyouknowme/assessment-tax/config"
	"github.com/labstack/echo/v4"
)

type Server struct {
	config *config.Config
	db     *sql.DB
	router *echo.Echo
}

func NewServer(config *config.Config, db *sql.DB) *Server {
	server := &Server{
		config: config,
		db:     db,
	}

	server.setupRouter()

	return server
}

func (server *Server) setupRouter() {
	e := echo.New()

	e.POST("/tax/calculations", CalculateTax)

	server.router = e
}

func (server *Server) Start(address string) error {
	return server.router.Start(address)
}

func (server *Server) Shutdown(ctx context.Context) error {
	return server.router.Shutdown(ctx)
}
