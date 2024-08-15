package cmd

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"github.com/wiliehidayat87/rmqp"
)

var (
	APP_NAME       string = getEnv("APP_NAME")
	APP_HOST       string = getEnv("APP_HOST")
	APP_PORT       string = getEnv("APP_PORT")
	APP_TZ         string = getEnv("APP_TZ")
	APP_URL        string = getEnv("APP_URL")
	URI_POSTGRES   string = getEnv("URI_POSTGRES")
	URI_REDIS      string = getEnv("URI_REDIS")
	RMQ_HOST       string = getEnv("RMQ_HOST")
	RMQ_USER       string = getEnv("RMQ_USER")
	RMQ_PASS       string = getEnv("RMQ_PASS")
	RMQ_PORT       string = getEnv("RMQ_PORT")
	RMQ_VHOST      string = getEnv("RMQ_VHOST")
	RMQ_URL        string = getEnv("RMQ_URL")
	ARPU_URL_SUB   string = getEnv("ARPU_URL_SUB")
	ARPU_URL_TRANS string = getEnv("ARPU_URL_TRANS")
	LOG_PATH       string = getEnv("LOG_PATH")
	PUBLIC_PATH    string = getEnv("PUBLIC_PATH")
	RMQ_PREFIX     string = getEnv("RMQ_PREFIX")
)

var (
	RMQ_EXCHANGE_TYPE         string = "direct"
	RMQ_DATA_TYPE             string = "application/json"
	RMQ_MO_EXCHANGE           string = "E_" + RMQ_PREFIX + "_MO"
	RMQ_MO_QUEUE              string = "Q_" + RMQ_PREFIX + "_MO"
	RMQ_RENEWAL_EXCHANGE      string = "E_" + RMQ_PREFIX + "_RENEWAL"
	RMQ_RENEWAL_QUEUE         string = "Q_" + RMQ_PREFIX + "_RENEWAL"
	RMQ_RETRY_DP_EXCHANGE     string = "E_" + RMQ_PREFIX + "_RETRY_DP"
	RMQ_RETRY_DP_QUEUE        string = "Q_" + RMQ_PREFIX + "_RETRY_DP"
	RMQ_RETRY_FP_EXCHANGE     string = "E_" + RMQ_PREFIX + "_RETRY_FP"
	RMQ_RETRY_FP_QUEUE        string = "Q_" + RMQ_PREFIX + "_RETRY_FP"
	RMQ_RETRY_INSUFF_EXCHANGE string = "E_" + RMQ_PREFIX + "_RETRY_INSUFF"
	RMQ_RETRY_INSUFF_QUEUE    string = "Q_" + RMQ_PREFIX + "_RETRY_INSUFF"
	RMQ_NOTIF_EXCHANGE        string = "E_" + RMQ_PREFIX + "_NOTIF"
	RMQ_NOTIF_QUEUE           string = "Q_" + RMQ_PREFIX + "_NOTIF"
	RMQ_PB_MO_EXCHANGE        string = "E_" + RMQ_PREFIX + "_POSTBACK_MO"
	RMQ_PB_MO_QUEUE           string = "Q_" + RMQ_PREFIX + "_POSTBACK_MO"
	RMQ_PB_MT_EXCHANGE        string = "E_" + RMQ_PREFIX + "_POSTBACK_MT"
	RMQ_PB_MT_QUEUE           string = "Q_" + RMQ_PREFIX + "_POSTBACK_MT"
	RMQ_TRAFFIC_EXCHANGE      string = "E_" + RMQ_PREFIX + "_TRAFFIC"
	RMQ_TRAFFIC_QUEUE         string = "Q_" + RMQ_PREFIX + "_TRAFFIC"
	RMQ_DAILYPUSH_EXCHANGE    string = "E_" + RMQ_PREFIX + "_BQ_DAILYPUSH"
	RMQ_DAILYPUSH_QUEUE       string = "Q_" + RMQ_PREFIX + "_BQ_DAILYPUSH"
	MT_FIRSTPUSH              string = "MT_FIRSTPUSH"
	ACT_RENEWAL               string = "RENEWAL"
	ACT_RETRY_FP              string = "RETRY_FP"
	ACT_RETRY_DP              string = "RETRY_DP"
	ACT_RETRY_INSUFF          string = "RETRY_INSUFF"
	ACT_CSV                   string = "CSV"
	ACT_UPLOAD_CSV            string = "UPLOAD_CSV"
)

var (
	rootCmd = &cobra.Command{
		Use:   "cobra-cli",
		Short: "A generator for Cobra based Applications",
		Long:  `Cobra is a CLI library for Go that empowers applications.`,
	}
)

func init() {
	// setup timezone
	loc, _ := time.LoadLocation(APP_TZ)
	time.Local = loc

	/**
	 * Listener Service
	 */
	rootCmd.AddCommand(listenerCmd)

	/**
	 * Consumer Service
	 */
	rootCmd.AddCommand(consumerMOCmd)
	rootCmd.AddCommand(consumerRenewalCmd)
	rootCmd.AddCommand(consumerPostbackMOCmd)
	rootCmd.AddCommand(consumerPostbackMTCmd)
	rootCmd.AddCommand(consumerNotifCmd)
	rootCmd.AddCommand(consumerTrafficCmd)
	rootCmd.AddCommand(consumerDailypushCmd)
	rootCmd.AddCommand(consumerTestCmd)

	/**
	 * Publisher Service
	 */
	rootCmd.AddCommand(publisherCSVCmd)
	rootCmd.AddCommand(publisherUploadCSVCmd)

}

func Execute() error {
	return rootCmd.Execute()
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		log.Panicf("Error %v", key)
	}
	return value
}

// Connect to postgresql
func connectPgsql() (*sql.DB, error) {
	db, err := sql.Open("postgres", URI_POSTGRES)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Connect to redis
func connectRedis() (*redis.Client, error) {
	opts, err := redis.ParseURL(URI_REDIS)
	if err != nil {
		return nil, err
	}
	return redis.NewClient(opts), nil
}

// Connect to rabbitmq
func connectRabbitMq() (rmqp.AMQP, error) {
	var rb rmqp.AMQP
	port, _ := strconv.Atoi(RMQ_PORT)
	rb.SetAmqpURL(RMQ_HOST, port, RMQ_USER, RMQ_PASS, RMQ_VHOST)
	errConn := rb.SetUpConnectionAmqp()
	if errConn != nil {
		return rb, errConn
	}
	return rb, nil
}
