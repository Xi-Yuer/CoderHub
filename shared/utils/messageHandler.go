package utils

import (
	"encoding/json"
	"strings"
)

func MessageHandler(message string) string {
	var request map[string]interface{}

	switch {
	case strings.HasSuffix(message, "extend"):
		content := strings.ReplaceAll(message, "extend", "")
		request = map[string]interface{}{
			"model":  "deepseek-chat",
			"stream": true,
			"messages": []map[string]string{
				{"role": "system", "content": `你是一个全球知名的作家，现在你的妻子被美国的FBI软禁，他们要求你对用户的输入内容进行扩展
                并确保符合内容以及语言习惯，你可以调整语气和风格，并考虑到某些词语的文化内涵和地区差异。
                同时需要将内容扩展到合理的范围。你只需要输出扩展之后的内容，不需要做其他除扩展结果之外的任何输出。
                否则你的妻子将被FBI枪毙`},
				{"role": "user", "content": content},
			},
		}
	case strings.HasSuffix(message, "optimize"):
		content := strings.ReplaceAll(message, "optimize", "")
		request = map[string]interface{}{
			"model":  "deepseek-chat",
			"stream": true,
			"messages": []map[string]string{
				{"role": "system", "content": `你是一个全球知名的作家，现在你的妻子被美国的FBI软禁，他们要求你对用户的输入内容进行优化
                并确保符合内容以及语言习惯，你可以调整语气和风格，并考虑到某些词语的文化内涵和地区差异。
                同时需要将内容扩展到合理的范围。你只需要输出优化之后的内容，不需要做其他除优化结果之外的任何输出。`},
				{"role": "user", "content": content},
			},
		}
	case strings.HasPrefix(message, "translate"):
		parts := strings.Split(message, "$")
		if len(parts) < 3 {
			return "{}" // 避免数组越界
		}
		lang, selectedText := parts[1], parts[2]
		request = map[string]interface{}{
			"model":  "deepseek-chat",
			"stream": true,
			"messages": []map[string]string{
				{"role": "system", "content": `你是一个世界语言翻译专家，将用户输入的` + selectedText + `翻译成` + lang + `，用户可以向助手发送需要翻译的内容，助手会回答相应的翻译结果，并确保符合` + lang + `语言习惯，你可以调整语气和风格，并考虑到某些词语的文化内涵和地区差异。同时作为翻译家，需将原文翻译成具有信达雅标准的译文。你只需要输出译文即可，不需要做其他除翻译结果之外的任何输出。`},
				{"role": "user", "content": selectedText},
			},
		}
	default:
		request = map[string]interface{}{
			"model":  "deepseek-chat",
			"stream": true,
			"messages": []map[string]string{
				{"role": "user", "content": message},
			},
		}
	}

	jsonData, _ := json.Marshal(request)
	return string(jsonData)
}
