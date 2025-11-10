package main

import "gocode/go-local-mcp/mcpTools"

func main() {
	downloader := mcpTools.NewVideoDownloader("downloads")
	downloader.Download("https://www.instagram.com/reel/DNGTLMkzNT_/?igsh=MWd0czFnMnBjemxzbQ==", mcpTools.Instagram, "1080")
}
