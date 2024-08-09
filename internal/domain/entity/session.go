package entity

type Session struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func (e *Session) GetTokenType() string {
	return e.TokenType
}

func (e *Session) GetAccessToken() string {
	return e.AccessToken
}

func (e *Session) GetExpiresIn() int {
	return e.ExpiresIn
}
