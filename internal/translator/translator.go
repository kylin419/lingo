package translator

type Translator interface {
	TranslateText(text string, targetLangs []string) (map[string]string, error)
	DetectLanguage(text string) (string, error)
}
