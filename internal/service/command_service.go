package service

import (
	"fmt"
	"strings"

	"line-translate-bot/internal/repository"
)

var supportedLangs = map[string]string{
	"zh-Hant": "繁體中文",
	"zh-Hans": "簡體中文",
	"en":      "英文",
	"ja":      "日文",
	"ko":      "韓文",
	"th":      "泰文",
	"vi":      "越南文",
	"es":      "西班牙文",
	"fr":      "法文",
	"de":      "德文",
}

type CommandService struct {
	repo repository.GroupRepository
}

func NewCommandService(repo repository.GroupRepository) *CommandService {
	return &CommandService{repo: repo}
}

func (s *CommandService) HandleSetCommand(sourceID, text string) string {
	argsStr := strings.TrimSpace(strings.TrimPrefix(text, "$set"))
	if argsStr == "" {
		return generateHelpMessage("請提供語言代碼！")
	}

	// 使用 strings.Fields 可以自動忽略多餘的空白字元
	langs := strings.Fields(argsStr)

	var validLangs []string
	var invalidLangs []string

	// 檢查使用者輸入的代碼是否存在於我們的字典中
	for _, l := range langs {
		if _, exists := supportedLangs[l]; exists {
			validLangs = append(validLangs, l)
		} else {
			invalidLangs = append(invalidLangs, l)
		}
	}

	// 如果有任何錯誤的代碼，直接拒絕設定並回傳教學
	if len(invalidLangs) > 0 {
		errMsg := fmt.Sprintf("❌ 找不到以下語言代碼：%s", strings.Join(invalidLangs, ", "))
		return generateHelpMessage(errMsg)
	}

	// 儲存合法的設定
	err := s.repo.SaveLanguages(sourceID, validLangs)
	if err != nil {
		return "設定失敗，請稍後再試。"
	}

	// 組合成功提示訊息
	var langNames []string
	for _, l := range validLangs {
		langNames = append(langNames, supportedLangs[l])
	}

	return fmt.Sprintf("✅ 語言已成功設定為：%s\n接下來的訊息我會在這幾個語言間互翻！", strings.Join(langNames, ", "))
}

// generateHelpMessage 負責產生包含可用代碼表的提示訊息
func generateHelpMessage(prefixMsg string) string {
	var helpBuilder strings.Builder
	helpBuilder.WriteString(prefixMsg + "\n\n")
	helpBuilder.WriteString("支援的常用語言代碼如下：\n")

	// 列出所有支援的語言
	for code, name := range supportedLangs {
		helpBuilder.WriteString(fmt.Sprintf("• %s (%s)\n", code, name))
	}
	return helpBuilder.String()
}
