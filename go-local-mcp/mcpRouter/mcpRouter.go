package mcpRouter

import (
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type MCPServer struct {
	server *server.MCPServer
}

func NewMCPServer() *MCPServer {
	mcpServer := server.NewMCPServer(
		"in-mcp-server",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	weather := mcp.NewTool("weather",
		mcp.WithDescription("Get weather information"),
		mcp.WithString("city", mcp.Required(), mcp.Description("City name to query weather")),
	)
	read_Excel := mcp.NewTool("readExcel",
		mcp.WithDescription("Read Excel file"),
		mcp.WithString("dirPath", mcp.Required(), mcp.Description("Directory path to Excel file")),
	)

	mcpServer.AddTool(weather, weatherHandler)
	mcpServer.AddTool(read_Excel, readExcelHandler)

	return &MCPServer{server: mcpServer}
}

func (s *MCPServer) ServeSSE(addr string) *server.SSEServer {
	return server.NewSSEServer(s.server, server.WithBaseURL(fmt.Sprintf("http://%s", addr)))
}
