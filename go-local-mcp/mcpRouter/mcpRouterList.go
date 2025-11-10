package mcpRouter

import (
	"context"
	"fmt"
	"gocode/go-local-mcp/mcpTools"
	"gocode/go-local-mcp/wordEdit"
	"log"
	"os/exec"

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
func convertTo2DStringArray(data interface{}) ([][]string, error) {
	// Handle case where data is already [][]string
	if arr, ok := data.([][]string); ok {
		return arr, nil
	}

	// Handle case where data is []interface{} containing []interface{}
	if arr, ok := data.([]interface{}); ok {
		result := make([][]string, len(arr))
		for i, row := range arr {
			if rowArr, ok := row.([]interface{}); ok {
				strRow := make([]string, len(rowArr))
				for j, cell := range rowArr {
					strRow[j] = fmt.Sprintf("%v", cell)
				}
				result[i] = strRow
			} else {
				return nil, fmt.Errorf("invalid row format at index %d", i)
			}
		}
		return result, nil
	}

	return nil, fmt.Errorf("unsupported data format")
}

func createExcelHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	fmt.Println("createExcel", req)
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid arguments format")
	}

	data, exists := args["content_array"]
	if !exists {
		return nil, fmt.Errorf("content_array parameter is required")
	}

	exceldata, err := convertTo2DStringArray(data)
	if err != nil {
		return nil, fmt.Errorf("createExcel invalid message parameter: %v", err)
	}

	fmt.Println("==============", exceldata)
	execErr := mcpTools.WriteExcelFile(exceldata)
	if execErr != nil {
		return nil, execErr
	}
	return mcp.NewToolResultText(fmt.Sprintf("Excel content: %v", exceldata)), nil
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

func videoDownloadHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Println("videoDownload", req)
	args := req.Params.Arguments.(map[string]interface{})

	url, ok := args["url"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid url parameter")
	}

	platform, ok := args["platform"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid platform parameter")
	}

	resolution := "1080" // default to 1080p
	if res, ok := args["resolution"].(string); ok {
		resolution = res
	}

	// Check if yt-dlp is installed with detailed installation guide
	if _, err := exec.LookPath("yt-dlp"); err != nil {
		guide := `
yt-dlp is required for video downloads. Please install it first:

Windows:
1. Install Python from python.org
2. Run: pip install yt-dlp
3. Add Python to PATH during installation

macOS:
1. Install Homebrew: /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
2. Run: brew install yt-dlp

Linux:
1. Run: sudo apt install python3-pip
2. Run: pip install yt-dlp
`
		return nil, fmt.Errorf("yt-dlp is not installed.\n%s", guide)
	}

	downloader := mcpTools.NewVideoDownloader("downloads")
	filePath, err := downloader.Download(url, mcpTools.VideoPlatform(platform), resolution)
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("video download failed (exit code %d): %s\n%s", 
				exitErr.ExitCode(), 
				exitErr.Stderr, 
				"Please ensure yt-dlp is properly installed and in your PATH")
		}
		return nil, fmt.Errorf("video download failed: %v", err)
	}

	return mcp.NewToolResultText(fmt.Sprintf("Video downloaded successfully: %s", filePath)), nil
}
