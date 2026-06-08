package handler

import (
	"io"
	"line-translate-bot/internal/service"
	"line-translate-bot/internal/vision"
	"log"
	"net/http"
	"strings"

	"github.com/line/line-bot-sdk-go/v8/linebot"
)

type LineHandler struct {
	bot          *linebot.Client
	cmdService   *service.CommandService
	transService *service.TranslationService
	visionSvc    vision.VisionService
}

func NewLineHandler(bot *linebot.Client, cmdService *service.CommandService, transService *service.TranslationService, visionSvc vision.VisionService) *LineHandler {
	return &LineHandler{
		bot:          bot,
		cmdService:   cmdService,
		transService: transService,
		visionSvc:    visionSvc,
	}
}

func (h *LineHandler) HandleRequest(w http.ResponseWriter, req *http.Request) {
	events, err := h.bot.ParseRequest(req)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	for _, event := range events {
		sourceID := getSourceID(event.Source)

		switch event.Type {
		// 處理加入群組事件
		case linebot.EventTypeJoin:
			welcomeMsg := "大家好！我是翻譯機器人。\n請輸入 `$set [語言1] [語言2]...` 來設定要互翻的語言。\n例如：`$set zh-Hant en ja`"
			if _, err := h.bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(welcomeMsg)).Do(); err != nil {
				log.Println("回覆歡迎訊息失敗:", err)
			}

		// 處理文字訊息事件
		case linebot.EventTypeMessage:
			switch msg := event.Message.(type) {
			case *linebot.TextMessage:
				if message, ok := event.Message.(*linebot.TextMessage); ok {
					text := strings.TrimSpace(message.Text)

					// 判斷是否為 $set 指令
					if strings.HasPrefix(text, "$set ") {
						replyText := h.cmdService.HandleSetCommand(sourceID, text)
						h.bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyText)).Do()
					} else {
						// 執行翻譯
						replyText, err := h.transService.TranslateMessage(sourceID, text)
						log.Println("翻譯結果:", replyText)
						if err != nil {
							log.Println("翻譯處理失敗:", err)
							continue
						}
						// 有翻譯結果才回傳
						if replyText != "" {
							h.bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyText)).Do()
						}
					}
				}
			case *linebot.ImageMessage:
				response, err := h.bot.GetMessageContent(msg.ID).Do()

				if err != nil {
					log.Println("圖片下載失敗:", err)
					continue
				}
				defer response.Content.Close()
				if response.ContentLength > 5*1024*1024 { // 限制 5MB
					h.bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("⚠️ 圖片太大了，請傳送 5MB 以下的圖片。")).Do()
					continue
				}
				contentBytes, err := io.ReadAll(response.Content)
				if err != nil {
					log.Println("讀取圖片內容失敗:", err)
					continue
				}
				// 2. OCR 辨識 (使用我們剛寫的 VisionService)
				text, err := h.visionSvc.ExtractTextFromImage(contentBytes)
				if err != nil {
					h.bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("⚠️ 無法辨識圖片中的文字")).Do()
					continue
				}

				// 3. 翻譯並回傳
				replyText, err := h.transService.TranslateMessage(sourceID, text)
				if err != nil || replyText == "" {
					h.bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("無法翻譯圖片內容")).Do()
					continue
				}
				h.bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyText)).Do()
			}
		}
	}
}

// getSourceID 取得唯一的來源 ID
func getSourceID(source *linebot.EventSource) string {
	if source.GroupID != "" {
		return source.GroupID
	} else if source.RoomID != "" {
		return source.RoomID
	}
	return source.UserID
}
