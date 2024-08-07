package telco

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/idprm/go-xl-direct/internal/domain/entity"
	"github.com/idprm/go-xl-direct/internal/logger"
	"github.com/idprm/go-xl-direct/internal/utils"
)

var (
	TELCO_URL_AUTH string = utils.GetEnv("TELCO_URL_AUTH")
	TELCO_KEY      string = utils.GetEnv("TELCO_KEY")
	TELCO_SECRET   string = utils.GetEnv("TELCO_SECRET")
	TELCO_CPNAME   string = utils.GetEnv("TELCO_CPNAME")
	TELCO_CPID     string = utils.GetEnv("TELCO_CPID")
	TELCO_PWD      string = utils.GetEnv("TELCO_PWD")
	TELCO_SENDER   string = utils.GetEnv("TELCO_SENDER")
)

type Telco struct {
	logger       *logger.Logger
	subscription *entity.Subscription
	service      *entity.Service
	content      *entity.Content
}

func NewTelco(
	logger *logger.Logger,
	subscription *entity.Subscription,
	service *entity.Service,
	content *entity.Content,
) *Telco {
	return &Telco{
		logger:       logger,
		subscription: subscription,
		service:      service,
		content:      content,
	}
}

type ITelco interface {
	Test() ([]byte, error)
	OAuth() ([]byte, error)
	CreateSubscription() ([]byte, error)
	ConfirmOTP() ([]byte, error)
	Refund() ([]byte, error)
	UnsubscribeSubscription() ([]byte, error)
	Notification() ([]byte, error)
}

func (t *Telco) Test() ([]byte, error) {
	params := url.Values{}
	params.Add("client_id", "3S7QIae30ToXBghLAdoQY8V8rWnlYqiA")
	params.Add("client_secret", "121c4kc9saaCKF9Hd7F3zesRUjOmSJs8")
	params.Add("grant_type", "client_credentials")
	resp, err := http.PostForm("https://staging-sdp.xlaxiata.id/dcb-nongoogle/oauth2/token",
		params)
	if err != nil {
		log.Printf("Request Failed: %s", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return nil, err
	}
	// Unmarshal result
	// post := Post{}
	// err = json.Unmarshal(body, &post)
	// if err != nil {
	// 	log.Printf("Reading body failed: %s", err)
	// 	return nil, err
	// }

	// log.Printf("Post added with ID %d", post.ID)
	return body, nil
}

func (t *Telco) OAuth() ([]byte, error) {
	var p = url.Values{
		"client_id":     {"3S7QIae30ToXBghLAdoQY8V8rWnlYqiA"},
		"client_secret": {"121c4kc9saaCKF9Hd7F3zesRUjOmSJs8"},
		"grant_type":    {"client_credentials"},
	}
	// p.Add("client_id", "3S7QIae30ToXBghLAdoQY8V8rWnlYqiA")
	// p.Add("client_secret", "121c4kc9saaCKF9Hd7F3zesRUjOmSJs8")
	// p.Add("grant_type", "client_credentials")

	req, err := http.NewRequest(http.MethodPost, "https://staging-sdp.xlaxiata.id/dcb-nongoogle/oauth2/token", strings.NewReader(p.Encode()))
	if err != nil {
		log.Println(err)
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
		return nil, errors.New(err.Error())
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return body, nil
}

func (t *Telco) CreateSubscription() ([]byte, error) {
	return nil, nil
}

func (t *Telco) ConfirmOTP() ([]byte, error) {
	return nil, nil
}

func (t *Telco) Refund() ([]byte, error) {
	return nil, nil
}

func (t *Telco) UnsubscribeSubscription() ([]byte, error) {
	return nil, nil
}

func (t *Telco) Notification() ([]byte, error) {
	return nil, nil
}
