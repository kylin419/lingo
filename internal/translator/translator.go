package translator

import "strings"

type Translator interface {
	TranslateText(text string, targetLangs []string) (map[string]string, error)
	DetectLanguage(text string) (string, error)
}

func GetLangFamily(lang string) string {
	// 統一將所有中文系變體歸類為 'zh'
	if strings.HasPrefix(lang, "zh") {
		return "zh"
	}
	// 你可以依此類推，將 'en-US', 'en-GB' 都歸類為 'en'
	if strings.HasPrefix(lang, "en") {
		return "en"
	}
	return lang
}
