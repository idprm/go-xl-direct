package handler

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/idprm/go-xl-direct/internal/domain/entity"
	"github.com/idprm/go-xl-direct/internal/domain/model"
	"github.com/idprm/go-xl-direct/internal/logger"
	"github.com/idprm/go-xl-direct/internal/services"
	"github.com/idprm/go-xl-direct/internal/utils"
	"github.com/wiliehidayat87/rmqp"
)

type IncomingHandler struct {
	rmq            rmqp.AMQP
	logger         *logger.Logger
	serviceService services.IServiceService
	verifyService  services.IVerifyService
}

func NewIncomingHandler(
	rmq rmqp.AMQP,
	logger *logger.Logger,
	serviceService services.IServiceService,
	verifyService services.IVerifyService,
) *IncomingHandler {
	return &IncomingHandler{
		rmq:            rmq,
		logger:         logger,
		serviceService: serviceService,
		verifyService:  verifyService,
	}
}

var (
	APP_HOST  string = utils.GetEnv("APP_HOST")
	APP_URL   string = utils.GetEnv("APP_URL")
	TELCO_SDC string = utils.GetEnv("TELCO_SDC")
)

const (
	RMQ_DATA_TYPE          string = "application/json"
	RMQ_MO_EXCHANGE        string = "E_MO"
	RMQ_MO_QUEUE           string = "Q_MO"
	RMQ_RENEWAL_EXCHANGE   string = "E_RENEWAL"
	RMQ_RENEWAL_QUEUE      string = "Q_RENEWAL"
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

	if !h.serviceService.IsServiceByCode(req.GetProductId()) {
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
	return c.Status(fiber.StatusOK).JSON(
		&model.WebResponse{
			Error:      false,
			StatusCode: fiber.StatusOK,
			Message:    "Successful",
		},
	)
}

func (h *IncomingHandler) ConfirmOTP(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(
		&model.WebResponse{
			Error:      false,
			StatusCode: fiber.StatusOK,
			Message:    "Successful",
		},
	)
}
