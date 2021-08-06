package cmd

import (
	"database/sql"
	"io"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/lib/pq" // init postgres database driver
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	feiralivreController "github.com/bgildson/unico-challenge/controller/feiralivre"
	feiralivreRepository "github.com/bgildson/unico-challenge/repository/feiralivre"
	"github.com/bgildson/unico-challenge/server"
	"github.com/bgildson/unico-challenge/server/parser"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the rest api server",
	Run: func(cmd *cobra.Command, args []string) {
		paginationDefaultLimit, _ := strconv.Atoi(os.Getenv("PAGINATION_DEFAULT_LIMIT"))
		paginationMaxLimit, _ := strconv.Atoi(os.Getenv("PAGINATION_MAX_LIMIT"))
		config := server.NewConfig(
			os.Getenv("ENVIRONMENT"),
			os.Getenv("PORT"),
			os.Getenv("DATABASE_URL"),
			os.Getenv("LOGS_PATH"),
			paginationDefaultLimit,
			paginationMaxLimit,
		)
		if err := config.Validate(); err != nil {
			logrus.Error(err)
		}

		logsFile, err := os.OpenFile(config.LogsPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
		if err != nil {
			logrus.Error(err)
		}
		defer logsFile.Close()

		logsOut := io.MultiWriter(os.Stdout, logsFile)

		logrus.SetOutput(logsOut)
		logrus.SetFormatter(&logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyFunc:  "caller",
				logrus.FieldKeyMsg:   "message",
			},
		})
		if config.Environment == server.ProductionEnvironment {
			logrus.SetLevel(logrus.WarnLevel)
		} else {
			logrus.SetReportCaller(true)
			logrus.SetLevel(logrus.DebugLevel)
		}

		app := fiber.New()

		app.Use(logger.New(
			logger.Config{
				Format:     `{"timestamp":"${time}", "method":"${method}", "path":"${path}", "query_params":"${queryParams}", "status":${status}, "latency":"${latency}", "pid":${pid}"}` + "\n",
				TimeFormat: "2006-01-02T15:04:05Z07:00",
				Output:     logsOut,
			},
		))

		db, err := sql.Open("postgres", config.DatabaseURL)
		if err != nil {
			logrus.Error(err)
		}
		defer db.Close()

		feiralivreRepo := feiralivreRepository.NewPostgresRepository(db)
		queryParamsParser := parser.NewQueryParamsParser(config.PaginationDefaultLimit, config.PaginationMaxLimit)
		feiralivreCtrl := feiralivreController.New(feiralivreRepo, queryParamsParser)
		feiralivreCtrl.Register(app, "/feiras-livres")

		if err := app.Listen(":" + config.Port); err != nil {
			logrus.Error(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
