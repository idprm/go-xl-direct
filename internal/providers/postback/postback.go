package postback

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/idprm/go-xl-direct/internal/domain/entity"
	"github.com/idprm/go-xl-direct/internal/logger"
	"github.com/idprm/go-xl-direct/internal/utils"
	"github.com/sirupsen/logrus"
)

type Postback struct {
	logger       *logger.Logger
	subscription *entity.Subscription
	service      *entity.Service
}

func NewPostback(
	logger *logger.Logger,
	subscription *entity.Subscription,
	service *entity.Service,
) *Postback {
	return &Postback{
		logger:       logger,
		subscription: subscription,
		service:      service,
	}
}

func (p *Postback) Send() ([]byte, error) {
	l := p.logger.Init("pb", true)

	start := time.Now()
	trxId := utils.GenerateTrxId()

	q := url.Values{}
	q.Add("partner", "linkitxl")
	q.Add("px", p.subscription.GetAdnet())
	q.Add("serv_id", p.service.GetCode())
	q.Add("msisdn", p.subscription.GetMsisdn())
	q.Add("trxid", p.subscription.GetLatestTrxId())

	req, err := http.NewRequest("GET", p.service.UrlPostback+"?"+q.Encode(), nil)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}

	p.logger.Writer(req)
	l.WithFields(logrus.Fields{
		"msisdn":  p.subscription.Msisdn,
		"request": p.service.UrlPostback + "?" + q.Encode(),
		"trx_id":  trxId,
	}).Info("POSTBACK")

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	duration := time.Since(start).Milliseconds()
	p.logger.Writer(string(body))
	l.WithFields(logrus.Fields{
		"msisdn":      p.subscription.GetMsisdn(),
		"response":    string(body),
		"trx_id":      trxId,
		"duration":    duration,
		"status_code": resp.StatusCode,
		"status_text": http.StatusText(resp.StatusCode),
	}).Info("POSTBACK")

	return body, nil
}

func (p *Postback) Billable() ([]byte, error) {
	l := p.logger.Init("pb", true)

	start := time.Now()
	trxId := utils.GenerateTrxId()

	q := url.Values{}
	q.Add("partner", "linkitxl")
	q.Add("px", p.subscription.GetAdnet())
	q.Add("serv_id", p.service.GetCode())
	q.Add("msisdn", p.subscription.GetMsisdn())
	q.Add("trxid", p.subscription.GetLatestTrxId())

	req, err := http.NewRequest("GET", p.service.UrlPostbackBillable+"?"+q.Encode(), nil)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}

	p.logger.Writer(req)
	l.WithFields(logrus.Fields{
		"msisdn":  p.subscription.GetMsisdn(),
		"request": p.service.UrlPostbackBillable + "?" + q.Encode(),
		"trx_id":  trxId,
	}).Info("BILLABLE")

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	duration := time.Since(start).Milliseconds()
	p.logger.Writer(string(body))
	l.WithFields(logrus.Fields{
		"msisdn":      p.subscription.GetMsisdn(),
		"response":    string(body),
		"trx_id":      trxId,
		"duration":    duration,
		"status_code": resp.StatusCode,
		"status_text": http.StatusText(resp.StatusCode),
	}).Info("BILLABLE")

	return body, nil
}
