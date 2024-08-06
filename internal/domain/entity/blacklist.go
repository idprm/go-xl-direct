package entity

import "time"

type Blacklist struct {
	ID        int       `json:"id"`
	Msisdn    string    `json:"msisdn"`
	CreatedAt time.Time `json:"created_at"`
}

func (e *Blacklist) GetMsisdn() string {
	return e.Msisdn
}
