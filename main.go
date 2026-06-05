package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/shaileshhb/mcp-tutorial/server"
	"github.com/shaileshhb/mcp-tutorial/tool"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	mcpServer := server.NewMcpServer(logger)

	mcpHandler := mcp.NewStreamableHTTPHandler(func(request *http.Request) *mcp.Server {
		return mcpServer
	}, &mcp.StreamableHTTPOptions{
		Logger:       logger,
		JSONResponse: true,
	})

	router := server.NewHttpServer(mcpServer, mcpHandler)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Title:       "Greeting",
		Name:        "greeting",
		Description: "Greet a person by name",
	}, tool.SayHi)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Title:       "Currency Conversion",
		Name:        "currency-conversion",
		Description: "Convert a currency to another currency",
	}, tool.ConvertCurrencyTool)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadTimeout:       5 * time.Minute,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      5 * time.Minute,
		IdleTimeout:       5 * time.Minute,
		MaxHeaderBytes:    1 << 20,
	}

	go func() {
		logger.Info("starting server", "addr", ":8080")

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {

			logger.Error("server failed", "err", err)
			os.Exit(1)
		}
	}()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	<-ctx.Done()

	logger.Info("shutdown initiated")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown failed", "err", err)
	}

	logger.Info("shutdown complete")
}
