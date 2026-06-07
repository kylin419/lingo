package main

import (
	"log"
	"net/http"
	"os"

	"line-translate-bot/internal/handler"
	"line-translate-bot/internal/repository"
	"line-translate-bot/internal/service"
	"line-translate-bot/internal/translator"

	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v8/linebot"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ 未偵測到 .env 檔案，將直接使用系統環境變數")
	} else {
		log.Println(".env 檔案載入成功")
	}
	secret := os.Getenv("LINE_CHANNEL_SECRET")
	token := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")
	azureKey := os.Getenv("AZURE_API_KEY")
	azureRegion := os.Getenv("AZURE_API_REGION")

	if secret == "" || token == "" || azureKey == "" || azureRegion == "" {
		if secret == "" {
			log.Println("LINE_CHANNEL_SECRET 未設定")
		}
		if token == "" {
			log.Println("LINE_CHANNEL_ACCESS_TOKEN 未設定")
		}
		if azureKey == "" {
			log.Println("AZURE_API_KEY 未設定")
		}
		if azureRegion == "" {
			log.Println("AZURE_API_REGION 未設定")
		}
	}

	// 初始化 LINE Bot 實體
	bot, err := linebot.New(secret, token)
	if err != nil {
		log.Fatal("LINE Bot 初始化失敗: ", err)
	}
	repo := repository.NewMemoryRepository()

	trans := translator.NewAzureTranslator(azureKey, azureRegion)

	cmdService := service.NewCommandService(repo)
	transService := service.NewTranslationService(repo, trans)

	// 4. 初始化接收層 (Handler)
	lineHandler := handler.NewLineHandler(bot, cmdService, transService)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Bot is alive"))
	})

	// 註冊 HTTP 路由
	http.HandleFunc("/callback", lineHandler.HandleRequest)

	// 設定監聽 Port
	port := "8081"

	log.Printf("LINE Bot 伺服器已啟動，強制監聽 Port: %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("監聽失敗: ", err)
	}
}
