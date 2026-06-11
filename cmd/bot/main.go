package main

import (
	"context"
	"line-translate-bot/internal/vision"
	"log"
	"net/http"
	"os"

	"line-translate-bot/internal/handler"
	"line-translate-bot/internal/repository"
	"line-translate-bot/internal/service"
	"line-translate-bot/internal/translator"

	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/redis/go-redis/v9"
)

var (
	ctx = context.Background()
	rdb *redis.Client
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
	azureVisionEndpoint := os.Getenv("AZURE_VISION_ENDPOINT")
	azureVisionKey := os.Getenv("AZURE_VISION_KEY")
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
	// 初始化 Redis 資料庫
	rdb := repository.NewRedisClient()
	repo := repository.NewRedisRepository(rdb)
	// 初始化 Azure 翻譯引擎
	trans := translator.NewAzureTranslator(azureKey, azureRegion)
	// 初始化 LINE Bot
	bot, err := messaging_api.NewMessagingApiAPI(token)
	blobbot, err := messaging_api.NewMessagingApiBlobAPI(token)
	if err != nil {
		log.Fatal("LINE Bot 初始化失敗: ", err)
	}
	// 依賴注入
	cmdService := service.NewCommandService(repo)
	transService := service.NewTranslationService(repo, trans)
	// Vision Service
	visionSvc := vision.NewAzureVision(azureVisionEndpoint, azureVisionKey)
	// 初始化接收層 (Handler)
	lineHandler := handler.NewLineHandler(bot, blobbot, cmdService, transService, visionSvc)

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
