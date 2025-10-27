package mcpRouter

import (
	"context"
	"fmt"
	"gocode/go-local-mcp/mcpTools"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
)

func weatherHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	city, ok := req.Params.Arguments.(map[string]interface{})["city"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid city parameter")
	}
	// 模拟天气数据 - 实际应用中这里应该调用天气API
	weatherData := fmt.Sprintf("Weather in %s: Sunny, 25°C", city)
	return mcp.NewToolResultText(weatherData), nil
}

func readExcelHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {

	log.Println("readExcel", req)
	exceldir, ok := req.Params.Arguments.(map[string]interface{})["dirPath"].(string)
	if !ok {
		return nil, fmt.Errorf("readExcel invalid message parameter")
	}
	Excel_content, readErr := mcpTools.ProcessFiles(exceldir)
	if readErr != nil {
		return nil, readErr
	}
	return mcp.NewToolResultText(fmt.Sprintf("Excel content: %v", Excel_content)), nil

}
