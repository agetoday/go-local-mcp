package writer

import (
	"fmt"
	"log"
	"github.com/fsnotify/fsnotify"
)

func SetupWatcher(dir string) {
	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Printf("📝 文件变更: %s\n", event.Name)
					go processFile(event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("监控错误:", err)
			}
		}
	}()

	err := watcher.Add(dir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
