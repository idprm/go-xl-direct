package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/idprm/go-xl-direct/internal/domain/entity"
	"github.com/idprm/go-xl-direct/internal/utils"
)

type IncomingHandler struct {
}

func NewIncomingHandler() *IncomingHandler {
	return &IncomingHandler{}
}

var (
	APP_HOST     string = utils.GetEnv("APP_HOST")
	APP_URL      string = utils.GetEnv("APP_URL")
	TELCO_SENDER string = utils.GetEnv("TELCO_SENDER")
)

const (
	RMQ_DATA_TYPE          string = "application/json"
	RMQ_MO_EXCHANGE        string = "E_MO"
	RMQ_MO_QUEUE           string = "Q_MO"
	RMQ_NOTIF_EXCHANGE     string = "E_NOTIF"
	RMQ_NOTIF_QUEUE        string = "Q_NOTIF"
	RMQ_PB_MO_EXCHANGE     string = "E_POSTBACK_MO"
	RMQ_PB_MO_QUEUE        string = "Q_POSTBACK_MO"
	RMQ_PB_MT_EXCHANGE     string = "E_POSTBACK_MT"
	RMQ_PB_MT_QUEUE        string = "Q_POSTBACK_MT"
	RMQ_TRAFFIC_EXCHANGE   string = "E_TRAFFIC"
	RMQ_TRAFFIC_QUEUE      string = "Q_TRAFFIC"
	RMQ_DAILYPUSH_EXCHANGE string = "E_BQ_DAILYPUSH"
	RMQ_DAILYPUSH_QUEUE    string = "Q_BQ_DAILYPUSH"
	MT_FIRSTPUSH           string = "FIRSTPUSH"
	MT_RENEWAL             string = "RENEWAL"
	MT_UNSUB               string = "UNSUB"
	STATUS_SUCCESS         string = "SUCCESS"
	STATUS_FAILED          string = "FAILED"
	SUBJECT_FIRSTPUSH      string = "FIRSTPUSH"
	SUBJECT_DAILYPUSH      string = "DAILYPUSH"
	SUBJECT_RENEWAL        string = "RENEWAL"
	SUBJECT_UNSUB          string = "UNSUB"
	SUBJECT_RETRY          string = "RETRY"
	SUBJECT_PURGE          string = "PURGE"
)

var validate = validator.New()

func ValidateStruct(data interface{}) []*entity.ErrorResponse {
	var errors []*entity.ErrorResponse
	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element entity.ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
