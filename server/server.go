package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/shaileshhb/mcp-tutorial/security"
)

// NewHttpServer
func NewHttpServer(server *mcp.Server, mcpHandler *mcp.StreamableHTTPHandler) *mux.Router {
	router := mux.NewRouter()
	router.Use(security.RecoveryMiddleware)

	router.HandleFunc("/health", func(
		w http.ResponseWriter,
		req *http.Request,
	) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}).Methods(http.MethodGet)

	router.HandleFunc("/mcp", func(
		w http.ResponseWriter,
		req *http.Request,
	) {
		mcpHandler.ServeHTTP(w, req)
	}).Methods(http.MethodPost, http.MethodGet)

	return router
}

func NewMcpServer(logger *slog.Logger) *mcp.Server {
	implementation := &mcp.Implementation{
		Title:   "MCP Tutorial For Currency Conversion",
		Name:    "currency-conversion-mcp",
		Version: "1.0.0",
	}

	serverOptions := &mcp.ServerOptions{
		Logger:    logger,
		KeepAlive: 5 * time.Minute,
	}

	return mcp.NewServer(implementation, serverOptions)
}
