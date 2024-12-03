package entity

type Verify struct {
	TrxId      string   `json:"trx_id,omitempty"`
	Service    *Service `json:"service,omitempty"`
	Msisdn     string   `json:"msisdn,omitempty"`
	Channel    string   `json:"channel,omitempty"`
	PIN        string   `json:"pin,omitempty"`
	Keyword    string   `json:"keyword,omitempty"`
	SubKeyword string   `json:"sub_keyword,omitempty"`
	Adnet      string   `json:"adnet,omitempty"`
	PubID      string   `json:"pub_id,omitempty"`
	AffSub     string   `json:"aff_sub,omitempty"`
	IpAddress  string   `json:"ip_address"`
}

func (v *Verify) GetTrxId() string {
	return v.TrxId
}

func (v *Verify) GetMsisdn() string {
	return v.Msisdn
}

func (r *Verify) GetIpAddress() string {
	return r.IpAddress
}

func (r *Verify) SetTrxId(val string) {
	r.TrxId = val
}

func (r *Verify) SetIpAddress(val string) {
	r.IpAddress = val
}
