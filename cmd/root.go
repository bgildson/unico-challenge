package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "unico-challenge",
	Short: "This application contains the solution for unico challenge",
	Long:  `The challenge was to create a CRUD rest api for feiras livres and an command to import the csv containing the feiras livres.`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
