package model

type WebResponse struct {
	Error      bool   `json:"error"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

type OAuthResponse struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func (r *OAuthResponse) GetTokenType() string {
	return r.TokenType
}

func (r *OAuthResponse) GetAccessToken() string {
	return r.AccessToken
}

func (r *OAuthResponse) GetExpiresIn() int {
	return r.ExpiresIn
}

type TelcoResponse struct {
	TransactionId string `json:"transactionId"`
	Status        string `json:"status"`
}

func (r *TelcoResponse) GetTransactionId() string {
	return r.TransactionId
}

func (r *TelcoResponse) GetStatus() string {
	return r.Status
}
