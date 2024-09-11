package model

type WebResponse struct {
	Error       bool   `json:"error"`
	StatusCode  int    `json:"status_code,omitempty"`
	Message     string `json:"message,omitempty"`
	RedirectUrl string `json:"redirect_url,omitempty"`
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
	TransactionId    string `json:"transactionId"`
	Status           string `json:"status"`
	ErrorTitle       string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (r *TelcoResponse) GetTransactionId() string {
	return r.TransactionId
}

func (r *TelcoResponse) GetStatus() string {
	return r.Status
}

func (r *TelcoResponse) IsSuccess() bool {
	return r.GetStatus() == "SUCCESS"
}

func (r *TelcoResponse) IsInvalid() bool {
	return r.GetStatus() == "INVALID_OR_EXPIRED_PIN"
}

func (r *TelcoResponse) GetErrorTitle() string {
	return r.ErrorTitle
}

func (r *TelcoResponse) GetErrorDescription() string {
	return r.ErrorDescription
}
