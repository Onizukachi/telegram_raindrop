package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Onizukachi/telegram_raindrop/internal/config"
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
	messages       config.Messages
}

func NewServer(addr, botName string, botApi *tgbotapi.BotAPI, raindropClient *raindrop.Client, userRepo storage.UserRepository, messages config.Messages) *Server {
	return &Server{
		addr:           addr,
		botName:        botName,
		botApi:         botApi,
		raindropClient: raindropClient,
		userRepo:       userRepo,
		messages:       messages,
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
	botUrl := fmt.Sprintf("https://t.me/%s", s.botName)
	code := r.URL.Query().Get("code")
	chatIDStr := r.URL.Query().Get("chat_id")
	chatID, _ := strconv.ParseInt(chatIDStr, 10, 64)

	exchangeResponse, err := s.raindropClient.ExchangeToken(code)
	if err != nil {
		s.sendMessage(chatID, s.messages.Errors.FailAuth)
		http.Redirect(w, r, botUrl, http.StatusMovedPermanently)
		return
	}

	expriresIn := time.Now().Add(time.Second * time.Duration(exchangeResponse.ExpiresIn))
	err = s.userRepo.Create(chatID, exchangeResponse.AccessToken, exchangeResponse.RefreshToken, expriresIn)
	if err != nil {
		s.sendMessage(chatID, s.messages.Errors.FailAuth)
		http.Redirect(w, r, botUrl, http.StatusMovedPermanently)
		return
	}

	s.sendMessage(chatID, s.messages.Responses.SuccessAuthorized)
	http.Redirect(w, r, botUrl, http.StatusMovedPermanently)
}

func (s *Server) sendMessage(chatID int64, message string) error {
	msg := tgbotapi.NewMessage(chatID, message)
	_, err := s.botApi.Send(msg)
	return err
}
