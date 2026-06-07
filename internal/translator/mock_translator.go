package translator

import "fmt"

// MockTranslator 開發測試用的假翻譯器
type MockTranslator struct{}

// NewMockTranslator 初始化 MockTranslator
func NewMockTranslator() Translator {
	return &MockTranslator{}
}

// TranslateText 模擬翻譯行為，直接在文字後方加上標籤
func (m *MockTranslator) TranslateText(text string, targetLangs []string) (map[string]string, error) {
	result := make(map[string]string)
	
	for _, lang := range targetLangs {
		result[lang] = fmt.Sprintf("%s (模擬 %s 翻譯結果)", text, lang)
	}
	
	return result, nil
}