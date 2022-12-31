/*
Copyright © 2022 futurama-dev

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
	"github.com/futurama-dev/oauth-commander/discovery"
	"github.com/futurama-dev/oauth-commander/oauth2"
	"github.com/futurama-dev/oauth-commander/oidc"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// discoverCmd represents the discover command
var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		issuer, err := cmd.Flags().GetString("issuer")
		layerType, err := cmd.Flags().GetString("type")

		if err != nil {
			log.Fatalln(err)
		}

		switch layerType {
		case "all":
			fmt.Println("trying all types")
			fmt.Println("OIDC")
			fmt.Println(oidcDiscovery(issuer))
			fmt.Println("OAuth2")
			fmt.Println(oauth2Discovery(issuer))
		case "oidc":
			fmt.Println("trying oidc")
			fmt.Println(oidcDiscovery(issuer))
		case "oauth2":
			fmt.Println("trying oauth 2")
			fmt.Println(oauth2Discovery(issuer))
		default:
			fmt.Println("invalid type inputted: ", layerType)
		}
	},
}

func init() {
	serverCmd.AddCommand(discoverCmd)

	discoverCmd.Flags().StringP("issuer", "i", "", "Authorization server issuer to run discovery against.")
	discoverCmd.Flags().StringP("type", "t", "all", "Which top layer to use on top of OAuth 2")
}

func oidcDiscovery(issuer string) string {
	discoveryUrl, err := oidc.BuildDiscoveryUrl(issuer)

	if err != nil {
		log.Fatalln(err)
	}

	oidcConfig, err := discovery.FetchDiscovery(discoveryUrl)

	if err == discovery.NotFoundErr {
		fmt.Println("OpenID Connect discovery not found!")
		os.Exit(0)
	}

	if err != nil {
		log.Fatalln(err)
	}

	return oidcConfig
}

func oauth2Discovery(issuer string) string {
	discoveryUrl, err := oauth2.BuildDiscoveryUrl(issuer)

	if err != nil {
		log.Fatalln(err)
	}

	oauth2Config, err := discovery.FetchDiscovery(discoveryUrl)

	if err == discovery.NotFoundErr {
		fmt.Println("OAuth 2 Connect discovery not found!")
		os.Exit(0)
	}

	if err != nil {
		log.Fatalln(err)
	}

	return oauth2Config
}
