package handler

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/idprm/go-xl-direct/internal/domain/entity"
	"github.com/idprm/go-xl-direct/internal/domain/model"
	"github.com/idprm/go-xl-direct/internal/logger"
	"github.com/idprm/go-xl-direct/internal/providers/telco"
	"github.com/idprm/go-xl-direct/internal/services"
	"github.com/idprm/go-xl-direct/internal/utils"
	"github.com/wiliehidayat87/rmqp"
)

const (
	MO_REG   = "REG"
	MO_UNREG = "UNREG"
	MO_OFF   = "OFF"
)

type IncomingHandler struct {
	rmq                 rmqp.AMQP
	logger              *logger.Logger
	serviceService      services.IServiceService
	subscriptionService services.ISubscriptionService
	verifyService       services.IVerifyService
}

func NewIncomingHandler(
	rmq rmqp.AMQP,
	logger *logger.Logger,
	serviceService services.IServiceService,
	subscriptionService services.ISubscriptionService,
	verifyService services.IVerifyService,
) *IncomingHandler {
	return &IncomingHandler{
		rmq:                 rmq,
		logger:              logger,
		serviceService:      serviceService,
		subscriptionService: subscriptionService,
		verifyService:       verifyService,
	}
}

var (
	APP_HOST  string = utils.GetEnv("APP_HOST")
	APP_URL   string = utils.GetEnv("APP_URL")
	TELCO_SDC string = utils.GetEnv("TELCO_SDC")
)

const (
	RMQ_DATA_TYPE        string = "application/json"
	RMQ_MO_EXCHANGE      string = "E_MO"
	RMQ_MO_QUEUE         string = "Q_MO"
	RMQ_RENEWAL_EXCHANGE string = "E_RENEWAL"
	RMQ_RENEWAL_QUEUE    string = "Q_RENEWAL"
	RMQ_NOTIF_EXCHANGE   string = "E_NOTIF"
	RMQ_NOTIF_QUEUE      string = "Q_NOTIF"
	RMQ_PB_MO_EXCHANGE   string = "E_POSTBACK_MO"
	RMQ_PB_MO_QUEUE      string = "Q_POSTBACK_MO"
	RMQ_PB_MT_EXCHANGE   string = "E_POSTBACK_MT"
	RMQ_PB_MT_QUEUE      string = "Q_POSTBACK_MT"
	MT_FIRSTPUSH         string = "FIRSTPUSH"
	MT_RENEWAL           string = "RENEWAL"
	MT_UNSUB             string = "UNSUB"
	STATUS_SUCCESS       string = "SUCCESS"
	STATUS_FAILED        string = "FAILED"
	SUBJECT_FIRSTPUSH    string = "FIRSTPUSH"
	SUBJECT_DAILYPUSH    string = "DAILYPUSH"
	SUBJECT_RENEWAL      string = "RENEWAL"
	SUBJECT_UNSUB        string = "UNSUB"
	SUBJECT_RETRY        string = "RETRY"
	SUBJECT_PURGE        string = "PURGE"
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

func (h *IncomingHandler) Campaign(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(
		&model.WebResponse{
			Error:      false,
			StatusCode: fiber.StatusOK,
			Message:    "OK",
		},
	)
}

func (h *IncomingHandler) MessageOriginated(c *fiber.Ctx) error {

	req := new(model.NotificationRequest)

	err := c.BodyParser(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusBadRequest,
				Message:    err.Error(),
			},
		)
	}

	if !h.serviceService.IsServiceByProductId(req.GetProductId()) {
		return c.Status(fiber.StatusNotFound).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusNotFound,
				Message:    "service_unavailable",
			},
		)
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusBadGateway,
				Message:    err.Error(),
			},
		)
	}

	h.rmq.IntegratePublish(
		RMQ_MO_EXCHANGE,
		RMQ_MO_QUEUE,
		RMQ_DATA_TYPE, "", string(jsonData),
	)

	return c.Status(fiber.StatusOK).JSON(
		&model.WebResponse{
			Error:      false,
			StatusCode: fiber.StatusOK,
			Message:    "Successful",
		},
	)
}

func (h *IncomingHandler) CreateSubscription(c *fiber.Ctx) error {

	req := new(model.WebSubRequest)

	err := c.BodyParser(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusBadRequest,
				Message:    err.Error(),
			},
		)
	}

	if !h.serviceService.IsServiceByCode(req.GetService()) {
		return c.Status(fiber.StatusNotFound).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusNotFound,
				Message:    "service_unavailable",
			},
		)
	}

	service, err := h.serviceService.GetServiceByCode(req.GetService())
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusBadGateway,
				Message:    err.Error(),
			},
		)
	}

	if c.Get("Cf-Connecting-Ip") != "" {
		req.SetIpAddress(c.Get("Cf-Connecting-Ip"))
	} else {
		req.SetIpAddress(c.Get("X-Forwarded-For"))
	}

	verify := &entity.Verify{
		Msisdn:    req.GetMsisdn(),
		Service:   service,
		IpAddress: req.IpAddress,
	}

	t := telco.NewTelco(
		h.logger,
		service,
		&entity.Subscription{},
		&entity.Session{},
		verify,
	)

	mt, err := t.CreateSubscription()
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusBadGateway,
				Message:    err.Error(),
			},
		)
	}

	var resp model.TelcoResponse
	json.Unmarshal(mt, &resp)

	if resp.IsSuccess() {

		verify.SetTrxId(resp.GetTransactionId())
		h.verifyService.Set(verify)

		return c.Status(fiber.StatusOK).JSON(
			&model.WebResponse{
				Error:      false,
				StatusCode: fiber.StatusOK,
				Message:    resp.GetStatus(),
			},
		)
	}

	return c.Status(fiber.StatusBadGateway).JSON(
		&model.WebResponse{
			Error:      true,
			StatusCode: fiber.StatusBadGateway,
			Message:    "error_bad_gateway",
		},
	)
}

func (h *IncomingHandler) ConfirmOTP(c *fiber.Ctx) error {

	req := new(model.WebOTPRequest)

	err := c.BodyParser(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusBadRequest,
				Message:    err.Error(),
			},
		)
	}

	if !h.serviceService.IsServiceByCode(req.GetService()) {
		return c.Status(fiber.StatusNotFound).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusNotFound,
				Message:    "service_unavailable",
			},
		)
	}

	service, err := h.serviceService.GetServiceByCode(req.GetService())
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusBadGateway,
				Message:    err.Error(),
			},
		)
	}

	verify, err := h.verifyService.Get(req.GetMsisdn())
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusBadGateway,
				Message:    err.Error(),
			},
		)
	}

	t := telco.NewTelco(
		h.logger,
		service,
		&entity.Subscription{},
		&entity.Session{},
		verify,
	)

	mt, err := t.ConfirmOTP(req.GetPin())
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusBadGateway,
				Message:    err.Error(),
			},
		)
	}

	var resp model.TelcoResponse
	json.Unmarshal(mt, &resp)

	if !resp.IsSuccess() {
		return c.Status(fiber.StatusOK).JSON(
			&model.WebResponse{
				Error:      false,
				StatusCode: fiber.StatusOK,
				Message:    resp.GetStatus(),
			},
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		&model.WebResponse{
			Error:      true,
			StatusCode: fiber.StatusOK,
			Message:    resp.GetStatus(),
		},
	)
}

func (h *IncomingHandler) Refund(c *fiber.Ctx) error {

	req := new(model.WebSubRequest)

	err := c.BodyParser(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusBadRequest,
				Message:    err.Error(),
			},
		)
	}

	if !h.serviceService.IsServiceByCode(req.GetService()) {
		return c.Status(fiber.StatusNotFound).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusNotFound,
				Message:    "service_unavailable",
			},
		)
	}

	service, err := h.serviceService.GetServiceByCode(req.GetService())
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusBadGateway,
				Message:    err.Error(),
			},
		)
	}

	if !h.subscriptionService.IsActiveSubscription(service.GetId(), req.GetMsisdn()) {
		return c.Status(fiber.StatusNotFound).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusNotFound,
				Message:    "msisdn_not_found",
			},
		)
	}

	subscription, err := h.subscriptionService.SelectSubscription(service.GetId(), req.GetMsisdn())
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusBadGateway,
				Message:    err.Error(),
			},
		)
	}

	t := telco.NewTelco(
		h.logger,
		service,
		subscription,
		&entity.Session{},
		&entity.Verify{},
	)

	mt, err := t.UnsubscribeSubscription()
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusBadGateway,
				Message:    err.Error(),
			},
		)
	}

	var resp model.TelcoResponse
	json.Unmarshal(mt, &resp)

	if !resp.IsSuccess() {
		return c.Status(fiber.StatusOK).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusOK,
				Message:    resp.GetStatus(),
			},
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		&model.WebResponse{
			Error:      false,
			StatusCode: fiber.StatusOK,
			Message:    resp.GetStatus(),
		},
	)

}

func (h *IncomingHandler) Unsubscribe(c *fiber.Ctx) error {

	req := new(model.WebSubRequest)

	err := c.BodyParser(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusBadRequest,
				Message:    err.Error(),
			},
		)
	}

	if !h.serviceService.IsServiceByCode(req.GetService()) {
		return c.Status(fiber.StatusNotFound).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusNotFound,
				Message:    "service_unavailable",
			},
		)
	}

	service, err := h.serviceService.GetServiceByCode(req.GetService())
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusBadGateway,
				Message:    err.Error(),
			},
		)
	}

	if !h.subscriptionService.IsActiveSubscription(service.GetId(), req.GetMsisdn()) {
		return c.Status(fiber.StatusNotFound).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusNotFound,
				Message:    "msisdn_not_found_or_already_unsub",
			},
		)
	}

	subscription, err := h.subscriptionService.SelectSubscription(service.GetId(), req.GetMsisdn())
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusBadGateway,
				Message:    err.Error(),
			},
		)
	}

	t := telco.NewTelco(
		h.logger,
		service,
		subscription,
		&entity.Session{},
		&entity.Verify{},
	)

	mt, err := t.UnsubscribeSubscription()
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusBadGateway,
				Message:    err.Error(),
			},
		)
	}

	var resp model.TelcoResponse
	json.Unmarshal(mt, &resp)

	if !resp.IsSuccess() {
		return c.Status(fiber.StatusOK).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusOK,
				Message:    resp.GetStatus(),
			},
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		&model.WebResponse{
			Error:      false,
			StatusCode: fiber.StatusOK,
			Message:    resp.GetStatus(),
		},
	)
}
