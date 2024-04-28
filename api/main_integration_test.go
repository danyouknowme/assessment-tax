//go:build integration

package api

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/danyouknowme/assessment-tax/config"
	"github.com/danyouknowme/assessment-tax/db"
	"github.com/stretchr/testify/require"
)

const serverPort = 1337

var testServer *Server

func setup(t *testing.T) func() {
	cfg := config.New()
	dbConn, err := sql.Open("postgres", cfg.DatabaseUrl)
	require.NoError(t, err)

	err = db.PrepareDatabase(dbConn)
	require.NoError(t, err)

	testServer = NewServer(cfg, db.NewStore(dbConn))
	go func(server *Server) {
		server.Start(fmt.Sprintf(":%d", serverPort))
	}(testServer)

	return func() {
		err = db.ResetDatabase(dbConn)
		require.NoError(t, err)

		err = dbConn.Close()
		require.NoError(t, err)

		err = testServer.Shutdown(context.Background())
		require.NoError(t, err)
	}
}
