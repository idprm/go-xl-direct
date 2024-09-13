package handler

import (
	"log"

	"github.com/idprm/go-xl-direct/internal/domain/entity"
	"github.com/idprm/go-xl-direct/internal/domain/model"
	"github.com/idprm/go-xl-direct/internal/logger"
	"github.com/idprm/go-xl-direct/internal/services"
	"github.com/idprm/go-xl-direct/internal/utils"
	"github.com/wiliehidayat87/rmqp"
)

type RefundHandler struct {
	rmq                 rmqp.AMQP
	logger              *logger.Logger
	serviceService      services.IServiceService
	subscriptionService services.ISubscriptionService
	transactionService  services.ITransactionService
	req                 *model.NotificationRequest
}

func NewRefundHandler(
	rmq rmqp.AMQP,
	logger *logger.Logger,
	serviceService services.IServiceService,
	subscriptionService services.ISubscriptionService,
	transactionService services.ITransactionService,
	req *model.NotificationRequest,
) *RefundHandler {
	return &RefundHandler{
		rmq:                 rmq,
		logger:              logger,
		serviceService:      serviceService,
		subscriptionService: subscriptionService,
		transactionService:  transactionService,
		req:                 req,
	}
}

func (h *RefundHandler) Refund() {

	trxId := utils.GenerateTrxId()

	service, err := h.getService()
	if err != nil {
		log.Println(err)
	}

	h.subscriptionService.UpdateLatest(
		&entity.Subscription{
			ServiceID:     service.GetId(),
			Msisdn:        h.req.GetUserIdentifier(),
			LatestTrxId:   trxId,
			LatestSubject: SUBJECT_REFUND,
			LatestStatus:  h.req.GetStatus(),
			LatestKeyword: SUBJECT_REFUND,
		},
	)

	transSuccess := &entity.Transaction{
		TxID:           trxId,
		ServiceID:      service.GetId(),
		Msisdn:         h.req.GetUserIdentifier(),
		SubID:          h.req.GetSubscriptionId(),
		Channel:        "",
		Adnet:          "",
		Keyword:        SUBJECT_REFUND,
		Amount:         h.req.GetAmount(),
		PIN:            "",
		Status:         h.req.Status,
		StatusCode:     "",
		StatusDetail:   "",
		Subject:        SUBJECT_REFUND,
		Payload:        "",
		CampKeyword:    "",
		CampSubKeyword: "",
		IpAddress:      "",
	}

	h.transactionService.SaveTransaction(transSuccess)
}

func (h *RefundHandler) getService() (*entity.Service, error) {
	return h.serviceService.GetServiceByProductId(h.req.GetProductId())
}
