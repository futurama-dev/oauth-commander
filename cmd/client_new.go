/*
Copyright © 2023 futurama-dev

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"github.com/futurama-dev/oauth-commander/client"
	"github.com/futurama-dev/oauth-commander/config"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// clientNewCmd represents the new command
var clientNewCmd = &cobra.Command{
	Use:   "new",
	Short: "create a new client, under the selected server",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		serverSlug := config.GetSelectedServer()
		if len(serverSlug) == 0 {
			fmt.Println("no server selected")
			os.Exit(1)
		}

		clientId, err := cmd.Flags().GetString("id")
		if err != nil {
			log.Fatalln(err)
		}
		if foundClient, ok := client.FindById(clientId); ok {
			fmt.Println("A client with this client id already exists:", foundClient.Slug)
			os.Exit(1)
		}

		clientSlug := client.NextSlug()
		redirectUri, err := cmd.Flags().GetString("redirect-uri")
		if err != nil {
			log.Fatalln(err)
		}
		clientX := client.Client{
			Slug:         clientSlug,
			ServerSlug:   serverSlug,
			CreatedAt:    time.Now().Truncate(time.Second),
			Id:           clientId,
			SecretHandle: serverSlug + "_" + clientSlug,
			RedirectURIs: []string{redirectUri},
		}

		err = client.Save(clientX)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("Client saved:", clientX.Slug)

		slugs := viper.GetStringMapString(config.SelectedClientSlugs)
		if len(slugs[serverSlug]) == 0 {
			fmt.Println("Currently no selected client, selecting this one")
			slugs[serverSlug] = clientSlug
			viper.Set(config.SelectedClientSlugs, slugs)
			err = viper.WriteConfig()
			if err != nil {
				log.Fatalln(err)
			}
		}
	},
}

func init() {
	clientCmd.AddCommand(clientNewCmd)

	clientNewCmd.Flags().String("id", "", "the client id")
	clientNewCmd.Flags().String("redirect-uri", "", "the first redirect URI")

}
