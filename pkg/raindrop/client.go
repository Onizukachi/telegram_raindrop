package raindrop

import "fmt"

const baseUrl = "https://raindrop.io"

type Client struct {
	clientID     string
	clientSecret string
	redirectUrl  string
}

func NewClient(clientID, clientSecret, redirectUrl string) *Client {
	return &Client{clientID: clientID, clientSecret: clientSecret, redirectUrl: redirectUrl}
}

func (client *Client) BuildOAuthLink() string {
	queryParams := fmt.Sprintf("?redirect_uri=%s&client_id=%s", client.redirectUrl, client.clientID)
	return baseUrl + "/oauth/authorize" + queryParams
}
