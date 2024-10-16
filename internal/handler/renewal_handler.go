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

type RenewalHandler struct {
	rmq                 rmqp.AMQP
	logger              *logger.Logger
	serviceService      services.IServiceService
	subscriptionService services.ISubscriptionService
	transactionService  services.ITransactionService
	req                 *model.NotificationRequest
}

func NewRenewalHandler(
	rmq rmqp.AMQP,
	logger *logger.Logger,
	serviceService services.IServiceService,
	subscriptionService services.ISubscriptionService,
	transactionService services.ITransactionService,
	req *model.NotificationRequest,
) *RenewalHandler {
	return &RenewalHandler{
		rmq:                 rmq,
		logger:              logger,
		serviceService:      serviceService,
		subscriptionService: subscriptionService,
		transactionService:  transactionService,
		req:                 req,
	}
}

func (h *RenewalHandler) Dailypush() {

	log.Println("in_renewal_handler")
	log.Println(h.req)

	trxId := utils.GenerateTrxId()

	service, err := h.serviceService.GetServiceByProductId(h.req.ProductId)
	if err != nil {
		log.Println(err.Error())
	}

	// check if active sub
	if h.subscriptionService.IsActiveSubscription(service.GetId(), h.req.GetUserIdentifier()) {

		service, err := h.serviceService.GetServiceByProductId(h.req.GetProductId())
		if err != nil {
			log.Println(err.Error())
		}

		sub, err := h.subscriptionService.SelectSubscription(service.GetId(), h.req.GetUserIdentifier())
		if err != nil {
			log.Println(err)
		}

		if h.req.IsActive() {
			subSuccess := &entity.Subscription{
				ServiceID:          service.GetId(),
				Msisdn:             h.req.GetUserIdentifier(),
				LatestTrxId:        h.req.GetTransactionId(),
				LatestSubject:      SUBJECT_RENEWAL,
				LatestStatus:       STATUS_SUCCESS,
				LatestPIN:          "",
				Amount:             h.req.GetAmount(),
				RenewalAt:          h.req.GetNextRenewalDate(),
				ChargeAt:           time.Now(),
				Success:            1,
				IsRetry:            false,
				TotalRenewal:       1,
				TotalAmountRenewal: h.req.GetAmount(),
				LatestPayload:      "-",
			}

			h.subscriptionService.UpdateSuccess(subSuccess)

			transSuccess := &entity.Transaction{
				TxID:           h.req.GetTransactionId(),
				ServiceID:      service.GetId(),
				Msisdn:         h.req.GetUserIdentifier(),
				SubID:          h.req.GetSubscriptionId(),
				Channel:        sub.GetChannel(),
				Adnet:          sub.GetAdnet(),
				Keyword:        MO_REG + " " + service.GetCode(),
				Amount:         h.req.GetAmount(),
				PIN:            "",
				Status:         STATUS_SUCCESS,
				StatusCode:     "",
				StatusDetail:   "",
				Subject:        SUBJECT_RENEWAL,
				Payload:        "",
				CampKeyword:    sub.GetCampKeyword(),
				CampSubKeyword: sub.GetCampSubKeyword(),
				IpAddress:      sub.GetIpAddress(),
			}

			h.transactionService.SaveTransaction(transSuccess)

			// insert to rabbitmq
			jsonData, _ := json.Marshal(
				&entity.NotifParamsRequest{
					Service:      service,
					Subscription: subSuccess,
					Action:       "RENEWAL",
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
		}

		if h.req.IsCancelled() {

			subFailed := &entity.Subscription{
				ServiceID:     service.GetId(),
				Msisdn:        h.req.GetUserIdentifier(),
				LatestTrxId:   trxId,
				LatestSubject: SUBJECT_RENEWAL,
				LatestStatus:  STATUS_FAILED,
				RenewalAt:     time.Now().AddDate(0, 0, 1),
				RetryAt:       time.Now(),
				Failed:        1,
				IsRetry:       true,
				LatestPayload: "",
			}
			h.subscriptionService.UpdateFailed(subFailed)

			transFailed := &entity.Transaction{
				TxID:           trxId,
				ServiceID:      service.GetId(),
				Msisdn:         h.req.GetUserIdentifier(),
				SubID:          h.req.GetSubscriptionId(),
				Channel:        sub.GetChannel(),
				Adnet:          sub.GetAdnet(),
				Keyword:        MO_REG + " " + service.GetCode(),
				Status:         STATUS_FAILED,
				StatusCode:     "",
				StatusDetail:   "",
				Subject:        SUBJECT_RENEWAL,
				Payload:        "",
				CampKeyword:    sub.GetCampKeyword(),
				CampSubKeyword: sub.GetCampSubKeyword(),
				IpAddress:      sub.GetIpAddress(),
			}
			h.transactionService.SaveTransaction(transFailed)
		}

		jsonDataPB, _ := json.Marshal(
			&entity.PostbackParamsRequest{
				Subscription: &entity.Subscription{
					LatestTrxId:    h.req.GetTransactionId(),
					ServiceID:      service.GetId(),
					Msisdn:         h.req.GetUserIdentifier(),
					LatestKeyword:  "",
					LatestSubject:  SUBJECT_RENEWAL,
					LatestPayload:  "",
					CampKeyword:    "",
					CampSubKeyword: "",
				},
				Service:   service,
				Action:    "MT_DAILYPUSH",
				Status:    "",
				AffSub:    "",
				IsSuccess: false,
			},
		)
		h.rmq.IntegratePublish(
			RMQ_PB_MT_EXCHANGE,
			RMQ_PB_MT_QUEUE,
			RMQ_DATA_TYPE, "", string(jsonDataPB),
		)
	}
}
