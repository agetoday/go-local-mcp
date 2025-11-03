package mcpRouter

import (
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type MCPServer struct {
	server *server.MCPServer
}

func MakeServer() *server.MCPServer {
	mcpServer := server.NewMCPServer(
		"in-mcp-server",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	return mcpServer
}

func NewMCPServer() *MCPServer {
	mcpServer := MakeServer()

	weather := mcp.NewTool("weather",
		mcp.WithDescription("Get weather information"),
		mcp.WithString("city", mcp.Required(), mcp.Description("City name to query weather")),
	)
	read_Excel := mcp.NewTool("readExcel",
		mcp.WithDescription("Read Excel file"),
		mcp.WithString("dirPath", mcp.Required(), mcp.Description("Directory path to Excel file")),
	)
	writer_tool := mcp.NewTool("writer",
		mcp.WithDescription("Analyze and suggest improvements for text content"),
		mcp.WithString("content", mcp.Required(), mcp.Description("Text content to analyze")),
	)

	create_wordprocessor := mcp.NewTool("createWordProcessor",
		mcp.WithDescription("Create a Word document with specified content"),
		mcp.WithString("content", mcp.Required(), mcp.Description("Content to include in the Word document")),
	)

	mcpServer.AddTool(weather, weatherHandler)
	mcpServer.AddTool(read_Excel, readExcelHandler)
	mcpServer.AddTool(writer_tool, writerHandler)
	mcpServer.AddTool(create_wordprocessor, createWordProcessorHandler)
	return &MCPServer{server: mcpServer}
}

func (s *MCPServer) ServeSSE(addr string) *server.SSEServer {
	return server.NewSSEServer(s.server, server.WithBaseURL(fmt.Sprintf("http://%s", addr)))
}

func (s *MCPServer) ServeStdio() *server.StdioServer {
	return server.NewStdioServer(s.server)
}
