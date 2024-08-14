package cmd

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html/v2"
	"github.com/idprm/go-xl-direct/internal/domain/repository"
	"github.com/idprm/go-xl-direct/internal/handler"
	"github.com/idprm/go-xl-direct/internal/logger"
	"github.com/idprm/go-xl-direct/internal/services"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"github.com/wiliehidayat87/rmqp"
)

var listenerCmd = &cobra.Command{
	Use:   "listener",
	Short: "Listener Service CLI",
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
		rmq.SetUpChannel(RMQ_EXCHANGE_TYPE, true, RMQ_TRAFFIC_EXCHANGE, true, RMQ_TRAFFIC_QUEUE)

		r := routerUrl(db, rds, rmq, logger)

		log.Fatal(r.Listen(":" + APP_PORT))

	},
}

func routerUrl(db *sql.DB, rds *redis.Client, rmq rmqp.AMQP, logger *logger.Logger) *fiber.App {

	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	engine := html.New(path+"/views", ".html")

	/**
	 * Init Fiber
	 */
	r := fiber.New(fiber.Config{
		Views: engine,
	})

	/**
	 * Initialize default config
	 */
	r.Use(cors.New())

	/**
	 * Access log on browser
	 */
	r.Use(LOG_PATH, filesystem.New(filesystem.Config{
		Root:         http.Dir(LOG_PATH),
		Browse:       true,
		Index:        "index.html",
		NotFoundFile: "404.html",
		MaxAge:       3600,
	}))

	r.Static("/static", path+"/"+PUBLIC_PATH)

	serviceRepo := repository.NewServiceRepository(db)
	serviceService := services.NewServiceService(serviceRepo)
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)
	verifyRepo := repository.NewVerifyRepository(rds)
	verifyService := services.NewVerifyService(verifyRepo)

	h := handler.NewIncomingHandler(
		rmq,
		logger,
		serviceService,
		subscriptionService,
		verifyService,
	)

	v1 := r.Group("v1")
	v1.Post("sub", h.CreateSubscription)
	v1.Post("otp", h.ConfirmOTP)
	v1.Post("refund", h.Refund)
	v1.Post("unsub", h.Unsubscribe)

	r.Post("notify", h.MessageOriginated)

	return r

}
