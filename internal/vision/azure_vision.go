package vision

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"
)

type AzureVision struct {
	endpoint string
	apiKey   string
	client   *http.Client
}

// 定義結構以擷取文字與座標
type Line struct {
	Text        string    `json:"text"`
	BoundingBox []float64 `json:"boundingBox"`
}

type readResultResponse struct {
	Status        string `json:"status"`
	AnalyzeResult struct {
		ReadResults []struct {
			Lines []Line `json:"lines"`
		} `json:"readResults"`
	} `json:"analyzeResult"`
}

func NewAzureVision(endpoint, apiKey string) VisionService {
	return &AzureVision{
		endpoint: endpoint,
		apiKey:   apiKey,
		client:   &http.Client{Timeout: 30 * time.Second},
	}
}

func (a *AzureVision) ExtractTextFromImage(imageBytes []byte) (string, error) {
	log.Printf("Debug: 準備發送圖片至 Azure，大小: %d bytes", len(imageBytes))

	req, _ := http.NewRequest("POST", a.endpoint+"/vision/v3.2/read/analyze", bytes.NewBuffer(imageBytes))
	req.Header.Set("Ocp-Apim-Subscription-Key", a.apiKey)
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := a.client.Do(req)
	if err != nil || (resp.StatusCode != 202 && resp.StatusCode != 200) {
		return "", fmt.Errorf("上傳失敗，狀態碼: %d", resp.StatusCode)
	}
	operationURL := resp.Header.Get("Operation-Location")
	resp.Body.Close()

	// 輪詢處理結果
	var result readResultResponse
	found := false
	for i := 0; i < 15; i++ {
		time.Sleep(1 * time.Second)
		req, _ := http.NewRequest("GET", operationURL, nil)
		req.Header.Set("Ocp-Apim-Subscription-Key", a.apiKey)
		resp, err := a.client.Do(req)
		if err != nil {
			continue
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		if result.Status == "succeeded" {
			found = true
			break
		}
		if result.Status == "failed" {
			return "", fmt.Errorf("Azure OCR 辨識失敗")
		}
	}

	if !found {
		return "", fmt.Errorf("OCR 處理超時")
	}

	// 整理並排序所有文字區塊
	var allLines []Line
	for _, read := range result.AnalyzeResult.ReadResults {
		allLines = append(allLines, read.Lines...)
	}

	// 排序邏輯：由上而下，同高度時由右至左 (針對漫畫/日文排版優化)
	sort.Slice(allLines, func(i, j int) bool {
		if len(allLines[i].BoundingBox) < 2 || len(allLines[j].BoundingBox) < 2 {
			return false
		}
		// Y 軸差距小於 20 像素視為同一行
		if abs(allLines[i].BoundingBox[1]-allLines[j].BoundingBox[1]) < 20 {
			return allLines[i].BoundingBox[0] > allLines[j].BoundingBox[0] // X 大的在右邊
		}
		return allLines[i].BoundingBox[1] < allLines[j].BoundingBox[1]
	})

	var sb strings.Builder
	for _, line := range allLines {
		// 雜訊過濾：長度過短或純特殊符號
		if len(line.Text) < 3 {
			continue
		}
		matched, _ := regexp.MatchString(`^[0-9\s!@#$%^&*()]+$`, line.Text)
		if matched {
			continue
		}

		sb.WriteString(line.Text)
		sb.WriteString("\n")
	}

	return strings.TrimSpace(sb.String()), nil
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
