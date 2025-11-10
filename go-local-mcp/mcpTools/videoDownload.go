package mcpTools

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"time"
)

type VideoPlatform string

const (
	YouTube   VideoPlatform = "youtube"
	Instagram VideoPlatform = "instagram"
)

type VideoDownloader struct {
	outputDir string
}

func NewVideoDownloader(outputDir string) *VideoDownloader {
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, 0755)
	}
	return &VideoDownloader{outputDir: outputDir}
}

func DetectExecuter() {
	if _, err := exec.LookPath("yt-dlp"); err!= nil {
		fmt.Println("请安装 yt-dlp 并配置环境变量")
	}

}
func InstallExecuter() {
	// 安装依赖 - python 环境
	if _, err := exec.LookPath("pip"); err!= nil {
		fmt.Println("请安装 python 环境")
		//  开始Windows 下安装python 环境
		startInstallPython := exec.Command("cmd", "/C", "start", "https://www.python.org/downloads/")
		startInstallPython.Stdout = os.Stdout
		startInstallPython.Stderr = os.Stderr
		err := startInstallPython.Run()
		if err != nil {
			log.Fatalf("打开浏览器失败: %v", err)
		}
		fmt.Println("请安装 python 环境")

		return
	}

	cmd := exec.Command("pip", "install", "yt-dlp")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("yt-dlp 安装失败: %v", err)
	}
	fmt.Println("yt-dlp 安装成功")	
}


func (vd *VideoDownloader) Download(url string, platform VideoPlatform, resolution string) (string, error) {
	var cmd *exec.Cmd

	switch platform {
	case YouTube:
		cmd = exec.Command("yt-dlp",
			"--no-warnings",
			"--quiet",
			"-f", fmt.Sprintf("bestvideo[height<=%s]+bestaudio/best[height<=%s]", resolution, resolution),
			"-o", filepath.Join(vd.outputDir, "%(title)s.%(ext)s"),
			"--",
			strings.TrimSpace(url))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	case Instagram:
		cmd = exec.Command("yt-dlp",
			"--no-warnings",
			"--quiet",
			"-f", "best",
			"-o", filepath.Join(vd.outputDir, "%(title)s.%(ext)s"),
			"--",
			strings.TrimSpace(url))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	default:
		return "", fmt.Errorf("unsupported platform: %s", platform)
	}
	err := cmd.Run()
	if err != nil {
		log.Fatalf("命令执行失败: %v", err)
	}
	fmt.Println("执行过程：\n",cmd.Stdout)
	fmt.Println("命令执行完毕")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("download failed: %v\nCommand output: %s", err, string(output))
	}

	// Find the most recently modified file in the output directory
	var newestFile string
	var newestTime time.Time

	files, err := os.ReadDir(vd.outputDir)
	if err != nil {
		return "", fmt.Errorf("failed to read output directory: %v", err)
	}

	for _, file := range files {
		info, err := file.Info()
	if err != nil {
			continue
		}

		if info.ModTime().After(newestTime) {
			newestTime = info.ModTime()
			newestFile = filepath.Join(vd.outputDir, file.Name())
		}
	}

	if newestFile == "" {
		return "", fmt.Errorf("no downloaded file found in %s", vd.outputDir)
	}

	return newestFile, nil
}
