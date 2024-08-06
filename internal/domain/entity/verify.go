package entity

import "strings"

type Verify struct {
	TxId           string `json:"tx_id,omitempty"`
	Token          string `json:"token,omitempty"`
	Service        string `json:"service,omitempty"`
	Adnet          string `json:"adnet,omitempty"`
	PubID          string `json:"pub_id,omitempty"`
	AffSub         string `json:"aff_sub,omitempty"`
	CampKeyword    string `json:"camp_keyword,omitempty"`
	CampSubKeyword string `json:"camp_sub_keyword,omitempty"`
	IpAddress      string `json:"ip_address,omitempty"`
	IsBillable     bool   `json:"is_billable,omitempty"`
	IsCampTool     bool   `json:"is_camptool,omitempty"`
}

func (v *Verify) GetTxId() string {
	return v.TxId
}

func (v *Verify) GetToken() string {
	return v.Token
}

func (v *Verify) GetService() string {
	return v.Service
}

func (v *Verify) GetAdnet() string {
	return v.Adnet
}

func (v *Verify) GetPubId() string {
	return v.PubID
}

func (v *Verify) GetAffSub() string {
	return v.AffSub
}

func (v *Verify) GetCampKeyword() string {
	return v.CampKeyword
}

func (v *Verify) GetCampSubKeyword() string {
	return v.CampSubKeyword
}

func (v *Verify) GetIpAddress() string {
	return v.IpAddress
}

func (v *Verify) GetIsBillable() bool {
	return v.IsBillable
}

func (v *Verify) GetIsCampTool() bool {
	return v.IsCampTool
}

func (v *Verify) SetCampKeyword(keyword string) {
	v.CampKeyword = strings.ToUpper(keyword)
}

func (v *Verify) SetCampSubKeyword(subkey string) {
	v.CampSubKeyword = strings.ToUpper(subkey)
}

func (v *Verify) IsCampKeyword() bool {
	return v.CampKeyword != ""
}
