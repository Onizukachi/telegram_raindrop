package raindrop

type baseAuthRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type exchangeRequest struct {
	baseAuthRequest
	Code        string `json:"code"`
	RedirectUrl string `json:"redirect_uri"`
}

type refreshRequest struct {
	baseAuthRequest
	RefreshToken string `json:"refresh_token"`
}

type createItemRequest struct {
	Link string `json:"link"`
}
