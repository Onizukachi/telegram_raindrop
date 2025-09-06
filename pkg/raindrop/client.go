package raindrop

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseUrl = "https://raindrop.io"
const grantType = "authorization_code"

type Client struct {
	clientID     string
	clientSecret string
	redirectUrl  string
}

type exchangeRequest struct {
	GrantType    string `json:"grant_type"`
	Code         string `json:"code"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectUrl  string `json:"redirect_uri"`
}

type ExchangeResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

func NewClient(clientID, clientSecret, redirectUrl string) *Client {
	return &Client{clientID: clientID, clientSecret: clientSecret, redirectUrl: redirectUrl}
}

func (client *Client) BuildOAuthLink() string {
	queryParams := fmt.Sprintf("?redirect_uri=%s&client_id=%s", client.redirectUrl, client.clientID)
	return baseUrl + "/oauth/authorize" + queryParams
}

func (client *Client) ExchangeToken(code string) (*ExchangeResponse, error) {
	data := exchangeRequest{
		GrantType:    grantType,
		Code:         code,
		ClientID:     client.clientID,
		ClientSecret: client.clientSecret,
		RedirectUrl:  client.redirectUrl,
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

	var response ExchangeResponse

	if err := json.Unmarshal(resBody, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
