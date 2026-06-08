package service

import (
	"fmt"
	"strings"

	"line-translate-bot/internal/repository"
	"line-translate-bot/internal/translator"
)

type TranslationService struct {
	repo  repository.GroupRepository
	trans translator.Translator
}

func NewTranslationService(repo repository.GroupRepository, trans translator.Translator) *TranslationService {
	return &TranslationService{
		repo:  repo,
		trans: trans,
	}
}
func (s *TranslationService) HandleMessage(groupID string, text string) (string, error) {
	sourceLang, _ := s.trans.DetectLanguage(text)
	targetLangs, _ := s.repo.GetLanguages(groupID)

	var finalResponse string
	for _, lang := range targetLangs {
		// 如果是原文語言，就跳過，不重複翻譯
		if lang == sourceLang {
			continue
		}

		// 翻譯成目標語言
		result, _ := s.trans.TranslateText(text, []string{lang})
		finalResponse += fmt.Sprintf("[%s]: %s\n", lang, result[lang])
	}
	return finalResponse, nil
}

func (s *TranslationService) TranslateMessage(sourceID string, text string) (string, error) {
	// 偵測原始語言
	sourceLang, err := s.trans.DetectLanguage(text)
	if err != nil {
		sourceLang = "" // 如果偵測失敗，將原始語言設為空字串，讓後續過濾邏輯能正常運作
	}

	// 取出目標語言設定
	langs, err := s.repo.GetLanguages(sourceID)
	if err != nil || len(langs) == 0 {
		return "⚠️ 請注意，您還沒設定翻譯語言喔！請輸入 `$set [語言1] [語言2]` 來進行設定。", nil
	}

	// 執行翻譯
	results, err := s.trans.TranslateText(text, langs)
	if err != nil {
		return "", fmt.Errorf("翻譯服務發生錯誤: %v", err)
	}

	// 組裝回傳訊息，過濾掉來源語言
	var replies []string
	for _, lang := range langs {
		// 如果目標語言與原始語言相同，就跳過
		if lang == sourceLang {
			continue
		}

		if translatedText, ok := results[lang]; ok {
			replies = append(replies, fmt.Sprintf("[%s] %s", lang, translatedText))
		}
	}

	if len(replies) == 0 {
		return "", nil
	}
	return strings.Join(replies, "\n\n"), nil
}
