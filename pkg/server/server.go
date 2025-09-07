package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Onizukachi/telegram_raindrop/pkg/raindrop"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Server struct {
	addr           string
	botName        string
	botApi         *tgbotapi.BotAPI
	raindropClient *raindrop.Client
}

// Доработать тут логику
func NewServer(addr, botName string, botApi *tgbotapi.BotAPI, raindropClient *raindrop.Client) *Server {
	return &Server{
		addr:           addr,
		botName:        botName,
		botApi:         botApi,
		raindropClient: raindropClient,
	}
}

func (s *Server) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleCallback)

	fmt.Println("HTTP server is running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		return err
	}

	return nil
}

func (s *Server) handleCallback(w http.ResponseWriter, r *http.Request) {
	botUrl := fmt.Sprintf("https://t.me/%s?start=code", s.botName)
	code := r.URL.Query().Get("code")
	chatIDStr := r.URL.Query().Get("chat_id")
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		log.Println(err)
		// textMsg := "Не удалось прочитать chatId"
		// Обработать когда chatId не валидный или не найден у нас такой чат репу прокинуть наверно надо
		http.Redirect(w, r, botUrl, http.StatusMovedPermanently)
	}

	echangeResponse, err := s.raindropClient.ExchangeToken(code)
	if err != nil {
		log.Println(err)
		log.Println("ERRRRORRRR")
		http.Redirect(w, r, botUrl, http.StatusMovedPermanently)
		return
	}

	log.Printf("%+v\n", echangeResponse)

	textMsg := "Поздравляю! Ты успешно авторизовался :)"
	msg := tgbotapi.NewMessage(chatID, textMsg)
	s.botApi.Send(msg)
	http.Redirect(w, r, botUrl, http.StatusMovedPermanently)
}
