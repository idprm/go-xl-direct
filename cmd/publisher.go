package cmd

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/idprm/go-xl-direct/internal/domain/entity"
	"github.com/idprm/go-xl-direct/internal/domain/repository"
	"github.com/idprm/go-xl-direct/internal/providers/rabbit"
	"github.com/idprm/go-xl-direct/internal/services"
	"github.com/spf13/cobra"
	"github.com/wiliehidayat87/rmqp"
)

var publisherRenewalCmd = &cobra.Command{
	Use:   "pub_renewal",
	Short: "Publisher Renewal Service CLI",
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
		 * SETUP CHANNEL
		 */
		rmq.SetUpChannel(
			RMQ_EXCHANGE_TYPE,
			true,
			RMQ_RENEWAL_EXCHANGE,
			true,
			RMQ_RENEWAL_QUEUE,
		)

		/**
		 * Looping schedule
		 */
		timeDuration := time.Duration(1)

		for {
			timeNow := time.Now().Format("15:04")

			scheduleRepo := repository.NewScheduleRepository(db)
			scheduleService := services.NewScheduleService(scheduleRepo)

			if scheduleService.GetUnlocked(ACT_RENEWAL, timeNow) {

				scheduleService.UpdateSchedule(false, ACT_RENEWAL)

				go func() {
					populateRenewal(db, rmq)
				}()
			}

			if scheduleService.GetLocked(ACT_RENEWAL, timeNow) {
				scheduleService.UpdateSchedule(true, ACT_RENEWAL)
			}

			time.Sleep(timeDuration * time.Minute)

		}
	},
}

var publisherRetryDpCmd = &cobra.Command{
	Use:   "pub_retry_dp",
	Short: "Publisher Retry Dailypush Service CLI",
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
		 * SETUP CHANNEL
		 */
		rmq.SetUpChannel(
			RMQ_EXCHANGE_TYPE,
			true,
			RMQ_RETRY_DP_EXCHANGE,
			true,
			RMQ_RETRY_DP_QUEUE,
		)

		/**
		 * Looping schedule
		 */
		timeDuration := time.Duration(1)

		for {

			/**
			** Populate retry if queue message is zero or 0
			**/
			p := rabbit.NewRabbitMQ()

			q, err := p.Queue(RMQ_RETRY_DP_QUEUE)
			if err != nil {
				log.Println(err)
			}

			var res *entity.RabbitMQResponse
			json.Unmarshal(q, &res)

			// if queue is empty
			if !res.IsRunning() {
				go func() {
					populateRetryDp(db, rmq)
				}()
			}

			time.Sleep(timeDuration * time.Minute)

		}
	},
}

var publisherRetryFpCmd = &cobra.Command{
	Use:   "pub_retry_fp",
	Short: "Publisher Retry Firstpush Service CLI",
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
		 * SETUP CHANNEL
		 */
		rmq.SetUpChannel(
			RMQ_EXCHANGE_TYPE,
			true,
			RMQ_RETRY_FP_EXCHANGE,
			true,
			RMQ_RETRY_FP_QUEUE,
		)

		/**
		 * Looping schedule
		 */
		timeDuration := time.Duration(1)

		for {

			/**
			** Populate retry if queue message is zero or 0
			**/
			p := rabbit.NewRabbitMQ()

			q, err := p.Queue(RMQ_RETRY_FP_QUEUE)
			if err != nil {
				log.Println(err)
			}

			var res *entity.RabbitMQResponse
			json.Unmarshal(q, &res)

			// if queue is empty
			if !res.IsRunning() {
				go func() {
					populateRetryFp(db, rmq)
				}()
			}

			time.Sleep(timeDuration * time.Minute)

		}
	},
}

var publisherRetryInsuffCmd = &cobra.Command{
	Use:   "pub_retry_insuff",
	Short: "Publisher Retry Insuff Service CLI",
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
		 * SETUP CHANNEL
		 */
		rmq.SetUpChannel(
			RMQ_EXCHANGE_TYPE,
			true,
			RMQ_RETRY_INSUFF_EXCHANGE,
			true,
			RMQ_RETRY_INSUFF_QUEUE,
		)

		/**
		 * Looping schedule
		 */
		timeDuration := time.Duration(1)

		for {
			timeNow := time.Now().Format("15:04")

			scheduleRepo := repository.NewScheduleRepository(db)
			scheduleService := services.NewScheduleService(scheduleRepo)

			if scheduleService.GetUnlocked(ACT_RETRY_INSUFF, timeNow) {

				scheduleService.UpdateSchedule(false, ACT_RETRY_INSUFF)

				go func() {
					populateRetryInsuff(db, rmq)
				}()
			}

			if scheduleService.GetLocked(ACT_RETRY_INSUFF, timeNow) {
				scheduleService.UpdateSchedule(true, ACT_RETRY_INSUFF)
			}

			time.Sleep(timeDuration * time.Minute)

		}
	},
}

var publisherCSVCmd = &cobra.Command{
	Use:   "pub_csv",
	Short: "Publisher CSV Service CLI",
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
		 * Looping schedule
		 */
		timeDuration := time.Duration(1)

		for {
			timeNow := time.Now().Format("15:04")

			scheduleRepo := repository.NewScheduleRepository(db)
			scheduleService := services.NewScheduleService(scheduleRepo)

			if scheduleService.GetUnlocked(ACT_CSV, timeNow) {

				scheduleService.UpdateSchedule(false, ACT_CSV)

				go func() {
					populateCSV(db)
				}()
			}

			if scheduleService.GetLocked(ACT_CSV, timeNow) {
				scheduleService.UpdateSchedule(true, ACT_CSV)
			}

			time.Sleep(timeDuration * time.Minute)
		}
	},
}

var publisherUploadCSVCmd = &cobra.Command{
	Use:   "pub_upload_csv",
	Short: "Upload CSV CLI",
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
		 * Looping schedule
		 */
		timeDuration := time.Duration(1)

		for {
			timeNow := time.Now().Format("15:04")

			scheduleRepo := repository.NewScheduleRepository(db)
			scheduleService := services.NewScheduleService(scheduleRepo)

			if scheduleService.GetUnlocked(ACT_UPLOAD_CSV, timeNow) {

				scheduleService.UpdateSchedule(false, ACT_UPLOAD_CSV)

				go func() {
					//
				}()
			}

			if scheduleService.GetLocked(ACT_UPLOAD_CSV, timeNow) {
				scheduleService.UpdateSchedule(true, ACT_UPLOAD_CSV)
			}

			time.Sleep(timeDuration * time.Minute)
		}

	},
}

func populateRenewal(db *sql.DB, rmq rmqp.AMQP) {
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)

	subs := subscriptionService.RenewalSubscription()
	for _, s := range *subs {
		var sub entity.Subscription

		sub.ID = s.ID
		sub.ServiceID = s.ServiceID
		sub.Msisdn = s.Msisdn
		sub.Channel = s.Channel
		sub.Adnet = s.Adnet
		sub.LatestKeyword = s.LatestKeyword
		sub.LatestSubject = s.LatestSubject
		sub.LatestPIN = s.LatestPIN
		sub.IpAddress = s.IpAddress
		sub.AffSub = s.AffSub
		sub.CampKeyword = s.CampKeyword
		sub.CampSubKeyword = s.CampSubKeyword
		sub.ContentSequence = s.ContentSequence
		sub.CreatedAt = s.CreatedAt

		json, err := json.Marshal(sub)
		if err != nil {
			log.Println(err)
		}

		pub := rmq.IntegratePublish(
			RMQ_RENEWAL_EXCHANGE,
			RMQ_RENEWAL_QUEUE,
			RMQ_DATA_TYPE,
			"",
			string(json),
		)

		if !pub {
			log.Println(json)
		}

		time.Sleep(100 * time.Microsecond)
	}
}

func populateRetryFp(db *sql.DB, rmq rmqp.AMQP) {
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)

	subs := subscriptionService.RetryFpSubscription()

	for _, s := range *subs {
		var sub entity.Subscription

		sub.ID = s.ID
		sub.ServiceID = s.ServiceID
		sub.Msisdn = s.Msisdn
		sub.Channel = s.Channel
		sub.Adnet = s.Adnet
		sub.LatestKeyword = s.LatestKeyword
		sub.LatestSubject = s.LatestSubject
		sub.LatestPIN = s.LatestPIN
		sub.IpAddress = s.IpAddress
		sub.AffSub = s.AffSub
		sub.CampKeyword = s.CampKeyword
		sub.CampSubKeyword = s.CampSubKeyword
		sub.ContentSequence = s.ContentSequence
		sub.RetryAt = s.RetryAt
		sub.CreatedAt = s.CreatedAt

		json, err := json.Marshal(sub)
		if err != nil {
			log.Println(err)
		}

		pub := rmq.IntegratePublish(
			RMQ_RETRY_FP_EXCHANGE,
			RMQ_RETRY_FP_QUEUE,
			RMQ_DATA_TYPE,
			"",
			string(json),
		)

		if !pub {
			log.Println(json)
		}

		time.Sleep(100 * time.Microsecond)
	}
}

func populateRetryDp(db *sql.DB, rmq rmqp.AMQP) {
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)

	subs := subscriptionService.RetryDpSubscription()

	for _, s := range *subs {
		var sub entity.Subscription

		sub.ID = s.ID
		sub.ServiceID = s.ServiceID
		sub.Msisdn = s.Msisdn
		sub.Channel = s.Channel
		sub.Adnet = s.Adnet
		sub.LatestKeyword = s.LatestKeyword
		sub.LatestSubject = s.LatestSubject
		sub.LatestPIN = s.LatestPIN
		sub.IpAddress = s.IpAddress
		sub.AffSub = s.AffSub
		sub.CampKeyword = s.CampKeyword
		sub.CampSubKeyword = s.CampSubKeyword
		sub.ContentSequence = s.ContentSequence
		sub.RetryAt = s.RetryAt
		sub.CreatedAt = s.CreatedAt

		json, err := json.Marshal(sub)
		if err != nil {
			log.Println(err)
		}

		pub := rmq.IntegratePublish(
			RMQ_RETRY_DP_EXCHANGE,
			RMQ_RETRY_DP_QUEUE,
			RMQ_DATA_TYPE,
			"",
			string(json),
		)

		if !pub {
			log.Println(json)
		}

		time.Sleep(100 * time.Microsecond)
	}
}

func populateRetryInsuff(db *sql.DB, rmq rmqp.AMQP) {
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)

	subs := subscriptionService.RetryInsuffSubscription()

	for _, s := range *subs {
		var sub entity.Subscription

		sub.ID = s.ID
		sub.ServiceID = s.ServiceID
		sub.Msisdn = s.Msisdn
		sub.Channel = s.Channel
		sub.Adnet = s.Adnet
		sub.LatestKeyword = s.LatestKeyword
		sub.LatestSubject = s.LatestSubject
		sub.LatestPIN = s.LatestPIN
		sub.IpAddress = s.IpAddress
		sub.AffSub = s.AffSub
		sub.CampKeyword = s.CampKeyword
		sub.CampSubKeyword = s.CampSubKeyword
		sub.ContentSequence = s.ContentSequence
		sub.RetryAt = s.RetryAt
		sub.CreatedAt = s.CreatedAt

		json, err := json.Marshal(sub)
		if err != nil {
			log.Println(err)
		}

		pub := rmq.IntegratePublish(
			RMQ_RETRY_INSUFF_EXCHANGE,
			RMQ_RETRY_INSUFF_QUEUE,
			RMQ_DATA_TYPE,
			"",
			string(json),
		)

		if !pub {
			log.Println(json)
		}

		time.Sleep(100 * time.Microsecond)
	}
}

func populateCSV(db *sql.DB) {

	fileSubs := LOG_PATH + "/csv/subscriptions_id_telkomsel_" + APP_NAME + ".csv"
	fileTrans := LOG_PATH + "/csv/transactions_id_telkomsel_" + APP_NAME + ".csv"

	subscriptionRepo := repository.NewSubscriptionRepository(db)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)
	transactionRepo := repository.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)

	// delete file
	os.Remove(fileSubs)

	subCsv, err := os.Create(fileSubs)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	defer subCsv.Close()
	subW := csv.NewWriter(subCsv)
	defer subW.Flush()

	subsHeaders := []string{
		"country", "operator", "service", "source", "msisdn",
		"status", "cycle", "adnet", "revenue", "subs_date",
		"renewal_date", "freemium_end_date", "unsubs_from", "unsubs_date",
		"service_price", "currency", "profile_status", "publisher",
		"trxid", "pixel", "handset", "browser", "attempt_charging",
		"success_billing",
	}
	subW.Write(subsHeaders)

	subRecords, err := subscriptionService.SelectSubcriptionToCSV()
	if err != nil {
		log.Fatalf("error load table subscriptions: %s", err)
	}

	var subsData [][]string
	for _, r := range *subRecords {
		row := []string{
			r.Country, r.Operator, r.Service, r.Source, r.Msisdn,
			r.LatestSubject, r.Cycle, r.Adnet, r.Revenue, r.SubsDate.String,
			r.RenewalDate.String, r.FreemiumEndDate, r.UnsubsFrom, r.UnsubsDate.String,
			r.ServicePrice, r.Currency, r.ProfileStatus, r.Publisher,
			r.Trxid, r.Pixel, r.Handset, r.Browser, r.AttemptCharging,
			r.SuccessBilling,
		}
		subsData = append(subsData, row)
	}

	err = subW.WriteAll(subsData) // calls Flush internally
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(15 * time.Second)

	// delete file
	os.Remove(fileTrans)

	transCsv, err := os.Create(fileTrans)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	defer transCsv.Close()
	transW := csv.NewWriter(transCsv)
	defer transW.Flush()

	transHeaders := []string{
		"country", "operator", "service", "source", "msisdn",
		"event", "event_date", "cycle", "revenue", "charge_date",
		"currency", "publisher", "handset",
		"browser", "trxid", "telco_api_url", "telco_api_response",
		"sms_content", "status_sms",
	}
	transW.Write(transHeaders)

	transRecords, err := transactionService.SelectTransactionToCSV()
	if err != nil {
		log.Fatalf("error load table transactions: %s", err)
	}
	var transData [][]string
	for _, r := range *transRecords {
		row := []string{
			r.Country, r.Operator, r.Service, r.Source, r.Msisdn,
			r.Event, r.EventDate.String, r.Cycle, r.Revenue, r.ChargeDate.String,
			r.Currency, r.Publisher, r.Handset,
			r.Browser, r.TrxId, r.TelcoApiUrl, r.TelcoApiResponse,
			r.SmsContent, r.StatusSms,
		}
		transData = append(transData, row)
	}

	err = transW.WriteAll(transData) // calls Flush internally
	if err != nil {
		log.Fatal(err)
	}
}
