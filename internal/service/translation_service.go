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

func (s *TranslationService) TranslateMessage(sourceID string, text string) (string, error) {
	langs, err := s.repo.GetLanguages(sourceID)
	if err != nil || len(langs) == 0 {
		return "", nil
	}
	results, err := s.trans.TranslateText(text, langs)
	if err != nil {
		return "", fmt.Errorf("翻譯服務發生錯誤: %v", err)
	}
	var replies []string
	for _, lang := range langs {
		if translatedText, ok := results[lang]; ok {
			replies = append(replies, fmt.Sprintf("[%s] %s", lang, translatedText))
		}
	}

	if len(replies) == 0 {
		return "", nil
	}
	return strings.Join(replies, "\n\n"), nil
}
