/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/futurama-dev/oauth-commander/client"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"log"
)

// infoCmd represents the info command
var clientInfoCmd = &cobra.Command{
	Use:   "info",
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

		clientId, err := cmd.Flags().GetString("clientId")
		if err != nil {
			log.Fatalln(err)
		}

		clients := client.Load()

		var client client.Client
		var ok bool

		if slug == "" && clientId == "" {
			// TODO use selected client
		} else if slug != "" {
			client, ok = clients.FindBySlug(slug)
		} else if clientId != "" {
			client, ok = clients.FindById(clientId)
		}

		if ok {
			info, err := yaml.Marshal(client)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(string(info))
		} else {
			fmt.Println("Client not found:", slug, " ", clientId)
		}
	},
}

func init() {
	clientCmd.AddCommand(clientInfoCmd)
	clientCmd.Flags().StringP("slug", "s", "", "Find info on client from slug")
	clientCmd.Flags().StringP("clientId", "d", "", "Find in on client from client ID")
}
