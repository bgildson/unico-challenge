package cmd

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	feiralivreRepo "github.com/bgildson/unico-challenge/repository/feiralivre"
	feiralivreServ "github.com/bgildson/unico-challenge/service/feiralivre"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Imports a CSV to the database",
	Run: func(cmd *cobra.Command, _ []string) {
		databaseURL := os.Getenv("DATABASE_URL")
		if databaseURL == "" {
			logrus.Error("could not load the database connection string")
		}
		if databaseURLFlag := cmd.Flag("dsn"); databaseURLFlag.Value.String() != "" {
			databaseURL = databaseURLFlag.Value.String()
		}

		db, err := sql.Open("postgres", databaseURL)
		if err != nil {
			logrus.Errorf("could not connect to the database: %v", err)
		}
		defer db.Close()

		fs := afero.NewOsFs()

		r := feiralivreRepo.NewPostgresRepository(db)

		s := feiralivreServ.New(fs, r)

		path := cmd.Flag("file")

		msg, err := s.Import(path.Value.String())
		if err != nil {
			logrus.Errorf("could not import: %v", err)
		}

		fmt.Println(msg)
	},
}

func init() {
	importCmd.Flags().StringP("dsn", "d", "", "The Data Source Name that should be used to connect in the database.")
	importCmd.Flags().StringP("file", "f", "", "CSV file path that should be imported.")
	importCmd.MarkFlagRequired("file")

	rootCmd.AddCommand(importCmd)
}
