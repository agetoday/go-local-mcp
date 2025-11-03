package writer

import (
	"strings"

	"github.com/kljensen/snowball"
)

type Suggestion struct {
	Type   string
	Detail string
}

func analyzeContent(text string) []Suggestion {
	var suggestions []Suggestion

	// 长句检测
	sentences := splitSentences(text)
	for _, s := range sentences {
		if wordCount(s) > 25 {
			suggestions = append(suggestions, Suggestion{
				Type:   "长句优化",
				Detail: s[:50] + "...",
			})
		}
	}

	// 关键词提取
	keywords := extractKeywords(text)
	if len(keywords) > 0 {
		suggestions = append(suggestions, Suggestion{
			Type:   "关键主题",
			Detail: strings.Join(keywords, ", "),
		})
	}

	return suggestions
}

func splitSentences(text string) []string {
	// 简易中文分句
	return strings.FieldsFunc(text, func(r rune) bool {
		return r == '。' || r == '!' || r == '?' || r == '\n'
	})
}

func wordCount(s string) int {
	return len(strings.Fields(s))
}

func extractKeywords(text string) []string {
	// 简易词干提取 (需安装 go get github.com/kljensen/snowball)
	words := strings.Fields(text)
	var keywords []string
	for _, w := range words {
		if len(w) > 3 {
			stemmed, _ := snowball.Stem(w, "english", true)
			keywords = append(keywords, stemmed)
		}
	}
	return unique(keywords)
}
func unique(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
