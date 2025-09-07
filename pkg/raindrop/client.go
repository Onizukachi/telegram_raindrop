package raindrop

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

const baseUrl = "https://raindrop.io"

type Client struct {
	clientID     string
	clientSecret string
	redirectUrl  string
}

func NewClient(clientID, clientSecret, redirectUrl string) *Client {
	return &Client{clientID: clientID, clientSecret: clientSecret, redirectUrl: redirectUrl}
}

func (client *Client) BuildOAuthLink(chatID int64) string {
	redirectUri := client.redirectUrl + fmt.Sprintf("?chat_id=%d", chatID)
	queryParams := fmt.Sprintf("?redirect_uri=%s&client_id=%s", redirectUri, client.clientID)
	return baseUrl + "/oauth/authorize" + queryParams
}

func (client *Client) ExchangeToken(code string) (*AuthResponse, error) {
	grantType := "authorization_code"

	data := exchangeRequest{
		baseAuthRequest: baseAuthRequest{
			GrantType:    grantType,
			ClientID:     client.clientID,
			ClientSecret: client.clientSecret,
		},
		Code:        code,
		RedirectUrl: client.redirectUrl,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	log.Printf("ReQUEST EXHANGE DATA %v", string(jsonData))
	url := baseUrl + "/oauth/access_token"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response AuthResponse

	log.Println("RESPONSE BODYY")
	log.Println(string(resBody))
	if err := json.Unmarshal(resBody, &response); err != nil {
		return nil, err
	}

	// проверить не валидный отве {"result":false,"status":400,"errorMessage":"Incorrect redirect_uri"}
	return &response, nil
}

func (client *Client) RefreshToken(token string) (*AuthResponse, error) {
	grantType := "refresh_token"

	data := refreshRequest{
		baseAuthRequest: baseAuthRequest{
			GrantType:    grantType,
			ClientID:     client.clientID,
			ClientSecret: client.clientSecret,
		},
		RefreshToken: token,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	url := baseUrl + "/oauth/access_token"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response AuthResponse

	if err := json.Unmarshal(resBody, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (client *Client) CreateItem(link, accessToken string) error {
	data := createItemRequest{Link: link}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	url := baseUrl + "/rest/v1/raindrop"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorizarion", "Bearer "+accessToken)

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response CreateItemResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return err
	}

	if !response.Result {
		return errors.New("response is not successfull")
	}

	return nil
}
