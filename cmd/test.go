package cmd

import (
	"log"

	"github.com/idprm/go-xl-direct/internal/domain/entity"
	"github.com/idprm/go-xl-direct/internal/logger"
	"github.com/idprm/go-xl-direct/internal/providers/telco"
	"github.com/spf13/cobra"
)

var consumerTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Consumer Test Service CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		/**
		 * SETUP LOG
		 */
		logger := logger.NewLogger()
		t := telco.NewTelco(logger, &entity.Subscription{}, &entity.Service{}, &entity.Content{})
		auth, err := t.OAuth()
		if err != nil {
			log.Println(err.Error())
		}

		log.Println(string(auth))
	},
}
