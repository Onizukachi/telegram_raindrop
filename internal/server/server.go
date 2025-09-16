package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Onizukachi/telegram_raindrop/internal/raindrop"
	"github.com/Onizukachi/telegram_raindrop/internal/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Server struct {
	addr           string
	botName        string
	botApi         *tgbotapi.BotAPI
	raindropClient *raindrop.Client
	userRepo       storage.UserRepository
}

// Доработать тут логику
func NewServer(addr, botName string, botApi *tgbotapi.BotAPI, raindropClient *raindrop.Client, userRepo storage.UserRepository) *Server {
	return &Server{
		addr:           addr,
		botName:        botName,
		botApi:         botApi,
		raindropClient: raindropClient,
		userRepo:       userRepo,
	}
}

func (s *Server) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleCallback)

	if err := http.ListenAndServe(s.addr, mux); err != nil {
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
		// textMsg := "Не удалось прочитать chatId"
		// Обработать когда chatId не валидный или не найден у нас такой чат репу прокинуть наверно надо
		http.Redirect(w, r, botUrl, http.StatusMovedPermanently)
	}

	exchangeResponse, err := s.raindropClient.ExchangeToken(code)
	if err != nil {
		http.Redirect(w, r, botUrl, http.StatusMovedPermanently)
		return
	}

	expriresIn := time.Now().Add(time.Second * time.Duration(exchangeResponse.ExpiresIn))
	err = s.userRepo.Create(chatID, exchangeResponse.AccessToken, exchangeResponse.RefreshToken, expriresIn)
	if err != nil {
		http.Redirect(w, r, botUrl, http.StatusMovedPermanently)
		return
	}
	log.Printf("%+v\n", &exchangeResponse)

	textMsg := "Поздравляю! Ты успешно авторизовался :)"
	msg := tgbotapi.NewMessage(chatID, textMsg)
	s.botApi.Send(msg)
	http.Redirect(w, r, botUrl, http.StatusMovedPermanently)
}
