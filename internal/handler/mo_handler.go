package handler

import (
	"encoding/json"
	"log"

	"github.com/idprm/go-xl-direct/internal/domain/model"
	"github.com/idprm/go-xl-direct/internal/logger"
	"github.com/idprm/go-xl-direct/internal/services"
	"github.com/wiliehidayat87/rmqp"
)

type MOHandler struct {
	rmq                 rmqp.AMQP
	logger              *logger.Logger
	blacklistService    services.IBlacklistService
	serviceService      services.IServiceService
	verifyService       services.IVerifyService
	subscriptionService services.ISubscriptionService
	transactionService  services.ITransactionService
	historyService      services.IHistoryService
	trafficService      services.ITrafficService
	req                 *model.NotificationRequest
}

func NewMOHandler(
	rmq rmqp.AMQP,
	logger *logger.Logger,
	blacklistService services.IBlacklistService,
	serviceService services.IServiceService,
	verifyService services.IVerifyService,
	subscriptionService services.ISubscriptionService,
	transactionService services.ITransactionService,
	historyService services.IHistoryService,
	trafficService services.ITrafficService,
	req *model.NotificationRequest,
) *MOHandler {
	return &MOHandler{
		rmq:                 rmq,
		logger:              logger,
		blacklistService:    blacklistService,
		serviceService:      serviceService,
		verifyService:       verifyService,
		subscriptionService: subscriptionService,
		transactionService:  transactionService,
		historyService:      historyService,
		trafficService:      trafficService,
		req:                 req,
	}
}

func (h *MOHandler) Firstpush() {

}

func (h *MOHandler) Unsub() {

}

func (h *MOHandler) Renewal() {
	jsonData, err := json.Marshal(h.req)
	if err != nil {
		log.Println(err.Error())
	}

	h.rmq.IntegratePublish(
		RMQ_RENEWAL_EXCHANGE,
		RMQ_RENEWAL_QUEUE,
		RMQ_DATA_TYPE, "", string(jsonData),
	)
}

func (h *MOHandler) Refund() {

}
