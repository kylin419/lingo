package translator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type AzureTranslator struct {
	apiKey string
	region string
	client *http.Client
}

type azureRequestItem struct {
	Text string `json:"Text"`
}

type azureResponseItem struct {
	Translations []struct {
		Text string `json:"text"`
		To   string `json:"to"`
	} `json:"translations"`
}

type detectResponse []struct {
	Language string  `json:"language"`
	Score    float64 `json:"score"`
}

func NewAzureTranslator(apiKey, region string) Translator {
	return &AzureTranslator{
		apiKey: apiKey,
		region: region,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (a *AzureTranslator) TranslateText(text string, targetLangs []string) (map[string]string, error) {
	if len(targetLangs) == 0 {
		return nil, nil
	}

	baseURL, err := url.Parse("https://api.cognitive.microsofttranslator.com/translate?api-version=3.0")
	if err != nil {
		return nil, err
	}

	q := baseURL.Query()
	for _, lang := range targetLangs {
		q.Add("to", lang)
	}
	baseURL.RawQuery = q.Encode()

	reqBody := []azureRequestItem{{Text: text}}
	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", baseURL.String(), bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", a.apiKey)
	req.Header.Set("Content-Type", "application/json")
	if a.region != "" && a.region != "global" {
		req.Header.Set("Ocp-Apim-Subscription-Region", a.region)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Azure API 回傳錯誤狀態碼 %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var azureResp []azureResponseItem
	if err := json.NewDecoder(resp.Body).Decode(&azureResp); err != nil {
		return nil, err
	}

	results := make(map[string]string)
	if len(azureResp) > 0 {
		for _, t := range azureResp[0].Translations {
			results[t.To] = t.Text
		}
	}

	return results, nil
}
func (a *AzureTranslator) DetectLanguage(text string) (string, error) {
	// 設定端點 URL
	endpoint := "https://api.cognitive.microsofttranslator.com/detect?api-version=3.0"

	// 準備請求 Body
	reqBody := []azureRequestItem{{Text: text}}
	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// 建立 POST 請求
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", err
	}

	//  設定 Header (關鍵：一定要包含 API Key 和 Region)
	req.Header.Set("Ocp-Apim-Subscription-Key", a.apiKey)
	req.Header.Set("Content-Type", "application/json")
	if a.region != "" && a.region != "global" {
		req.Header.Set("Ocp-Apim-Subscription-Region", a.region)
	}

	// 5. 發送請求
	resp, err := a.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 6. 解析回應
	var detectResp detectResponse
	if err := json.NewDecoder(resp.Body).Decode(&detectResp); err != nil {
		return "", err
	}

	// 7. 取出第一個結果的語言代碼
	if len(detectResp) > 0 {
		return detectResp[0].Language, nil
	}
	return "", fmt.Errorf("無法偵測到語言")
}
