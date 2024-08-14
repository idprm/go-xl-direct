package model

import (
	"strings"
)

type CampaignDirectQueryRequest struct {
	Service        string `json:"service" query:"service"`
	Adnet          string `json:"adnet" query:"ad"`
	PubId          string `json:"pub_id" query:"pubid"`
	AffSub         string `json:"aff_sub" query:"aff_sub"`
	CampKeyword    string `json:"keyword" query:"keyword"`
	CampSubKeyword string `json:"subkey" query:"subkey"`
	IpAddress      string `json:"ip_address" query:"ip"`
}

func (r *CampaignDirectQueryRequest) GetService() string {
	return strings.ToUpper(r.Service)
}

func (r *CampaignDirectQueryRequest) GetAdnet() string {
	return r.Adnet
}

func (r *CampaignDirectQueryRequest) GetPubId() string {
	return r.PubId
}

func (r *CampaignDirectQueryRequest) GetAffSub() string {
	return r.AffSub
}

func (r *CampaignDirectQueryRequest) GetCampKeyword() string {
	return strings.ToUpper(r.CampKeyword)
}

func (r *CampaignDirectQueryRequest) GetCampSubKeyword() string {
	return strings.ToUpper(r.CampSubKeyword)
}

func (r *CampaignDirectQueryRequest) GetIpAddress() string {
	return r.IpAddress
}

func (r *CampaignDirectQueryRequest) SetService(data string) {
	r.Service = data
}

func (r *CampaignDirectQueryRequest) SetIpAddress(ip string) {
	r.IpAddress = ip
}

type CampaignToolsQueryRequest struct {
	Service   string `json:"srv" query:"srv"`
	Dynamic   string `json:"dyn" query:"dyn"`
	Adnet     string `json:"adnet" query:"ad"`
	PubId     string `json:"pub_id" query:"pubid"`
	AffSub    string `json:"aff_sub" query:"aff_sub"`
	IpAddress string `json:"ip_address" query:"ip"`
	OS        string `json:"os" query:"os"`
	Browser   string `json:"browser" query:"browser"`
	UA        string `json:"useragent" query:"useragent"`
	Referer   string `json:"referer" query:"referer"`
}

func (r *CampaignToolsQueryRequest) GetService() string {
	message := strings.ToUpper(r.Service)
	index := strings.Split(message, " ")
	if len(index[0]) > 0 {
		return index[0]
	}
	return ""
}

func (r *CampaignToolsQueryRequest) GetDynamic() string {
	message := strings.ToUpper(r.Dynamic)
	index := strings.Split(message, " ")
	if len(index[0]) > 0 {
		return index[0]
	}
	return ""
}

func (r *CampaignToolsQueryRequest) GetSubKeyword() string {
	message := strings.ToUpper(r.Service)
	index := strings.Split(message, " ")
	if len(index) > 1 {
		return index[1]
	}
	return ""
}

func (r *CampaignToolsQueryRequest) GetSubDynamic() string {
	message := strings.ToUpper(r.Dynamic)
	index := strings.Split(message, " ")
	if len(index) > 1 {
		return index[1]
	}
	return ""
}

func (r *CampaignToolsQueryRequest) GetAdnet() string {
	return r.Adnet
}

func (r *CampaignToolsQueryRequest) GetPubId() string {
	return r.PubId
}

func (r *CampaignToolsQueryRequest) GetAffSub() string {
	return r.AffSub
}

func (r *CampaignToolsQueryRequest) GetIpAddress() string {
	return r.IpAddress
}

func (r *CampaignToolsQueryRequest) SetIpAddress(ip string) {
	r.IpAddress = ip
}

func (r *CampaignToolsQueryRequest) IsBillable() bool {
	return r.GetSubKeyword() == "LNK" ||
		strings.Contains(r.GetSubKeyword(), "BLB") ||
		strings.Contains(r.GetSubKeyword(), "BIL")
}

func (r *CampaignToolsQueryRequest) IsSam() bool {
	return r.GetSubKeyword() == "SAM"
}

type CampaignToolsResponse struct {
	StatusCode int    `json:"status_code" xml:"status_code"`
	Token      string `json:"token" xml:"token"`
	UrlPromo   string `json:"url_promo" xml:"url_promo"`
}

type SuccessQueryParamsRequest struct {
	Token     string `query:"token"`
	TrxId     string `query:"trx_id"`
	IpAddress string `query:"ip" json:"ip_address"`
}

func (e *SuccessQueryParamsRequest) GetToken() string {
	return e.Token
}

func (e *SuccessQueryParamsRequest) GetTrxId() string {
	return e.TrxId
}

func (e *SuccessQueryParamsRequest) GetIpAddress() string {
	return e.IpAddress
}

func (e *SuccessQueryParamsRequest) SetIpAddress(data string) {
	e.IpAddress = data
}

type OAuthRequest struct {
	ClientId     string `form:"client_id" json:"client_id"`
	ClientSecret string `form:"client_secret" json:"client_secret"`
	GrantType    string `form:"grant_type" json:"grant_type"`
}

type CreateSubscriptionRequest struct {
	RequestId      string `json:"requestId"`
	ProductId      string `json:"productId"`
	UserIdentifier string `json:"userIdentifier"`
	Amount         string `json:"amount"`
}

type ConfirmOTPRequest struct {
	RequestId string `json:"requestId"`
	PIN       string `json:"pin"`
}

func (r *ConfirmOTPRequest) GetRequestId() string {
	return r.RequestId
}

func (r *ConfirmOTPRequest) GetPIN() string {
	return r.PIN
}

func (r *ConfirmOTPRequest) SetRequestId(val string) {
	r.RequestId = val
}

func (r *ConfirmOTPRequest) SetPIN(val string) {
	r.PIN = val
}

type RefundRequest struct {
	RequestId     string `json:"requestId"`
	TransactionId string `json:"transactionId"`
}

func (r *RefundRequest) SetRequestId(val string) {
	r.RequestId = val
}

func (r *RefundRequest) SetTransactionId(val string) {
	r.TransactionId = val
}

type UnsubscribeRequest struct {
	RequestId string `json:"requestId"`
}

func (r *UnsubscribeRequest) SetRequestId(val string) {
	r.RequestId = val
}

type NotificationRequest struct {
	UserIdentifier  string  `validate:"required" json:"userIdentifier"`
	SubscriptionId  int     `validate:"required" json:"subscriptionId"`
	ProductId       string  `validate:"required" json:"productId"`
	Amount          float64 `validate:"required" json:"amount"`
	StartDate       string  `json:"startDate"`
	NextRenewalDate string  `json:"nextRenewalDate"`
	Status          string  `validate:"required" json:"status"`
	Context         string  `validate:"required" json:"context"`
	TransactionId   string  `json:"transactionId"`
}

func (m *NotificationRequest) GetUserIdentifier() string {
	return m.UserIdentifier
}

func (m *NotificationRequest) GetSubscriptionId() int {
	return m.SubscriptionId
}

func (m *NotificationRequest) GetProductId() string {
	return m.ProductId
}

func (m *NotificationRequest) GetAmount() float64 {
	return m.Amount
}

func (m *NotificationRequest) GetStartDate() string {
	return m.StartDate
}

func (m *NotificationRequest) GetNextRenewalDate() string {
	return m.NextRenewalDate
}

func (m *NotificationRequest) GetStatus() string {
	return m.Status
}

func (m *NotificationRequest) GetContext() string {
	return m.Context
}

func (m *NotificationRequest) GetTransactionId() string {
	return m.TransactionId
}

func (m *NotificationRequest) IsSubscription() bool {
	return m.Context == "SUBSCRIPTION"
}

func (m *NotificationRequest) IsRenewal() bool {
	return m.Context == "RENEWAL"
}

func (m *NotificationRequest) IsRefund() bool {
	return m.Context == "REFUND"
}

func (m *NotificationRequest) IsActive() bool {
	return m.Status == "ACTIVE"
}

func (m *NotificationRequest) IsCancelled() bool {
	return m.Status == "CANCELLED"
}

type NotificationSubcriptionRequest struct {
	UserIdentifier  string  `json:"userIdentifier"`
	SubscriptionId  int     `json:"subscriptionId"`
	ProductId       string  `json:"productId"`
	Amount          float64 `json:"amount"`
	StartDate       string  `json:"startDate"`
	NextRenewalDate string  `json:"nextRenewalDate"`
	Status          string  `json:"status"`
	Context         string  `json:"context"`
	TransactionId   string  `json:"transactionId"`
}

type NotificationUnSubcriptionRequest struct {
	UserIdentifier string `json:"userIdentifier"`
	SubscriptionId int    `json:"subscriptionId"`
	ProductId      string `json:"productId"`
	Status         string `json:"status"`
	Context        string `json:"context"`
}

type NotificationRenewalRequest struct {
	UserIdentifier  string  `json:"userIdentifier"`
	SubscriptionId  int     `json:"subscriptionId"`
	ProductId       string  `json:"productId"`
	Amount          float64 `json:"amount"`
	NextRenewalDate string  `json:"nextRenewalDate"`
	Status          string  `json:"status"`
	Context         string  `json:"context"`
	TransactionId   string  `json:"transactionId"`
}

type NotificationRefundRequest struct {
	UserIdentifier string `json:"userIdentifier"`
	SubscriptionId int    `json:"subscriptionId"`
	ProductId      string `json:"productId"`
	Status         string `json:"status"`
	Context        string `json:"context"`
}

type WebSubRequest struct {
	Msisdn    string `json:"msisdn"`
	Service   string `json:"service"`
	IpAddress string `json:"ip_address"`
}

func (r *WebSubRequest) SetIpAddress(val string) {
	r.IpAddress = val
}

func (r *WebSubRequest) GetMsisdn() string {
	return r.Msisdn
}

func (r *WebSubRequest) GetService() string {
	return r.Service
}

type WebOTPRequest struct {
	Msisdn  string `json:"msisdn"`
	Service string `json:"service"`
	Pin     string `json:"pin"`
}

func (r *WebOTPRequest) GetMsisdn() string {
	return r.Msisdn
}

func (r *WebOTPRequest) GetService() string {
	return r.Service
}

func (r *WebOTPRequest) GetPin() string {
	return r.Pin
}
