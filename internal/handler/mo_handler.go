package handler

import (
	"encoding/json"
	"log"
	"time"

	"github.com/idprm/go-xl-direct/internal/domain/entity"
	"github.com/idprm/go-xl-direct/internal/domain/model"
	"github.com/idprm/go-xl-direct/internal/logger"
	"github.com/idprm/go-xl-direct/internal/services"
	"github.com/idprm/go-xl-direct/internal/utils"
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

	trxId := utils.GenerateTrxId()

	service, err := h.getService()
	if err != nil {
		log.Println(err)
	}

	verify, err := h.verifyService.Get(h.req.GetTransactionId())
	if err != nil {
		log.Println(err)
	}

	subscription := &entity.Subscription{
		ServiceID:     service.GetId(),
		Category:      service.GetCategory(),
		Msisdn:        h.req.GetUserIdentifier(),
		SubID:         h.req.GetSubscriptionId(),
		LatestTrxId:   trxId,
		LatestKeyword: MO_REG + " " + service.GetCode(),
		LatestSubject: SUBJECT_FIRSTPUSH,
		Channel:       "",
		IsActive:      true,
	}

	if verify != nil {
		subscription.Channel = verify.Channel
		subscription.Adnet = verify.Adnet
		subscription.PubID = verify.PubID
		subscription.AffSub = verify.AffSub
		subscription.CampKeyword = verify.Keyword
		subscription.CampSubKeyword = verify.SubKeyword
		subscription.IpAddress = verify.IpAddress
	}

	if h.IsSub() {
		h.subscriptionService.UpdateEnable(subscription)
	} else {
		h.subscriptionService.SaveSubscription(subscription)
	}

	if h.req.IsActive() {
		subSuccess := &entity.Subscription{
			ServiceID:            service.GetId(),
			Msisdn:               h.req.GetUserIdentifier(),
			LatestTrxId:          h.req.GetTransactionId(),
			LatestSubject:        SUBJECT_FIRSTPUSH,
			LatestStatus:         STATUS_SUCCESS,
			LatestPIN:            "",
			Amount:               h.req.GetAmount(),
			RenewalAt:            h.req.GetNextRenewalDate(),
			ChargeAt:             time.Now(),
			Success:              1,
			IsRetry:              false,
			TotalFirstpush:       1,
			TotalAmountFirstpush: h.req.GetAmount(),
			LatestPayload:        "",
		}
		if verify != nil {
			subSuccess.Channel = verify.Channel
			subSuccess.LatestPIN = verify.PIN
		}
		h.subscriptionService.UpdateSuccess(subSuccess)

		transSuccess := &entity.Transaction{
			TxID:         trxId,
			ServiceID:    service.GetId(),
			Msisdn:       h.req.GetUserIdentifier(),
			SubID:        h.req.GetSubscriptionId(),
			Channel:      "",
			Keyword:      MO_REG + " " + service.GetCode(),
			Amount:       h.req.GetAmount(),
			PIN:          "",
			Status:       STATUS_SUCCESS,
			StatusCode:   "",
			StatusDetail: "",
			Subject:      SUBJECT_FIRSTPUSH,
			Payload:      "",
			IpAddress:    "",
		}
		if verify != nil {
			transSuccess.Channel = verify.Channel
			transSuccess.CampKeyword = verify.Keyword
			transSuccess.CampSubKeyword = verify.SubKeyword
			transSuccess.Adnet = verify.Adnet
			transSuccess.PubID = verify.PubID
			transSuccess.AffSub = verify.AffSub
			transSuccess.PIN = verify.PIN
			transSuccess.IpAddress = verify.GetIpAddress()
		}
		h.transactionService.SaveTransaction(transSuccess)

		historySuccess := &entity.History{
			ServiceID: service.GetId(),
			Msisdn:    h.req.GetUserIdentifier(),
			SubID:     h.req.GetSubscriptionId(),
			Channel:   "",
			Keyword:   MO_REG + " " + service.GetCode(),
			Subject:   SUBJECT_FIRSTPUSH,
			Status:    STATUS_SUCCESS,
		}

		if verify != nil {
			historySuccess.Channel = verify.Channel
			historySuccess.Adnet = verify.Adnet
			historySuccess.IpAddress = verify.GetIpAddress()
		}

		h.historyService.SaveHistory(historySuccess)

		// insert to rabbitmq
		jsonData, _ := json.Marshal(
			&entity.NotifParamsRequest{
				Service:      service,
				Subscription: subscription,
				Action:       "SUB",
				Pin:          "",
			},
		)
		h.rmq.IntegratePublish(
			RMQ_NOTIF_EXCHANGE,
			RMQ_NOTIF_QUEUE,
			RMQ_DATA_TYPE,
			"",
			string(jsonData),
		)

	} else {
		subFailed := &entity.Subscription{
			ServiceID:     service.GetId(),
			Msisdn:        h.req.GetUserIdentifier(),
			LatestTrxId:   trxId,
			LatestSubject: SUBJECT_FIRSTPUSH,
			LatestStatus:  STATUS_FAILED,
			RenewalAt:     time.Now().AddDate(0, 0, 1),
			RetryAt:       time.Now(),
			Failed:        1,
			IsRetry:       true,
			LatestPayload: "",
		}
		h.subscriptionService.UpdateFailed(subFailed)

		transFailed := &entity.Transaction{
			TxID:         trxId,
			ServiceID:    service.GetId(),
			Msisdn:       h.req.GetUserIdentifier(),
			SubID:        h.req.GetSubscriptionId(),
			Channel:      "",
			Keyword:      MO_REG + " " + service.GetCode(),
			Status:       STATUS_FAILED,
			StatusCode:   "",
			StatusDetail: "",
			Subject:      SUBJECT_FIRSTPUSH,
			Payload:      "",
			IpAddress:    "",
		}
		if verify != nil {
			transFailed.Channel = verify.Channel
			transFailed.CampKeyword = verify.Keyword
			transFailed.CampSubKeyword = verify.SubKeyword
			transFailed.Adnet = verify.Adnet
			transFailed.PubID = verify.PubID
			transFailed.AffSub = verify.AffSub
			transFailed.PIN = verify.PIN
			transFailed.IpAddress = verify.GetIpAddress()
		}
		h.transactionService.SaveTransaction(transFailed)

		historyFailed := &entity.History{
			ServiceID: service.GetId(),
			Msisdn:    h.req.GetUserIdentifier(),
			SubID:     h.req.GetSubscriptionId(),
			Channel:   "",
			Keyword:   MO_REG + " " + service.GetCode(),
			Subject:   SUBJECT_FIRSTPUSH,
			Status:    STATUS_FAILED,
		}
		if verify != nil {
			historyFailed.Adnet = verify.Adnet
			historyFailed.Channel = verify.Channel
			historyFailed.IpAddress = verify.GetIpAddress()
		}
		h.historyService.SaveHistory(historyFailed)

	}

	if verify != nil {
		// insert to rabbitmq
		jsonData, _ := json.Marshal(
			&entity.PostbackParamsRequest{
				Verify:       verify,
				Subscription: subscription,
				Service:      service,
				Action:       "MO",
			},
		)
		h.rmq.IntegratePublish(
			RMQ_PB_MO_EXCHANGE,
			RMQ_PB_MO_QUEUE,
			RMQ_DATA_TYPE,
			"",
			string(jsonData),
		)
	}
}

func (h *MOHandler) Unsub() {

	trxId := utils.GenerateTrxId()

	service, err := h.getService()
	if err != nil {
		log.Println(err)
	}

	subscription := &entity.Subscription{
		ServiceID:     service.GetId(),
		Msisdn:        h.req.GetUserIdentifier(),
		LatestTrxId:   trxId,
		LatestKeyword: MO_UNREG + " " + service.GetCode(),
		LatestSubject: SUBJECT_UNSUB,
		LatestStatus:  STATUS_SUCCESS,
		UnsubAt:       time.Now(),
		IpAddress:     "",
		IsRetry:       false,
		IsActive:      false,
	}
	h.subscriptionService.UpdateDisable(subscription)

	transaction := &entity.Transaction{
		TxID:         trxId,
		ServiceID:    service.GetId(),
		Msisdn:       h.req.GetUserIdentifier(),
		Channel:      "WAP",
		Adnet:        "",
		Keyword:      MO_UNREG + " " + service.GetCode(),
		Status:       STATUS_SUCCESS,
		StatusCode:   "",
		StatusDetail: "",
		Subject:      SUBJECT_UNSUB,
		Payload:      "",
		IpAddress:    "",
	}
	h.transactionService.SaveTransaction(transaction)

	history := &entity.History{
		ServiceID: service.GetId(),
		Msisdn:    h.req.GetUserIdentifier(),
		Channel:   "WAP",
		Adnet:     "",
		Keyword:   MO_UNREG + " " + service.GetCode(),
		Subject:   SUBJECT_UNSUB,
		Status:    STATUS_SUCCESS,
		IpAddress: "",
	}
	h.historyService.SaveHistory(history)

	// insert to rabbitmq
	jsonDataNotif, _ := json.Marshal(
		&entity.NotifParamsRequest{
			Service:      service,
			Subscription: subscription,
			Action:       "UNSUB",
		},
	)

	h.rmq.IntegratePublish(
		RMQ_NOTIF_EXCHANGE,
		RMQ_NOTIF_QUEUE,
		RMQ_DATA_TYPE,
		"",
		string(jsonDataNotif),
	)
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
	jsonData, err := json.Marshal(h.req)
	if err != nil {
		log.Println(err.Error())
	}

	h.rmq.IntegratePublish(
		RMQ_REFUND_EXCHANGE,
		RMQ_REFUND_QUEUE,
		RMQ_DATA_TYPE, "", string(jsonData),
	)
}

func (h *MOHandler) IsSub() bool {
	service, err := h.getService()
	if err != nil {
		log.Println(err)
	}
	return h.subscriptionService.IsSubscription(service.GetId(), h.req.GetUserIdentifier())
}

func (h *MOHandler) IsService() bool {
	return h.serviceService.IsServiceByProductId(h.req.GetProductId())
}

func (h *MOHandler) getService() (*entity.Service, error) {
	return h.serviceService.GetServiceByProductId(h.req.GetProductId())
}
