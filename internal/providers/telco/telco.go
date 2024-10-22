package telco

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/idprm/go-xl-direct/internal/domain/entity"
	"github.com/idprm/go-xl-direct/internal/domain/model"
	"github.com/idprm/go-xl-direct/internal/logger"
	"github.com/idprm/go-xl-direct/internal/utils"
	"github.com/sirupsen/logrus"
)

var (
	TELCO_URL           string = utils.GetEnv("TELCO_URL")
	TELCO_CLIENT_ID     string = utils.GetEnv("TELCO_CLIENT_ID")
	TELCO_CLIENT_SECRET string = utils.GetEnv("TELCO_CLIENT_SECRET")
	TELCO_GRANT_TYPE    string = utils.GetEnv("TELCO_GRANT_TYPE")
)

type Telco struct {
	logger       *logger.Logger
	service      *entity.Service
	subscription *entity.Subscription
	session      *entity.Session
	verify       *entity.Verify
}

func NewTelco(
	logger *logger.Logger,
	service *entity.Service,
	subscription *entity.Subscription,
	session *entity.Session,
	verify *entity.Verify,
) *Telco {
	return &Telco{
		logger:       logger,
		service:      service,
		subscription: subscription,
		session:      session,
		verify:       verify,
	}
}

type ITelco interface {
	OAuth() ([]byte, error)
	CreateSubscription() ([]byte, error)
	ConfirmOTP(string) ([]byte, error)
	Refund() ([]byte, error)
	UnsubscribeSubscription() ([]byte, error)
	Notification() ([]byte, error)
}

func (t *Telco) OAuth() ([]byte, error) {

	var p = url.Values{}

	p.Add("client_id", TELCO_CLIENT_ID)
	p.Add("client_secret", TELCO_CLIENT_SECRET)
	p.Add("grant_type", TELCO_GRANT_TYPE)

	req, err := http.NewRequest(http.MethodPost, TELCO_URL+"/oauth2/token", strings.NewReader(p.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(p.Encode())))

	tr := &http.Transport{
		MaxIdleConns:       30,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: tr,
	}

	log.Println(req)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (t *Telco) CreateSubscription() ([]byte, error) {
	l := t.logger.Init("mt", true)

	start := time.Now()
	trxId := utils.GenerateTrxId()

	r := &model.CreateSubscriptionRequest{
		RequestId:      trxId,
		ProductId:      t.service.GetProductId(),
		UserIdentifier: t.verify.GetMsisdn(),
		Amount:         t.service.GetPriceToString(),
	}
	r.SetPartnerId(t.service.GetSidMt())

	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, TELCO_URL+"/subscription/create", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	oauth, err := t.OAuth()
	if err != nil {
		return nil, err
	}

	var respOauth model.OAuthResponse
	json.Unmarshal(oauth, &respOauth)

	var bearer = "Bearer " + respOauth.GetAccessToken()
	req.Header.Add("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")

	tr := &http.Transport{
		MaxIdleConns:       30,
		IdleConnTimeout:    60 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   60 * time.Second,
		Transport: tr,
	}

	t.logger.Writer(req)
	l.WithFields(logrus.Fields{
		"request": req,
		"trx_id":  trxId,
	}).Info("CREATE_SUB")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	duration := time.Since(start).Milliseconds()
	t.logger.Writer(string(body))
	l.WithFields(logrus.Fields{
		"response":    string(body),
		"trx_id":      trxId,
		"duration":    duration,
		"status_code": resp.StatusCode,
		"status_text": http.StatusText(resp.StatusCode),
	}).Info("CREATE_SUB")

	return body, nil
}

func (t *Telco) ConfirmOTP(pin string) ([]byte, error) {
	l := t.logger.Init("mt", true)

	start := time.Now()
	trxId := utils.GenerateTrxId()

	urlTelco := TELCO_URL + "/subscription/" + t.verify.GetMsisdn() + "/" + t.service.GetProductId() + "/otp/confirm"
	jsonData, err := json.Marshal(
		&model.ConfirmOTPRequest{
			RequestId: trxId,
			PIN:       pin,
		},
	)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, urlTelco, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	oauth, err := t.OAuth()
	if err != nil {
		return nil, err
	}

	var respOauth model.OAuthResponse
	json.Unmarshal(oauth, &respOauth)

	var bearer = "Bearer " + respOauth.GetAccessToken()
	req.Header.Add("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")

	tr := &http.Transport{
		MaxIdleConns:       30,
		IdleConnTimeout:    60 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   60 * time.Second,
		Transport: tr,
	}

	t.logger.Writer(req)
	l.WithFields(logrus.Fields{
		"msisdn":  t.verify.GetMsisdn(),
		"request": req,
		"trx_id":  trxId,
	}).Info("CONFIRM_OTP")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	duration := time.Since(start).Milliseconds()
	t.logger.Writer(string(body))
	l.WithFields(logrus.Fields{
		"msisdn":      t.verify.GetMsisdn(),
		"response":    string(body),
		"trx_id":      trxId,
		"duration":    duration,
		"status_code": resp.StatusCode,
		"status_text": http.StatusText(resp.StatusCode),
	}).Info("CONFIRM_OTP")

	return body, nil
}

func (t *Telco) Refund() ([]byte, error) {
	l := t.logger.Init("mt", true)

	start := time.Now()
	trxId := utils.GenerateTrxId()

	urlTelco := TELCO_URL + "/subscription/" + t.subscription.GetMsisdn() + "/" + t.service.GetProductId() + "/refund"
	jsonData, err := json.Marshal(
		&model.RefundRequest{
			RequestId:     trxId,
			TransactionId: t.subscription.GetLatestTrxId(),
		},
	)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, urlTelco, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	oauth, err := t.OAuth()
	if err != nil {
		return nil, err
	}

	var respOauth model.OAuthResponse
	json.Unmarshal(oauth, &respOauth)

	var bearer = "Bearer " + respOauth.GetAccessToken()
	req.Header.Add("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")

	tr := &http.Transport{
		MaxIdleConns:       30,
		IdleConnTimeout:    60 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   60 * time.Second,
		Transport: tr,
	}

	t.logger.Writer(req)
	l.WithFields(logrus.Fields{
		"msisdn":  t.subscription.GetMsisdn(),
		"request": req,
		"trx_id":  trxId,
	}).Info("REFUND")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	duration := time.Since(start).Milliseconds()
	t.logger.Writer(string(body))
	l.WithFields(logrus.Fields{
		"msisdn":      t.subscription.GetMsisdn(),
		"response":    string(body),
		"trx_id":      trxId,
		"duration":    duration,
		"status_code": resp.StatusCode,
		"status_text": http.StatusText(resp.StatusCode),
	}).Info("REFUND")

	return body, nil
}

func (t *Telco) UnsubscribeSubscription() ([]byte, error) {
	l := t.logger.Init("mt", true)

	start := time.Now()
	trxId := utils.GenerateTrxId()

	urlTelco := TELCO_URL + "/subscription/" + t.subscription.GetMsisdn() + "/" + t.service.GetProductId()
	jsonData, err := json.Marshal(
		&model.RefundRequest{
			RequestId:     trxId,
			TransactionId: t.subscription.GetLatestTrxId(),
		},
	)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodDelete, urlTelco, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	oauth, err := t.OAuth()
	if err != nil {
		return nil, err
	}

	var respOauth model.OAuthResponse
	json.Unmarshal(oauth, &respOauth)

	var bearer = "Bearer " + respOauth.GetAccessToken()
	req.Header.Add("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")

	tr := &http.Transport{
		MaxIdleConns:       30,
		IdleConnTimeout:    60 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   60 * time.Second,
		Transport: tr,
	}

	t.logger.Writer(req)
	l.WithFields(logrus.Fields{
		"msisdn":  t.subscription.GetMsisdn(),
		"request": req,
		"trx_id":  trxId,
	}).Info("UNSUB")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	duration := time.Since(start).Milliseconds()
	t.logger.Writer(string(body))
	l.WithFields(logrus.Fields{
		"msisdn":      t.subscription.GetMsisdn(),
		"response":    string(body),
		"trx_id":      trxId,
		"duration":    duration,
		"status_code": resp.StatusCode,
		"status_text": http.StatusText(resp.StatusCode),
	}).Info("UNSUB")

	return body, nil
}
