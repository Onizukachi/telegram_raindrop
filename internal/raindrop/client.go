package raindrop

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const baseUrl = "https://raindrop.io"
const apiBaseUrl = "https://api.raindrop.io/rest/"

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
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if response.AccessToken == "" {
		return nil, fmt.Errorf("auth error: %s (status %d)", response.ErrorMessage, response.Status)
	}

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

	url := apiBaseUrl + "v1/raindrop"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response: %s", resp.Status)
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
