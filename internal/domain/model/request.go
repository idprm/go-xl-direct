package model

import (
	"strings"
	"time"

	"github.com/idprm/go-xl-direct/internal/domain/entity"
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
	RequestId       string `json:"requestId"`
	ProductId       string `json:"productId"`
	UserIdentifier  string `json:"userIdentifier"`
	Amount          string `json:"amount"`
	TransactionInfo struct {
		PartnerId string `json:"partnerId"`
	} `json:"transactionInfo"`
}

func (m *CreateSubscriptionRequest) SetPartnerId(v string) {
	m.TransactionInfo.PartnerId = v
}

type ConfirmOTPRequest struct {
	RequestId       string `json:"requestId"`
	PIN             string `json:"pin"`
	TransactionInfo struct {
		PartnerId string `json:"partnerId"`
	} `json:"transactionInfo"`
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

func (m *ConfirmOTPRequest) SetPartnerId(v string) {
	m.TransactionInfo.PartnerId = v
}

type RefundRequest struct {
	RequestId       string `json:"requestId"`
	TransactionId   string `json:"transactionId"`
	TransactionInfo struct {
		PartnerId string `json:"partnerId"`
	} `json:"transactionInfo"`
}

func (r *RefundRequest) SetRequestId(val string) {
	r.RequestId = val
}

func (r *RefundRequest) SetTransactionId(val string) {
	r.TransactionId = val
}

func (m *RefundRequest) SetPartnerId(v string) {
	m.TransactionInfo.PartnerId = v
}

type UnsubscribeRequest struct {
	RequestId       string `json:"requestId"`
	TransactionInfo struct {
		PartnerId string `json:"partnerId"`
	} `json:"transactionInfo"`
}

func (r *UnsubscribeRequest) SetRequestId(val string) {
	r.RequestId = val
}

func (m *UnsubscribeRequest) SetPartnerId(v string) {
	m.TransactionInfo.PartnerId = v
}

type NotificationRequest struct {
	UserIdentifier  string  `validate:"required" json:"userIdentifier"`
	SubscriptionId  int64   `validate:"required" json:"subscriptionId"`
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

func (m *NotificationRequest) GetSubscriptionId() int64 {
	return m.SubscriptionId
}

func (m *NotificationRequest) GetProductId() string {
	return m.ProductId
}

func (m *NotificationRequest) GetAmount() float64 {
	return m.Amount
}

func (m *NotificationRequest) GetStartDate() time.Time {
	date, _ := time.Parse(time.RFC3339, m.StartDate)
	return date
}

func (m *NotificationRequest) GetNextRenewalDate() time.Time {
	date, _ := time.Parse(time.RFC3339, m.NextRenewalDate)
	return date
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
	Service    string `validate:"required" json:"service"`
	Msisdn     string `validate:"required" json:"msisdn"`
	Channel    string `json:"channel"`
	Keyword    string `json:"keyword"`
	SubKeyword string `json:"subkey"`
	Adnet      string `json:"adnet"`
	PubId      string `json:"pubid"`
	AffSub     string `json:"aff_sub"`
	IpAddress  string `json:"ip_address"`
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
	Msisdn     string `validate:"required" json:"msisdn"`
	Service    string `validate:"required" json:"service"`
	Pin        string `validate:"required" json:"pin"`
	Channel    string `json:"channel"`
	Keyword    string `json:"keyword"`
	SubKeyword string `json:"subkey"`
	Adnet      string `json:"adnet"`
	PubId      string `json:"pubid"`
	AffSub     string `json:"aff_sub"`
	IpAddress  string `json:"ip_address"`
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

func (r *WebOTPRequest) GetKeyword() string {
	replacer := strings.NewReplacer("+", " ")
	return strings.ToUpper(replacer.Replace(r.Keyword))
}

func (r *WebOTPRequest) GetSubKeyword() string {
	return strings.ToUpper(r.SubKeyword)
}

type WebUnRegRequest struct {
	Msisdn string `query:"msisdn" json:"msisdn"`
}

func (r *WebUnRegRequest) GetMsisdn() string {
	return r.Msisdn
}

type NotifParamsRequest struct {
	Subscription *entity.Subscription
	Service      *entity.Service
	Action       string `json:"action"`
	Pin          string `json:"pin"`
}

func (e *NotifParamsRequest) GetAction() string {
	return strings.ToUpper(e.Action)
}

func (e *NotifParamsRequest) IsSub() bool {
	return e.GetAction() == "SUB"
}

func (e *NotifParamsRequest) IsRenewal() bool {
	return e.GetAction() == "RENEWAL"
}

func (e *NotifParamsRequest) IsUnsub() bool {
	return e.GetAction() == "UNSUB"
}

func (e *NotifParamsRequest) GetPin() string {
	return strings.ToLower(e.Pin)
}
