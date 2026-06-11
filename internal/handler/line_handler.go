package handler

import (
	"io"
	"line-translate-bot/internal/service"
	"line-translate-bot/internal/vision"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

type LineHandler struct {
	bot          *messaging_api.MessagingApiAPI
	blobbot      *messaging_api.MessagingApiBlobAPI
	cmdService   *service.CommandService
	transService *service.TranslationService
	visionSvc    vision.VisionService
}

func NewLineHandler(bot *messaging_api.MessagingApiAPI, blobbot *messaging_api.MessagingApiBlobAPI, cmdService *service.CommandService, transService *service.TranslationService, visionSvc vision.VisionService) *LineHandler {
	return &LineHandler{
		bot:          bot,
		blobbot:      blobbot,
		cmdService:   cmdService,
		transService: transService,
		visionSvc:    visionSvc,
	}
}

func (h *LineHandler) HandleRequest(w http.ResponseWriter, req *http.Request) {
	events, err := webhook.ParseRequest(os.Getenv("LINE_CHANNEL_SECRET"), req)
	if err != nil {
		if err == webhook.ErrInvalidSignature {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	for _, event := range events.Events {

		switch e := event.(type) {
		// 處理加入群組事件
		case webhook.JoinEvent:
			welcomeMsg := "大家好！我是翻譯機器人。\n請輸入 `$set [語言1] [語言2]...` 來設定要互翻的語言。\n例如：`$set zh-Hant en ja`"
			if _, err := h.bot.ReplyMessage(&messaging_api.ReplyMessageRequest{
				ReplyToken: e.ReplyToken,
				Messages: []messaging_api.MessageInterface{
					messaging_api.TextMessage{
						Text: welcomeMsg,
					},
				},
			}); err != nil {
				log.Println("回覆歡迎訊息失敗:", err)
			}

		// 處理文字訊息事件
		case webhook.MessageEvent:
			sourceID := getSourceID(e.Source)
			switch msg := e.Message.(type) {
			case webhook.TextMessageContent:
				if message, ok := e.Message.(webhook.TextMessageContent); ok {
					text := strings.TrimSpace(message.Text)

					// 判斷是否為 $set 指令
					if strings.HasPrefix(text, "$set ") {
						replyText := h.cmdService.HandleSetCommand(sourceID, text)
						h.bot.ReplyMessage(&messaging_api.ReplyMessageRequest{
							ReplyToken: e.ReplyToken,
							Messages: []messaging_api.MessageInterface{
								messaging_api.TextMessage{
									Text: replyText,
								},
							},
						})
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
							h.bot.ReplyMessage(&messaging_api.ReplyMessageRequest{
								ReplyToken: e.ReplyToken,
								Messages: []messaging_api.MessageInterface{
									messaging_api.TextMessage{
										Text: replyText,
									},
								},
							})
						}
					}
				}
			case webhook.ImageMessageContent:
				response, err := h.blobbot.GetMessageContent(msg.Id)
				if err != nil {
					log.Println("圖片下載失敗:", err)
					continue
				}
				defer response.Body.Close()
				if response.ContentLength > 5*1024*1024 { // 限制 5MB
					h.bot.ReplyMessage(&messaging_api.ReplyMessageRequest{
						ReplyToken: e.ReplyToken,
						Messages: []messaging_api.MessageInterface{
							messaging_api.TextMessage{
								Text: "⚠️ 圖片太大了，請傳送 5MB 以下的圖片。",
							},
						},
					})
					continue
				}
				contentBytes, err := io.ReadAll(response.Body)
				if err != nil {
					log.Println("讀取圖片內容失敗:", err)
					continue
				}
				// 2. OCR 辨識 (使用我們剛寫的 VisionService)
				text, err := h.visionSvc.ExtractTextFromImage(contentBytes)
				if err != nil {
					h.bot.ReplyMessage(&messaging_api.ReplyMessageRequest{
						ReplyToken: e.ReplyToken,
						Messages: []messaging_api.MessageInterface{
							messaging_api.TextMessage{
								Text: "⚠️ 無法辨識圖片中的文字",
							},
						},
					})
					continue
				}

				// 3. 翻譯並回傳
				replyText, err := h.transService.TranslateMessage(sourceID, text)
				if err != nil || replyText == "" {
					h.bot.ReplyMessage(&messaging_api.ReplyMessageRequest{
						ReplyToken: e.ReplyToken,
						Messages: []messaging_api.MessageInterface{
							messaging_api.TextMessage{
								Text: "無法翻譯圖片內容",
							},
						},
					})
					continue
				}
				h.bot.ReplyMessage(&messaging_api.ReplyMessageRequest{
					ReplyToken: e.ReplyToken,
					Messages: []messaging_api.MessageInterface{
						messaging_api.TextMessage{
							Text: replyText,
						},
					},
				})
			}
		}
	}
}

// getSourceID 取得唯一的來源 ID
func getSourceID(source webhook.SourceInterface) string {
	if source == nil {
		return ""
	}
	switch v := source.(type) {
	case *webhook.GroupSource:
		return v.GroupId
	case webhook.GroupSource:
		return v.GroupId
	case *webhook.RoomSource:
		return v.RoomId
	case webhook.RoomSource:
		return v.RoomId
	case *webhook.UserSource:
		return v.UserId
	case webhook.UserSource:
		return v.UserId
	default:
		return ""
	}
}
