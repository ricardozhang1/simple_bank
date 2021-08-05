package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/techschool/simplebank/db/sqlc"
	"os"
	"testing"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	server := NewServer(store)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

