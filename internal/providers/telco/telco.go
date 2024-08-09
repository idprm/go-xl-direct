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
)

var (
	TELCO_URL           string = utils.GetEnv("TELCO_URL")
	TELCO_CLIENT_ID     string = utils.GetEnv("TELCO_CLIENT_ID")
	TELCO_SECRET_SECRET string = utils.GetEnv("TELCO_SECRET_SECRET")
	TELCO_GRANT_TYPE    string = utils.GetEnv("TELCO_GRANT_TYPE")
)

type Telco struct {
	logger       *logger.Logger
	session      *entity.Session
	subscription *entity.Subscription
	service      *entity.Service
}

func NewTelco(
	logger *logger.Logger,
	session *entity.Session,
	subscription *entity.Subscription,
	service *entity.Service,
) *Telco {
	return &Telco{
		logger:       logger,
		session:      session,
		subscription: subscription,
		service:      service,
	}
}

type ITelco interface {
	OAuth() ([]byte, error)
	CreateSubscription() ([]byte, error)
	ConfirmOTP() ([]byte, error)
	Refund() ([]byte, error)
	UnsubscribeSubscription() ([]byte, error)
	Notification() ([]byte, error)
}

func (t *Telco) OAuth() ([]byte, error) {
	var p = url.Values{}

	p.Add("client_id", TELCO_CLIENT_ID)
	p.Add("client_secret", TELCO_SECRET_SECRET)
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
	jsonData, err := json.Marshal(
		&model.CreateSubscriptionRequest{
			RequestId:      "",
			ProductId:      "",
			UserIdentifier: "",
			Amount:         "",
		},
	)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, TELCO_URL+"/subscription/create", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	var bearer = "Bearer " + t.session.GetAccessToken()
	req.Header.Add("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")

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

func (t *Telco) ConfirmOTP() ([]byte, error) {

	urlTelco := TELCO_URL + "/subscription/{msisdn}/{productId}/otp/confirm"
	jsonData, err := json.Marshal(
		&model.ConfirmOTPRequest{
			RequestId: "",
			PIN:       "",
		},
	)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, urlTelco, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	var bearer = "Bearer " + t.session.GetAccessToken()
	req.Header.Add("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")

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

func (t *Telco) Refund() ([]byte, error) {

	urlTelco := TELCO_URL + "/subscription/{msisdn}/{productId}/refund"
	jsonData, err := json.Marshal(
		&model.RefundRequest{
			RequestId:     "",
			TransactionId: "",
		},
	)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, urlTelco, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	var bearer = "Bearer " + t.session.GetAccessToken()
	req.Header.Add("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")

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

func (t *Telco) UnsubscribeSubscription() ([]byte, error) {
	urlTelco := TELCO_URL + "/subscription/{msisdn}/{productId}"
	jsonData, err := json.Marshal(
		&model.RefundRequest{
			RequestId:     "",
			TransactionId: "",
		},
	)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodDelete, urlTelco, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	var bearer = "Bearer " + t.session.GetAccessToken()
	req.Header.Add("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")

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

func (t *Telco) Notification() ([]byte, error) {
	urlTelco := TELCO_URL + "/subscription/{msisdn}/{productId}"
	jsonData, err := json.Marshal(
		&model.RefundRequest{
			RequestId:     "",
			TransactionId: "",
		},
	)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodDelete, urlTelco, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	var bearer = "Bearer " + t.session.GetAccessToken()
	req.Header.Add("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")

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
