package mcpRouter

import (
	"context"
	"fmt"
	"gocode/go-local-mcp/mcpTools"
	"gocode/go-local-mcp/wordEdit"
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
func writerHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	content, ok := req.Params.Arguments.(map[string]interface{})["content"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid content parameter")
	}
	// 模拟文件写入 - 实际应用中这里应该将内容写入文件
	log.Printf("Writing content: %s", content)
	return mcp.NewToolResultText("Content written successfully"), nil
}

func createWordProcessorHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Println("createWordProcessor", req)
	docName, ok := req.Params.Arguments.(map[string]interface{})["docName"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid docName parameter")
	}
	wordDoc := wordEdit.NewDoc()
	wordDoc.AddParagraph(docName)
	wordDoc.SaveWord()
	// 模拟创建word处理器 - 实际应用中这里应该创建word处理器
	log.Printf("Creating word processor for %s", docName)
	return mcp.NewToolResultText("Word processor created successfully"), nil
}
