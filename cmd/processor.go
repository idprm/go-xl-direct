package cmd

import (
	"database/sql"
	"sync"

	"github.com/idprm/go-xl-direct/internal/logger"
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

	wg.Done()
}

func (p *Processor) Renewal(wg *sync.WaitGroup, message []byte) {

	wg.Done()
}

func (p *Processor) RetryFp(wg *sync.WaitGroup, message []byte) {

	wg.Done()
}

func (p *Processor) RetryDp(wg *sync.WaitGroup, message []byte) {

	wg.Done()
}

func (p *Processor) RetryInsuff(wg *sync.WaitGroup, message []byte) {

	wg.Done()
}

func (p *Processor) PostbackMO(wg *sync.WaitGroup, message []byte) {

	wg.Done()
}

func (p *Processor) PostbackMT(wg *sync.WaitGroup, message []byte) {

	wg.Done()
}

func (p *Processor) Notif(wg *sync.WaitGroup, message []byte) {

	wg.Done()
}

func (p *Processor) Traffic(wg *sync.WaitGroup, message []byte) {

	wg.Done()
}

func (p *Processor) Dailypush(wg *sync.WaitGroup, message []byte) {

	wg.Done()
}
