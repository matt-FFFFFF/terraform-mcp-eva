package main

import (
	"context"
	"flag"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"log"
	"terraform-mcp-eva/pkg"
)

var httpAddr = flag.String("http", "", "if set, use streamable HTTP at this address, instead of stdin/stdout")

func main() {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "mcp-ever",
		Version: "0.1.0",
	}, nil)
	pkg.RegisterMcpServer(server)
	if err := server.Run(context.Background(), mcp.NewStdioTransport()); err != nil {
		log.Fatal(err)
	}
}
