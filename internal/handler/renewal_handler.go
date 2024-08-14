package handler

import (
	"encoding/json"
	"log"
	"time"

	"github.com/idprm/go-xl-direct/internal/domain/entity"
	"github.com/idprm/go-xl-direct/internal/domain/model"
	"github.com/idprm/go-xl-direct/internal/logger"
	"github.com/idprm/go-xl-direct/internal/services"
	"github.com/wiliehidayat87/rmqp"
)

type RenewalHandler struct {
	rmq                 rmqp.AMQP
	logger              *logger.Logger
	sub                 *entity.Subscription
	serviceService      services.IServiceService
	subscriptionService services.ISubscriptionService
	transactionService  services.ITransactionService
	req                 *model.NotificationRequest
}

func NewRenewalHandler(
	rmq rmqp.AMQP,
	logger *logger.Logger,
	sub *entity.Subscription,
	serviceService services.IServiceService,
	subscriptionService services.ISubscriptionService,
	transactionService services.ITransactionService,
	req *model.NotificationRequest,
) *RenewalHandler {
	return &RenewalHandler{
		rmq:                 rmq,
		logger:              logger,
		sub:                 sub,
		serviceService:      serviceService,
		subscriptionService: subscriptionService,
		transactionService:  transactionService,
		req:                 req,
	}
}

func (h *RenewalHandler) Dailypush() {
	// check if active sub
	if h.subscriptionService.IsActiveSubscription(h.sub.GetServiceId(), h.sub.GetMsisdn()) {

		service, err := h.serviceService.GetServiceByCode(h.req.GetProductId())
		if err != nil {
			log.Println(err.Error())
		}

		if h.req.IsActive() {
			h.subscriptionService.UpdateSuccess(
				&entity.Subscription{
					ServiceID:          h.sub.GetServiceId(),
					Msisdn:             h.sub.GetMsisdn(),
					LatestTrxId:        "",
					LatestSubject:      "",
					LatestPIN:          "",
					Amount:             service.GetPrice(),
					RenewalAt:          time.Now(),
					ChargeAt:           time.Now(),
					Success:            1,
					IsRetry:            false,
					TotalRenewal:       1,
					TotalAmountRenewal: service.GetPrice(),
					LatestPayload:      "-",
				},
			)

			h.transactionService.SaveTransaction(
				&entity.Transaction{
					ServiceID:      h.sub.GetServiceId(),
					Msisdn:         h.sub.GetMsisdn(),
					Channel:        "",
					Adnet:          "",
					Keyword:        "",
					Amount:         service.GetPrice(),
					PIN:            "-",
					Status:         "",
					StatusCode:     "",
					StatusDetail:   "",
					Subject:        "",
					Payload:        "",
					CampKeyword:    "",
					CampSubKeyword: "",
					IpAddress:      "",
				},
			)

			// insert to rabbitmq
			jsonData, _ := json.Marshal(
				&entity.NotifParamsRequest{
					Service:      service,
					Subscription: h.sub,
					Action:       "RENEWAL",
					Pin:          "",
				},
			)
			h.rmq.IntegratePublish(
				RMQ_NOTIF_EXCHANGE,
				RMQ_NOTIF_QUEUE,
				RMQ_DATA_TYPE,
				"", string(jsonData),
			)
		}

		if h.req.IsCancelled() {

			h.subscriptionService.UpdateFailed(
				&entity.Subscription{
					ServiceID:     h.sub.GetServiceId(),
					Msisdn:        h.sub.GetMsisdn(),
					LatestTrxId:   "",
					LatestSubject: "",
					LatestStatus:  "",
					RenewalAt:     time.Now(),
					RetryAt:       time.Now(),
					Failed:        1,
					IsRetry:       true,
					LatestPayload: "",
				},
			)

			h.transactionService.SaveTransaction(
				&entity.Transaction{
					ServiceID:      h.sub.GetServiceId(),
					Msisdn:         h.sub.GetMsisdn(),
					Channel:        "",
					Adnet:          "",
					Keyword:        "",
					Status:         "",
					StatusCode:     "",
					StatusDetail:   "",
					Subject:        "",
					Payload:        "",
					CampKeyword:    "",
					CampSubKeyword: "",
					IpAddress:      "",
				},
			)
		}

		jsonDataPB, _ := json.Marshal(
			&entity.PostbackParamsRequest{
				Subscription: &entity.Subscription{
					LatestTrxId:    "",
					ServiceID:      h.sub.GetServiceId(),
					Msisdn:         h.sub.GetMsisdn(),
					LatestKeyword:  h.sub.GetLatestKeyword(),
					LatestSubject:  SUBJECT_RENEWAL,
					LatestPayload:  "",
					CampKeyword:    h.sub.GetCampKeyword(),
					CampSubKeyword: h.sub.GetCampSubKeyword(),
				},
				Service:   service,
				Action:    "MT_DAILYPUSH",
				Status:    "",
				AffSub:    h.sub.GetAffSub(),
				IsSuccess: false,
			},
		)
		h.rmq.IntegratePublish(
			RMQ_PB_EXCHANGE,
			RMQ_PB_QUEUE,
			RMQ_DATA_TYPE, "", string(jsonDataPB),
		)

		jsonDataDP, _ := json.Marshal(
			&entity.DailypushBodyRequest{
				TxId:           "",
				SubscriptionId: h.sub.GetId(),
				ServiceId:      h.sub.GetServiceId(),
				Msisdn:         h.sub.GetMsisdn(),
				Channel:        h.sub.GetChannel(),
				CampKeyword:    h.sub.GetCampKeyword(),
				CampSubKeyword: h.sub.GetCampSubKeyword(),
				Adnet:          h.sub.GetAdnet(),
				PubID:          h.sub.GetPubId(),
				AffSub:         h.sub.GetPubId(),
				Subject:        "",
				StatusCode:     "",
				StatusDetail:   "",
				IsCharge:       false,
				IpAddress:      h.sub.GetIpAddress(),
				Action:         SUBJECT_RENEWAL,
			},
		)

		h.rmq.IntegratePublish(
			RMQ_DAILYPUSH_EXCHANGE,
			RMQ_DAILYPUSH_QUEUE,
			RMQ_DATA_TYPE,
			"", string(jsonDataDP),
		)
	}

}
