/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/futurama-dev/oauth-commander/server"
	"github.com/spf13/cobra"
	"log"
)

// removeCmd represents the remove command
var serverRemoveCmd = &cobra.Command{
	Use:   "remove",
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
		fmt.Println("Attempting to remove saved server from given slug.")
		err = server.RemoveBySlug(slug)
		if err != nil {
			fmt.Println("Failed to remove saved server: ", err)
		} else {
			fmt.Println("Success, server removed: ", slug)
		}
	},
}

func init() {
	serverCmd.AddCommand(serverRemoveCmd)
	serverRemoveCmd.Flags().StringP("slug", "s", "", "remove saved server by slug")
}
