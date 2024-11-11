package handler

import (
	"encoding/json"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/idprm/go-xl-direct/internal/domain/entity"
	"github.com/idprm/go-xl-direct/internal/domain/model"
	"github.com/idprm/go-xl-direct/internal/logger"
	"github.com/idprm/go-xl-direct/internal/providers/telco"
	"github.com/idprm/go-xl-direct/internal/services"
	"github.com/idprm/go-xl-direct/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/wiliehidayat87/rmqp"
)

const (
	MO_REG   = "REG"
	MO_UNREG = "UNREG"
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
	APP_HOST   string = utils.GetEnv("APP_HOST")
	APP_URL    string = utils.GetEnv("APP_URL")
	TELCO_SDC  string = utils.GetEnv("TELCO_SDC")
	RMQ_PREFIX string = utils.GetEnv("RMQ_PREFIX")
)

var (
	RMQ_DATA_TYPE          string = "application/json"
	RMQ_MO_EXCHANGE        string = "E_" + RMQ_PREFIX + "_MO"
	RMQ_MO_QUEUE           string = "Q_" + RMQ_PREFIX + "_MO"
	RMQ_RENEWAL_EXCHANGE   string = "E_" + RMQ_PREFIX + "_RENEWAL"
	RMQ_RENEWAL_QUEUE      string = "Q_" + RMQ_PREFIX + "_RENEWAL"
	RMQ_REFUND_EXCHANGE    string = "E_" + RMQ_PREFIX + "_REFUND"
	RMQ_REFUND_QUEUE       string = "Q_" + RMQ_PREFIX + "_REFUND"
	RMQ_NOTIF_EXCHANGE     string = "E_" + RMQ_PREFIX + "_NOTIF"
	RMQ_NOTIF_QUEUE        string = "Q_" + RMQ_PREFIX + "_NOTIF"
	RMQ_PB_MO_EXCHANGE     string = "E_" + RMQ_PREFIX + "_POSTBACK_MO"
	RMQ_PB_MO_QUEUE        string = "Q_" + RMQ_PREFIX + "_POSTBACK_MO"
	RMQ_PB_MT_EXCHANGE     string = "E_" + RMQ_PREFIX + "_POSTBACK_MT"
	RMQ_PB_MT_QUEUE        string = "Q_" + RMQ_PREFIX + "_POSTBACK_MT"
	RMQ_TRAFFIC_EXCHANGE   string = "E_" + RMQ_PREFIX + "_TRAFFIC"
	RMQ_TRAFFIC_QUEUE      string = "Q_" + RMQ_PREFIX + "_TRAFFIC"
	RMQ_DAILYPUSH_EXCHANGE string = "E_" + RMQ_PREFIX + "_BQ_DAILYPUSH"
	RMQ_DAILYPUSH_QUEUE    string = "Q_" + RMQ_PREFIX + "_BQ_DAILYPUSH"
	MT_FIRSTPUSH           string = "FIRSTPUSH"
	MT_RENEWAL             string = "RENEWAL"
	MT_REFUND              string = "REFUND"
	MT_UNSUB               string = "UNSUB"
	STATUS_SUCCESS         string = "SUCCESS"
	STATUS_FAILED          string = "FAILED"
	SUBJECT_FIRSTPUSH      string = "FIRSTPUSH"
	SUBJECT_DAILYPUSH      string = "DAILYPUSH"
	SUBJECT_RENEWAL        string = "RENEWAL"
	SUBJECT_REFUND         string = "REFUND"
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

func (h *IncomingHandler) LandingPage(c *fiber.Ctx) error {
	paramService := strings.ToUpper(c.Params("service"))

	if !h.serviceService.IsServiceByCode(paramService) {
		return c.Status(fiber.StatusNotFound).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusNotFound,
				Message:    "service_not_found",
			},
		)
	}

	service, err := h.serviceService.GetServiceByCode(paramService)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusInternalServerError,
				Message:    err.Error(),
			},
		)
	}

	if service.IsEmagz() {
		return c.Render("emagz/sub", fiber.Map{
			"app_url":      APP_URL,
			"service_code": paramService,
		})
	}

	return c.Redirect(APP_URL)
}

func (h *IncomingHandler) UnsubPage(c *fiber.Ctx) error {
	paramService := strings.ToUpper(c.Params("service"))

	if !h.serviceService.IsServiceByCode(paramService) {
		return c.Status(fiber.StatusNotFound).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusNotFound,
				Message:    "service_not_found",
			},
		)
	}

	service, err := h.serviceService.GetServiceByCode(paramService)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusInternalServerError,
				Message:    err.Error(),
			},
		)
	}

	if service.IsEmagz() {
		return c.Render("emagz/unsub", fiber.Map{
			"app_url":      APP_URL,
			"service_code": paramService,
		})
	}

	return c.Redirect(APP_URL)
}

func (h *IncomingHandler) UnRegPage(c *fiber.Ctx) error {

	req := new(model.WebUnRegRequest)

	err := c.QueryParser(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusBadRequest,
				Message:    err.Error(),
			},
		)
	}

	paramService := strings.ToUpper(c.Params("service"))

	if !h.serviceService.IsServiceByCode(paramService) {
		return c.Status(fiber.StatusNotFound).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusNotFound,
				Message:    "service_not_found",
			},
		)
	}

	service, err := h.serviceService.GetServiceByCode(paramService)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusInternalServerError,
				Message:    err.Error(),
			},
		)
	}

	if !h.subscriptionService.IsActiveSubscription(service.GetId(), req.GetMsisdn()) {
		return c.Redirect(service.GetUrlPortal())
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
		return c.Redirect(service.GetUrlPortal())
	}

	return c.Redirect(service.GetUrlPortal())
}

func (h *IncomingHandler) Campaign(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(
		&model.WebResponse{
			Error:       false,
			StatusCode:  fiber.StatusOK,
			Message:     "OK",
			RedirectUrl: APP_URL,
		},
	)
}

func (h *IncomingHandler) MessageOriginated(c *fiber.Ctx) error {
	l := h.logger.Init("mo", true)

	c.Accepts("application/json")

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

	l.WithFields(logrus.Fields{"request": req}).Info("MO")

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
	l := h.logger.Init("web", true)

	c.Accepts("application/json")

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

	errors := ValidateStruct(*req)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	l.WithFields(logrus.Fields{"request": req}).Info("CREATE_SUB")

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

	if !resp.IsSuccess() {
		return c.Status(fiber.StatusOK).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusOK,
				Message:    resp.GetErrorDescription(),
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
	l := h.logger.Init("web", true)

	c.Accepts("application/json")

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

	errors := ValidateStruct(*req)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	l.WithFields(logrus.Fields{"request": req}).Info("CONFIRM_OTP")

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

	if resp.IsSuccess() {
		return c.Status(fiber.StatusOK).JSON(
			&model.WebResponse{
				Error:       false,
				StatusCode:  fiber.StatusOK,
				Message:     resp.GetStatus(),
				RedirectUrl: service.GetUrlPortal(),
			},
		)
	}

	if !resp.IsSuccess() {
		return c.Status(fiber.StatusOK).JSON(
			&model.WebResponse{
				Error:      true,
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

func (h *IncomingHandler) Refund(c *fiber.Ctx) error {
	l := h.logger.Init("web", true)

	c.Accepts("application/json")

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

	errors := ValidateStruct(*req)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	l.WithFields(logrus.Fields{"request": req}).Info("REFUND")

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
	l := h.logger.Init("web", true)

	c.Accepts("application/json")

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

	errors := ValidateStruct(*req)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	l.WithFields(logrus.Fields{"request": req}).Info("UNSUBSCRIBE")

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
		return c.Status(fiber.StatusOK).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusOK,
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
			Error:       false,
			StatusCode:  fiber.StatusOK,
			Message:     resp.GetStatus(),
			RedirectUrl: service.GetUrlPortal(),
		},
	)
}
