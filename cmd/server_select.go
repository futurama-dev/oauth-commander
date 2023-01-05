/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/futurama-dev/oauth-commander/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

// selectCmd represents the select command
var serverSelectCmd = &cobra.Command{
	Use:   "select",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		slug, err := cmd.Flags().GetString("slug")
		if err != nil {
			log.Fatalln(err)
		}
		pathToFile := filepath.Join(config.ServerDir(), slug+".yaml")
		_, err = os.Stat(pathToFile)
		if err == os.ErrNotExist {
			log.Fatalln(err)
		} else if err == nil {
			viper.Set(config.SelectedServerSlug, slug)
			viper.WriteConfig()
		} else {
			log.Fatalln(err)
		}

	},
}

func init() {
	serverCmd.AddCommand(serverSelectCmd)
	serverSelectCmd.Flags().StringP("slug", "s", "", "select current server based on slug")
}
