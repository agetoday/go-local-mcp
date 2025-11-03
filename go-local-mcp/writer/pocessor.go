package writer

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type Document struct {
	Meta    map[string]interface{}
	Content string
}

func processFile(path string) {
	doc := readMarkdown(path)

	// 自动更新修改时间
	if doc.Meta["modified"] == nil {
		doc.Meta["modified"] = time.Now().Format(time.RFC3339)
	}

	// 生成写作建议
	suggestions := analyzeContent(doc.Content)

	// 保存处理结果
	suggestionStrings := make([]string, len(suggestions))
	for i, suggestion := range suggestions {
		suggestionStrings[i] = fmt.Sprintf("Suggestion: %v", suggestion) // 根据Suggestion的字段格式化字符串
	}

	saveSuggestions(path+".suggestions", suggestionStrings)

	writeMarkdown(path, doc)
}

func readMarkdown(path string) Document {
	file, _ := os.Open(path)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	doc := Document{Meta: make(map[string]interface{})}

	inMeta := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "---") {
			inMeta = !inMeta
			continue
		}
		if inMeta {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				doc.Meta[strings.TrimSpace(parts[0])] = strings.Trim(parts[1], ` "`)
			}
		} else {
			doc.Content += line + "\n"
		}
	}
	return doc
}
func saveSuggestions(path string, suggestions []string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, suggestion := range suggestions {
		_, err = fmt.Fprintln(file, suggestion)
		if err != nil {
			return err
		}
	}
	return nil
}
func writeMarkdown(path string, doc Document) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入元数据
	for key, value := range doc.Meta {
		_, err := fmt.Fprintf(file, "%s: %v\n", key, value)
		if err != nil {
			return err
		}
	}
	// 写入空行以分隔元数据和内容
	_, err = file.WriteString("\n")
	if err != nil {
		return err
	}

	// 写入文档内容
	_, err = file.WriteString(doc.Content)
	if err != nil {
		return err
	}

	return nil
}
