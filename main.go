package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/lonegunmanb/terraform-mcp-eva/pkg"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	mode := flag.String("mode", getenv("TRANSPORT_MODE", "stdio"), "transport mode, can be `stdio` or `streamable-http`")
	host := flag.String("host", getenv("TRANSPORT_HOST", "127.0.0.1"), "host for streamable-http server")
	port := flag.String("port", getenv("TRANSPORT_PORT", "8080"), "port for streamable-http server")
	flag.Parse()

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "mcp-ever",
		Version: "0.1.0",
	}, nil)
	pkg.RegisterMcpServer(server)

	switch *mode {
	case "stdio":
		if err := server.Run(context.Background(), mcp.NewStdioTransport()); err != nil {
			log.Fatal(err)
		}
	case "streamable-http":
		addr := fmt.Sprintf("%s:%s", *host, *port)
		log.Printf("MCP server serving at %s", addr)
		handler := mcp.NewSSEHandler(func(request *http.Request) *mcp.Server {
			return server
		})
		if err := http.ListenAndServe(addr, handler); err != nil {
			log.Fatalf("failed to start streamable-http server: %v", err)
		}
	default:
		log.Fatalf("unknown mode: %s", *mode)
	}
}

func getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
