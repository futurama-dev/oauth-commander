/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/futurama-dev/oauth-commander/client"
	"github.com/futurama-dev/oauth-commander/config"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var clientListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		clients := client.Load()

		fmt.Println("Selected server:", config.GetSelectedServer())
		fmt.Println("Number of clients:", len(clients))

		for _, client := range clients {
			fmt.Println(client.Slug, client.ServerSlug, client.Id, client.CreatedAt)
		}
	},
}

func init() {
	clientCmd.AddCommand(clientListCmd)
}
