package cmd

import (
	"database/sql"
	"encoding/csv"
	"log"
	"os"
	"time"

	"github.com/idprm/go-xl-direct/internal/domain/repository"
	"github.com/idprm/go-xl-direct/internal/services"
	"github.com/spf13/cobra"
)

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
