package cmd

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/idprm/go-xl-direct/internal/logger"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
)

var consumerMOCmd = &cobra.Command{
	Use:   "mo",
	Short: "Consumer MO Service CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		/**
		 * SETUP PGSQL
		 */
		db, err := connectPgsql()
		if err != nil {
			panic(err)
		}

		/**
		 * SETUP REDIS
		 */
		rds, err := connectRedis()
		if err != nil {
			panic(err)
		}

		/**
		 * SETUP RMQ
		 */
		rmq, err := connectRabbitMq()
		if err != nil {
			panic(err)
		}

		/**
		 * SETUP LOG
		 */
		logger := logger.NewLogger()

		/**
		 * SETUP CHANNEL
		 */
		rmq.SetUpChannel(RMQ_EXCHANGE_TYPE, true, RMQ_MO_EXCHANGE, true, RMQ_MO_QUEUE)
		rmq.SetUpChannel(RMQ_EXCHANGE_TYPE, true, RMQ_NOTIF_EXCHANGE, true, RMQ_NOTIF_QUEUE)
		rmq.SetUpChannel(RMQ_EXCHANGE_TYPE, true, RMQ_PB_MO_EXCHANGE, true, RMQ_PB_MO_QUEUE)

		messagesData, errSub := rmq.Subscribe(1, false, RMQ_MO_QUEUE, RMQ_MO_EXCHANGE, RMQ_MO_QUEUE)
		if errSub != nil {
			panic(errSub)
		}

		// Initial sync waiting group
		var wg sync.WaitGroup

		// Loop forever listening incoming data
		forever := make(chan bool)

		processor := NewProcessor(db, rds, rmq, logger)

		// Set into goroutine this listener
		go func() {

			// Loop every incoming data
			for d := range messagesData {

				wg.Add(1)
				processor.MO(&wg, d.Body)
				wg.Wait()

				// Manual consume queue
				d.Ack(false)

			}

		}()

		fmt.Println("[*] Waiting for data...")

		<-forever
	},
}

var consumerRenewalCmd = &cobra.Command{
	Use:   "renewal",
	Short: "Consumer Renewal Service CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		/**
		 * SETUP PGSQL
		 */
		db, err := connectPgsql()
		if err != nil {
			panic(err)
		}

		/**
		 * SETUP RMQ
		 */
		rmq, err := connectRabbitMq()
		if err != nil {
			panic(err)
		}

		/**
		 * SETUP LOG
		 */
		logger := logger.NewLogger()

		/**
		 * SETUP CHANNEL
		 */
		rmq.SetUpChannel(RMQ_EXCHANGE_TYPE, true, RMQ_RENEWAL_EXCHANGE, true, RMQ_RENEWAL_QUEUE)
		rmq.SetUpChannel(RMQ_EXCHANGE_TYPE, true, RMQ_NOTIF_EXCHANGE, true, RMQ_NOTIF_QUEUE)
		rmq.SetUpChannel(RMQ_EXCHANGE_TYPE, true, RMQ_PB_MT_EXCHANGE, true, RMQ_PB_MT_QUEUE)

		messagesData, errSub := rmq.Subscribe(1, false, RMQ_RENEWAL_QUEUE, RMQ_RENEWAL_EXCHANGE, RMQ_RENEWAL_QUEUE)
		if errSub != nil {
			panic(errSub)
		}
		// Initial sync waiting group
		var wg sync.WaitGroup

		// Loop forever listening incoming data
		forever := make(chan bool)

		// don't open redis connection if not needed
		processor := NewProcessor(db, &redis.Client{}, rmq, logger)

		// Set into goroutine this listener
		go func() {

			// Loop every incoming data
			for d := range messagesData {

				wg.Add(1)
				processor.Renewal(&wg, d.Body)
				wg.Wait()

				// Manual consume queue
				d.Ack(false)

			}

		}()

		fmt.Println("[*] Waiting for data...")

		<-forever
	},
}

var consumerRetryFpCmd = &cobra.Command{
	Use:   "retry_fp",
	Short: "Consumer Retry Firstpush Service CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		/**
		 * SETUP PGSQL
		 */
		db, err := connectPgsql()
		if err != nil {
			panic(err)
		}

		/**
		 * SETUP RMQ
		 */
		rmq, err := connectRabbitMq()
		if err != nil {
			panic(err)
		}

		/**
		 * SETUP LOG
		 */
		logger := logger.NewLogger()

		/**
		 * SETUP CHANNEL
		 */
		rmq.SetUpChannel(RMQ_EXCHANGE_TYPE, true, RMQ_RETRY_FP_EXCHANGE, true, RMQ_RETRY_FP_QUEUE)
		rmq.SetUpChannel(RMQ_EXCHANGE_TYPE, true, RMQ_NOTIF_EXCHANGE, true, RMQ_NOTIF_QUEUE)
		rmq.SetUpChannel(RMQ_EXCHANGE_TYPE, true, RMQ_PB_MT_EXCHANGE, true, RMQ_PB_MT_QUEUE)

		messagesData, errSub := rmq.Subscribe(1, false, RMQ_RETRY_FP_QUEUE, RMQ_RETRY_FP_EXCHANGE, RMQ_RETRY_FP_QUEUE)
		if errSub != nil {
			panic(errSub)
		}

		// Initial sync waiting group
		var wg sync.WaitGroup

		// Loop forever listening incoming data
		forever := make(chan bool)

		processor := NewProcessor(db, &redis.Client{}, rmq, logger)

		// Set into goroutine this listener
		go func() {

			// Loop every incoming data
			for d := range messagesData {

				wg.Add(1)
				processor.RetryFp(&wg, d.Body)
				wg.Wait()

				// Manual consume queue
				d.Ack(false)

			}

		}()

		fmt.Println("[*] Waiting for data...")

		<-forever
	},
}

var consumerRetryDpCmd = &cobra.Command{
	Use:   "retry_dp",
	Short: "Consumer Retry Dailypush Service CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		/**
		 * SETUP PGSQL
		 */
		db, err := connectPgsql()
		if err != nil {
			panic(err)
		}

		/**
		 * SETUP RMQ
		 */
		rmq, err := connectRabbitMq()
		if err != nil {
			panic(err)
		}

		/**
		 * SETUP LOG
		 */
		logger := logger.NewLogger()

		/**
		 * SETUP CHANNEL
		 */
		rmq.SetUpChannel(RMQ_EXCHANGE_TYPE, true, RMQ_RETRY_DP_EXCHANGE, true, RMQ_RETRY_DP_QUEUE)
		rmq.SetUpChannel(RMQ_EXCHANGE_TYPE, true, RMQ_NOTIF_EXCHANGE, true, RMQ_NOTIF_QUEUE)
		rmq.SetUpChannel(RMQ_EXCHANGE_TYPE, true, RMQ_PB_MT_EXCHANGE, true, RMQ_PB_MT_QUEUE)

		messagesData, errSub := rmq.Subscribe(1, false, RMQ_RETRY_DP_QUEUE, RMQ_RETRY_DP_EXCHANGE, RMQ_RETRY_DP_QUEUE)
		if errSub != nil {
			panic(errSub)
		}

		// Initial sync waiting group
		var wg sync.WaitGroup

		// Loop forever listening incoming data
		forever := make(chan bool)

		// don't open redis connection if not needed
		processor := NewProcessor(db, &redis.Client{}, rmq, logger)

		// Set into goroutine this listener
		go func() {

			// Loop every incoming data
			for d := range messagesData {

				wg.Add(1)
				processor.RetryDp(&wg, d.Body)
				wg.Wait()

				// Manual consume queue
				d.Ack(false)

			}

		}()

		fmt.Println("[*] Waiting for data...")

		<-forever
	},
}

var consumerRetryInsuffCmd = &cobra.Command{
	Use:   "retry_insuff",
	Short: "Consumer Retry Insuff Service CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		/**
		 * SETUP PGSQL
		 */
		db, err := connectPgsql()
		if err != nil {
			panic(err)
		}

		/**
		 * SETUP RMQ
		 */
		rmq, err := connectRabbitMq()
		if err != nil {
			panic(err)
		}

		/**
		 * SETUP LOG
		 */
		logger := logger.NewLogger()

		/**
		 * SETUP CHANNEL
		 */
		rmq.SetUpChannel(RMQ_EXCHANGE_TYPE, true, RMQ_RETRY_INSUFF_EXCHANGE, true, RMQ_RETRY_INSUFF_QUEUE)
		rmq.SetUpChannel(RMQ_EXCHANGE_TYPE, true, RMQ_NOTIF_EXCHANGE, true, RMQ_NOTIF_QUEUE)
		rmq.SetUpChannel(RMQ_EXCHANGE_TYPE, true, RMQ_PB_MT_EXCHANGE, true, RMQ_PB_MT_QUEUE)

		messagesData, errSub := rmq.Subscribe(1, false, RMQ_RETRY_INSUFF_QUEUE, RMQ_RETRY_INSUFF_EXCHANGE, RMQ_RETRY_INSUFF_QUEUE)
		if errSub != nil {
			panic(errSub)
		}

		// Initial sync waiting group
		var wg sync.WaitGroup

		// Loop forever listening incoming data
		forever := make(chan bool)

		// don't open redis connection if not needed
		processor := NewProcessor(db, &redis.Client{}, rmq, logger)

		// Set into goroutine this listener
		go func() {

			// Loop every incoming data
			for d := range messagesData {

				wg.Add(1)
				processor.RetryInsuff(&wg, d.Body)
				wg.Wait()

				// Manual consume queue
				d.Ack(false)

			}

		}()

		fmt.Println("[*] Waiting for data...")

		<-forever
	},
}

var consumerPostbackMOCmd = &cobra.Command{
	Use:   "postback_mo",
	Short: "Consumer Postback MO CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		/**
		 * SETUP RMQ
		 */
		rmq, err := connectRabbitMq()
		if err != nil {
			panic(err)
		}

		/**
		 * SETUP LOG
		 */
		logger := logger.NewLogger()

		/**
		 * SETUP CHANNEL
		 */
		rmq.SetUpChannel(RMQ_EXCHANGE_TYPE, true, RMQ_PB_MO_EXCHANGE, true, RMQ_PB_MO_QUEUE)

		messagesData, errSub := rmq.Subscribe(1, false, RMQ_PB_MO_QUEUE, RMQ_PB_MO_EXCHANGE, RMQ_PB_MO_QUEUE)
		if errSub != nil {
			panic(errSub)
		}

		// Initial sync waiting group
		var wg sync.WaitGroup

		// Loop forever listening incoming data
		forever := make(chan bool)

		// don't open db connection if not needed
		processor := NewProcessor(&sql.DB{}, &redis.Client{}, rmq, logger)

		// Set into goroutine this listener
		go func() {

			// Loop every incoming data
			for d := range messagesData {

				wg.Add(1)
				processor.PostbackMO(&wg, d.Body)
				wg.Wait()

				// Manual consume queue
				d.Ack(false)
			}

		}()

		fmt.Println("[*] Waiting for data...")

		<-forever
	},
}

var consumerPostbackMTCmd = &cobra.Command{
	Use:   "postback_mt",
	Short: "Consumer Postback MT CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		/**
		 * SETUP RMQ
		 */
		rmq, err := connectRabbitMq()
		if err != nil {
			panic(err)
		}

		/**
		 * SETUP LOG
		 */
		logger := logger.NewLogger()

		/**
		 * SETUP CHANNEL
		 */
		rmq.SetUpChannel(RMQ_EXCHANGE_TYPE, true, RMQ_PB_MT_EXCHANGE, true, RMQ_PB_MT_QUEUE)

		messagesData, errSub := rmq.Subscribe(1, false, RMQ_PB_MT_QUEUE, RMQ_PB_MT_EXCHANGE, RMQ_PB_MT_QUEUE)
		if errSub != nil {
			panic(errSub)
		}

		// Initial sync waiting group
		var wg sync.WaitGroup

		// Loop forever listening incoming data
		forever := make(chan bool)

		// don't open db connection if not needed
		processor := NewProcessor(&sql.DB{}, &redis.Client{}, rmq, logger)

		// Set into goroutine this listener
		go func() {

			// Loop every incoming data
			for d := range messagesData {

				wg.Add(1)
				processor.PostbackMT(&wg, d.Body)
				wg.Wait()

				// Manual consume queue
				d.Ack(false)
			}

		}()

		fmt.Println("[*] Waiting for data...")

		<-forever
	},
}

var consumerNotifCmd = &cobra.Command{
	Use:   "notif",
	Short: "Consumer Notif Portal CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		/**
		 * SETUP RMQ
		 */
		rmq, err := connectRabbitMq()
		if err != nil {
			panic(err)
		}

		/**
		 * SETUP LOG
		 */
		logger := logger.NewLogger()

		/**
		 * SETUP CHANNEL
		 */
		rmq.SetUpChannel(RMQ_EXCHANGE_TYPE, true, RMQ_NOTIF_EXCHANGE, true, RMQ_NOTIF_QUEUE)

		messagesData, errSub := rmq.Subscribe(1, false, RMQ_NOTIF_QUEUE, RMQ_NOTIF_EXCHANGE, RMQ_NOTIF_QUEUE)
		if errSub != nil {
			panic(errSub)
		}

		// Initial sync waiting group
		var wg sync.WaitGroup

		// Loop forever listening incoming data
		forever := make(chan bool)

		// don't open db and redis connection if not needed
		processor := NewProcessor(&sql.DB{}, &redis.Client{}, rmq, logger)

		// Set into goroutine this listener
		go func() {

			// Loop every incoming data
			for d := range messagesData {

				wg.Add(1)
				processor.Notif(&wg, d.Body)
				wg.Wait()

				// Manual consume queue
				d.Ack(false)
			}

		}()

		fmt.Println("[*] Waiting for data...")

		<-forever
	},
}

var consumerTrafficCmd = &cobra.Command{
	Use:   "traffic",
	Short: "Consumer Traffic Service CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		/**
		 * connect pgsql
		 */
		db, err := connectPgsql()
		if err != nil {
			panic(err)
		}

		/**
		 * connect rabbitmq
		 */
		rmq, err := connectRabbitMq()
		if err != nil {
			panic(err)
		}

		/**
		 * SETUP LOG
		 */
		logger := logger.NewLogger()

		/**
		 * SETUP CHANNEL
		 */
		rmq.SetUpChannel(RMQ_EXCHANGE_TYPE, true, RMQ_TRAFFIC_EXCHANGE, true, RMQ_TRAFFIC_QUEUE)

		messagesData, errSub := rmq.Subscribe(1, false, RMQ_TRAFFIC_QUEUE, RMQ_TRAFFIC_EXCHANGE, RMQ_TRAFFIC_QUEUE)
		if errSub != nil {
			panic(errSub)
		}

		// Initial sync waiting group
		var wg sync.WaitGroup

		// Loop forever listening incoming data
		forever := make(chan bool)

		processor := NewProcessor(db, &redis.Client{}, rmq, logger)

		// Set into goroutine this listener
		go func() {

			// Loop every incoming data
			for d := range messagesData {

				wg.Add(1)
				processor.Traffic(&wg, d.Body)
				wg.Wait()

				// Manual consume queue
				d.Ack(false)
			}

		}()

		fmt.Println("[*] Waiting for data...")

		<-forever
	},
}

var consumerDailypushCmd = &cobra.Command{
	Use:   "dailypush",
	Short: "Consumer Dailypush Service CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		/**
		 * connect pgsql
		 */
		db, err := connectPgsql()
		if err != nil {
			panic(err)
		}

		/**
		 * connect rabbitmq
		 */
		rmq, err := connectRabbitMq()
		if err != nil {
			panic(err)
		}
		/**
		 * SETUP LOG
		 */
		logger := logger.NewLogger()

		/**
		 * SETUP CHANNEL
		 */
		rmq.SetUpChannel(RMQ_EXCHANGE_TYPE, true, RMQ_DAILYPUSH_EXCHANGE, true, RMQ_DAILYPUSH_QUEUE)

		messagesData, errSub := rmq.Subscribe(1, false, RMQ_DAILYPUSH_QUEUE, RMQ_DAILYPUSH_EXCHANGE, RMQ_DAILYPUSH_QUEUE)
		if errSub != nil {
			panic(errSub)
		}

		// Initial sync waiting group
		var wg sync.WaitGroup

		// Loop forever listening incoming data
		forever := make(chan bool)

		processor := NewProcessor(db, &redis.Client{}, rmq, logger)

		// Set into goroutine this listener
		go func() {

			// Loop every incoming data
			for d := range messagesData {

				wg.Add(1)
				processor.Dailypush(&wg, d.Body)
				wg.Wait()

				// Manual consume queue
				d.Ack(false)
			}

		}()

		fmt.Println("[*] Waiting for data...")

		<-forever

	},
}
