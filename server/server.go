package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// NewHttpServer
func NewHttpServer(server *mcp.Server, mcpHandler *mcp.StreamableHTTPHandler) *mux.Router {
	router := mux.NewRouter()
	// router.Use(security.RecoveryMiddleware)
	// router.Use(security.HandleCors)

	router.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}).Methods(http.MethodGet)

	// Optional root route to avoid confusion when testing in browser.
	router.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("MCP server is running. Use /mcp endpoint."))
	}).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)

	router.HandleFunc("/mcp", func(w http.ResponseWriter, req *http.Request) {
		fmt.Printf(
			"method=%s path=%s auth=%q\n",
			req.Method,
			req.URL.Path,
			req.Header.Get("Authorization"),
		)
		mcpHandler.ServeHTTP(w, req)
	}).Methods(http.MethodPost, http.MethodGet, http.MethodOptions)

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
