package writer

import (
	"flag"
	"fmt"
	"log"
)

func RunWriter(content string) {
	watchDir := flag.String("dir", "./writing", "监控目录路径")
	lang := flag.String("lang", "zh", "默认语言")
	flag.Parse()

	fmt.Printf("🚀 启动MCP写作助手 (语言: %s)\n", *lang)

	// 初始化语言处理模块
	initLanguageProcessor(*lang)
	// 启动文件监控
	SetupWatcher(*watchDir)
}

func initLanguageProcessor(lang string) {
	switch lang {
	case "zh":
		log.Println("加载中文处理模块")
	case "en":
		log.Println("Loading English processor")
	default:
		log.Println("使用默认处理器")
	}
}
