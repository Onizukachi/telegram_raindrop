package raindrop

type AuthResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	Result       bool   `json:"result,omitempty"`
	Status       int    `json:"status,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

type CreateItemResponse struct {
	Result bool `json:"result"`
}
