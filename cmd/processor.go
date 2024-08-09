package cmd

import (
	"database/sql"
	"encoding/json"
	"sync"

	"github.com/idprm/go-xl-direct/internal/domain/entity"
	"github.com/idprm/go-xl-direct/internal/domain/model"
	"github.com/idprm/go-xl-direct/internal/domain/repository"
	"github.com/idprm/go-xl-direct/internal/handler"
	"github.com/idprm/go-xl-direct/internal/logger"
	"github.com/idprm/go-xl-direct/internal/services"
	"github.com/redis/go-redis/v9"
	"github.com/wiliehidayat87/rmqp"
)

type Processor struct {
	db     *sql.DB
	rds    *redis.Client
	rmq    rmqp.AMQP
	logger *logger.Logger
}

func NewProcessor(
	db *sql.DB,
	rds *redis.Client,
	rmq rmqp.AMQP,
	logger *logger.Logger,
) *Processor {
	return &Processor{
		db:     db,
		rds:    rds,
		rmq:    rmq,
		logger: logger,
	}
}

func (p *Processor) MO(wg *sync.WaitGroup, message []byte) {
	/**
	 * -. Filter REG / UNREG
	 * -. Check Blacklist
	 * -. Check Active Sub
	 * -. MT API
	 * -. Save Sub
	 * -/ Save Transaction
	 */

	blacklistRepo := repository.NewBlacklistRepository(p.db)
	blacklistService := services.NewBlacklistService(blacklistRepo)
	serviceRepo := repository.NewServiceRepository(p.db)
	serviceService := services.NewServiceService(serviceRepo)
	verifyRepo := repository.NewVerifyRepository(p.rds)
	verifyService := services.NewVerifyService(verifyRepo)
	subscriptionRepo := repository.NewSubscriptionRepository(p.db)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)
	transactionRepo := repository.NewTransactionRepository(p.db)
	transactionService := services.NewTransactionService(transactionRepo)
	historyRepo := repository.NewHistoryRepository(p.db)
	historyService := services.NewHistoryService(historyRepo)
	trafficRepo := repository.NewTrafficRepository(p.db)
	trafficService := services.NewTrafficService(trafficRepo)

	var req *model.NotificationRequest
	json.Unmarshal([]byte(message), &req)

	h := handler.NewMOHandler(
		p.rmq,
		p.logger,
		blacklistService,
		serviceService,
		verifyService,
		subscriptionService,
		transactionService,
		historyService,
		trafficService,
		req,
	)

	if req.IsSubscription() {

		if req.IsActive() {
			h.Firstpush()
		}

		if req.IsCancelled() {
			h.Unsub()
		}

	}

	if req.IsRenewal() {
		h.Renewal()
	}

	if req.IsRefund() {
		h.Refund()
	}

	wg.Done()
}

func (p *Processor) Renewal(wg *sync.WaitGroup, message []byte) {

	serviceRepo := repository.NewServiceRepository(p.db)
	serviceService := services.NewServiceService(serviceRepo)
	subscriptionRepo := repository.NewSubscriptionRepository(p.db)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)
	transactionRepo := repository.NewTransactionRepository(p.db)
	transactionService := services.NewTransactionService(transactionRepo)

	var req *model.NotificationRequest
	json.Unmarshal([]byte(message), &req)

	handler.NewRenewalHandler(
		p.rmq,
		p.logger,
		&entity.Subscription{},
		serviceService,
		subscriptionService,
		transactionService,
		req,
	)

	wg.Done()
}

func (p *Processor) PostbackMO(wg *sync.WaitGroup, message []byte) {

	var req *entity.PostbackParamsRequest
	json.Unmarshal(message, &req)

	h := handler.NewPostbackHandler(p.logger, req)

	if req.IsMO() {
		h.Postback()
	}

	wg.Done()
}

func (p *Processor) PostbackMT(wg *sync.WaitGroup, message []byte) {
	var req *entity.PostbackParamsRequest
	json.Unmarshal(message, &req)

	if req.IsMTDailypush() {

	}

	if req.IsMTFirstpush() {

	}
	wg.Done()
}

func (p *Processor) Notif(wg *sync.WaitGroup, message []byte) {

	wg.Done()
}

func (p *Processor) Traffic(wg *sync.WaitGroup, message []byte) {
	/**
	 * load repo
	 */
	trafficRepo := repository.NewTrafficRepository(p.db)
	trafficService := services.NewTrafficService(trafficRepo)

	var req *entity.ReqTrafficParams
	json.Unmarshal(message, &req)

	h := handler.NewTrafficHandler(trafficService, req)

	h.Campaign()

	wg.Done()
}

func (p *Processor) Dailypush(wg *sync.WaitGroup, message []byte) {
	/**
	 * load repo
	 */
	dailypushRepo := repository.NewDailypushRepository(p.db)
	dailypushService := services.NewDailypushService(dailypushRepo)

	var req *entity.DailypushBodyRequest
	json.Unmarshal(message, &req)

	h := handler.NewDailypushHandler(dailypushService, req)

	h.Dailypush()

	wg.Done()
}
