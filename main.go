package main

import (
	"github.com/duchenhao/backend-demo/internal/bus"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/duchenhao/backend-demo/internal/conf"
	"github.com/duchenhao/backend-demo/internal/dao"
	"github.com/duchenhao/backend-demo/internal/log"
	"github.com/duchenhao/backend-demo/internal/server"
)

var (
	configPath = flag.StringP("config", "c", "", "path to config file")
)

var (
	rootCmd = &cobra.Command{
		Use: conf.Core.Name,
	}

	serverCmd = &cobra.Command{
		Use:   "run",
		Short: "start http server",
		Run: func(cmd *cobra.Command, args []string) {
			srv := server.NewHttpServer()
			srv.Run()
		},
	}

	migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "migrate database schema",
		Run: func(cmd *cobra.Command, args []string) {
			dao.Migrate()
		},
	}
)

func main() {
	flag.Parse()

	conf.Init(*configPath)

	log.Init()
	defer log.Sync()

	dao.Init()
	defer dao.Close()

	defer bus.Close()

	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.Execute()
}
