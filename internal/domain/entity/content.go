package entity

import (
	"strings"
)

type Content struct {
	ID        int `json:"id"`
	ServiceID int `json:"service_id"`
	Service   *Service
	Name      string `json:"name"`
	Value     string `json:"value"`
	Tid       string `json:"tid"`
	Sequence  int    `json:"sequence"`
}

func (c *Content) GetName() string {
	return c.Name
}

func (c *Content) GetValue() string {
	return c.Value
}

func (c *Content) GetTid() string {
	return c.Tid
}

func (c *Content) GetSequence() int {
	return c.Sequence
}

func (c *Content) SetPIN(pin string) {
	replacer := strings.NewReplacer("@pin", pin)
	c.Value = replacer.Replace(c.Value)
}
