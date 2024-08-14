package entity

type Verify struct {
	TrxId     string   `json:"trx_id,omitempty"`
	Msisdn    string   `json:"msisdn,omitempty"`
	Service   *Service `json:"service,omitempty"`
	IpAddress string   `json:"ip_address" query:"ip_address"`
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
