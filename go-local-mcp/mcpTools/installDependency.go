package mcpTools

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// InstallPython installs Python and yt-dlp on Windows systems
// useMirror: true to use domestic mirrors (Tsinghua for Python, Aliyun for pip)
func InstallPython(useMirror bool) error {
	// 1. Download Python installer
	installerPath, err := downloadPython(useMirror)
	if err != nil {
		return fmt.Errorf("failed to download Python: %v", err)
	}

	// 2. Install Python silently
	err = installPython(installerPath)
	if err != nil {
		return fmt.Errorf("failed to install Python: %v", err)
	}

	// 3. Install yt-dlp
	err = installYtdlp(useMirror)
	if err != nil {
		return fmt.Errorf("failed to install yt-dlp: %v", err)
	}

	return nil
}

func downloadPython(useMirror bool) (string, error) {
	var pythonUrl string
	if useMirror {
		pythonUrl = "https://mirrors.tuna.tsinghua.edu.cn/python/3.11.4/python-3.11.4-amd64.exe"
	} else {
		pythonUrl = "https://www.python.org/ftp/python/3.11.4/python-3.11.4-amd64.exe"
	}
	tempDir := os.TempDir()
	installerPath := filepath.Join(tempDir, "python-installer.exe")

	fmt.Println("Downloading Python installer...")

	// Create the file
	out, err := os.Create(installerPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(pythonUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Println("Python installer downloaded successfully")
	return installerPath, nil
}

func installPython(installerPath string) error {
	fmt.Println("Installing Python...")

	cmd := exec.Command(installerPath, "/quiet", "InstallAllUsers=1", "PrependPath=1")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	fmt.Println("Python installed successfully")
	return nil
}

func installYtdlp(useMirror bool) error {
	fmt.Println("Installing yt-dlp...")

	var cmd *exec.Cmd
	if useMirror {
		cmd = exec.Command("pip", "install", "--upgrade", "yt-dlp", "-i", "https://mirrors.aliyun.com/pypi/simple/", "--trusted-host", "mirrors.aliyun.com")
	} else {
		cmd = exec.Command("pip", "install", "--upgrade", "yt-dlp")
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	fmt.Println("yt-dlp installed successfully")
	return nil
}
