package model

import "strings"

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

type WebResponse struct {
	Error      bool   `json:"error"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}
