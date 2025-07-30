package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/lonegunmanb/terraform-mcp-eva/pkg"
	"github.com/matt-FFFFFF/tfpluginschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	}))

	mode := flag.String("mode", getenv("TRANSPORT_MODE", "stdio"), "transport mode, can be `stdio` or `streamable-http`")
	host := flag.String("host", getenv("TRANSPORT_HOST", "127.0.0.1"), "host for streamable-http server")
	port := flag.String("port", getenv("TRANSPORT_PORT", "8080"), "port for streamable-http server")
	flag.Parse()

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "mcp-ever",
		Version: "0.1.0",
		Title:   "Terraform provider MCP Server",
	}, nil)

	pkg.RegisterMcpServer(server)

	providerSchemaServer := tfpluginschema.NewServer(nil)

	switch *mode {
	case "stdio":
		ctx := context.Background()
		ctx = context.WithValue(ctx, tfpluginschema.ContextKey{}, providerSchemaServer)
		if err := server.Run(ctx, mcp.NewStdioTransport()); err != nil {
			l.Error(err.Error())
		}
	case "streamable-http":
		addr := fmt.Sprintf("%s:%s", *host, *port)
		l.Info("MCP server serving", "address", addr)
		handler := mcp.NewSSEHandler(func(request *http.Request) *mcp.Server {
			// Add context with dependencies to the request
			ctxWithDeps := context.WithValue(request.Context(), tfpluginschema.ContextKey{}, providerSchemaServer)
			request = request.WithContext(ctxWithDeps)
			return server
		})
		if err := http.ListenAndServe(addr, handler); err != nil {
			l.Error(err.Error())
		}
	default:
		l.Error("unknown mode", "mode", *mode)
		os.Exit(1)
	}
}

func getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
