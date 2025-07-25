package web

import "github.com/mark3labs/mcp-go/server"

func NewMCPServer() *server.MCPServer {
	return server.NewMCPServer(
		"mcpcities",
		"0.0.1",
		server.WithToolCapabilities(false),
	)
}
