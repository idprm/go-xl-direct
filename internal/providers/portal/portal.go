package portal

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

type Portal struct {
	logger       *logger.Logger
	subscription *entity.Subscription
	service      *entity.Service
	pin          string
	status       string
}

func NewPortal(
	logger *logger.Logger,
	subscription *entity.Subscription,
	service *entity.Service,
	pin string,
	status string,
) *Portal {
	return &Portal{
		logger:       logger,
		subscription: subscription,
		service:      service,
		pin:          pin,
		status:       status,
	}
}

func (p *Portal) Subscription() ([]byte, error) {
	l := p.logger.Init("notif", true)

	start := time.Now()
	trxId := utils.GenerateTrxId()

	q := url.Values{}
	q.Add("telco", "telkomsel")
	q.Add("msisdn", p.subscription.Msisdn)
	q.Add("event", "reg")
	q.Add("password", p.pin)
	q.Add("package", "daily")
	q.Add("status", p.status)
	q.Add("time", time.Now().String())

	req, err := http.NewRequest("GET", p.service.UrlNotifSub+"?"+q.Encode(), nil)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    5 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: tr,
	}
	p.logger.Writer(req)
	l.WithFields(logrus.Fields{
		"msisdn":  p.subscription.GetMsisdn(),
		"request": p.service.UrlNotifSub + "?" + q.Encode(),
		"trx_id":  trxId,
	}).Info("SUBSCRIPTION")

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
		"msisdn":      p.subscription.Msisdn,
		"response":    string(body),
		"trx_id":      trxId,
		"duration":    duration,
		"status_code": resp.StatusCode,
		"status_text": http.StatusText(resp.StatusCode),
	}).Info("SUBSCRIPTION")

	return body, nil
}

func (p *Portal) Unsubscription() ([]byte, error) {
	l := p.logger.Init("notif", true)

	start := time.Now()
	trxId := utils.GenerateTrxId()

	q := url.Values{}
	q.Add("telco", "telkomsel")
	q.Add("msisdn", p.subscription.Msisdn)
	q.Add("event", "unreg")

	req, err := http.NewRequest("GET", p.service.UrlNotifUnsub+"?"+q.Encode(), nil)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    5 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: tr,
	}

	p.logger.Writer(req)
	l.WithFields(logrus.Fields{
		"msisdn":  p.subscription.Msisdn,
		"request": p.service.UrlNotifUnsub + "?" + q.Encode(),
		"trx_id":  trxId,
	}).Info("UNSUBSCRIPTION")

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
		"msisdn":      p.subscription.Msisdn,
		"response":    string(body),
		"trx_id":      trxId,
		"duration":    duration,
		"status_code": resp.StatusCode,
		"status_text": http.StatusText(resp.StatusCode),
	}).Info("UNSUBSCRIPTION")

	return body, nil
}

func (p *Portal) Renewal() ([]byte, error) {
	l := p.logger.Init("notif", true)

	start := time.Now()
	trxId := utils.GenerateTrxId()

	q := url.Values{}
	q.Add("telco", "telkomsel")
	q.Add("msisdn", p.subscription.Msisdn)
	q.Add("event", "renewal")
	q.Add("password", p.pin)
	q.Add("package", "daily")
	q.Add("status", p.status)

	req, err := http.NewRequest("GET", p.service.UrlNotifRenewal+"?"+q.Encode(), nil)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    5 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: tr,
	}

	p.logger.Writer(req)
	l.WithFields(logrus.Fields{
		"msisdn":  p.subscription.Msisdn,
		"request": p.service.UrlNotifRenewal + "?" + q.Encode(),
		"trx_id":  trxId,
	}).Info("RENEWAL")

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
		"msisdn":      p.subscription.Msisdn,
		"response":    string(body),
		"trx_id":      trxId,
		"duration":    duration,
		"status_code": resp.StatusCode,
		"status_text": http.StatusText(resp.StatusCode),
	}).Info("RENEWAL")

	return body, nil
}

func (p *Portal) Callback() string {
	callbackUrl := p.service.UrlPortal + "?msisdn=" + p.subscription.Msisdn
	return callbackUrl
}
