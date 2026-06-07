package handler

import (
	"log"
	"net/http"
	"strings"

	"line-translate-bot/internal/service"

	"github.com/line/line-bot-sdk-go/v8/linebot"
)

type LineHandler struct {
	bot          *linebot.Client
	cmdService   *service.CommandService
	transService *service.TranslationService
}

func NewLineHandler(bot *linebot.Client, cmdService *service.CommandService, transService *service.TranslationService) *LineHandler {
	return &LineHandler{
		bot:          bot,
		cmdService:   cmdService,
		transService: transService,
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
