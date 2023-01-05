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
	"github.com/futurama-dev/oauth-commander/server"
	"github.com/spf13/cobra"
	"log"
	"time"
)

// serverDiscoverCmd represents the discover command
var serverDiscoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		issuer, err := cmd.Flags().GetString("issuer")
		if err != nil {
			log.Fatalln(err)
		}

		layerType, err := cmd.Flags().GetString("type")
		if err != nil {
			log.Fatalln(err)
		}

		save, err := cmd.Flags().GetBool("save")
		if err != nil {
			log.Fatalln(err)
		}

		update, err := cmd.Flags().GetBool("update")
		if err != nil {
			log.Fatalln(err)
		}

		switch layerType {
		case "all":
			var continueSave bool
			fmt.Println("trying all types")
			fmt.Println("---------------")
			fmt.Println("OIDC")
			oidcConfig, success := oidcDiscovery(issuer, save, update)
			if success {
				continueSave = false
			} else {
				continueSave = true
			}
			if oidcConfig != "" {
				fmt.Println(oidcConfig)
			}
			fmt.Println("---------------")
			fmt.Println("OAuth2")
			if !continueSave && (save || update) {
				save, update = false, false
			}
			oauth2Config, success := oauth2Discovery(issuer, save, update)
			if success {
				continueSave = false
			} else {
				continueSave = true
			}
			if oauth2Config != "" {
				fmt.Println(oauth2Config)
			}
		case "oidc":
			fmt.Println("trying oidc")
			oidcConfig, success := oidcDiscovery(issuer, save, update)
			if success {

			}
			if oidcConfig != "" {
				fmt.Println(oidcConfig)
			}
		case "oauth2":
			fmt.Println("trying oauth 2")
			oauth2Config, success := oauth2Discovery(issuer, save, update)
			if success {

			}
			if oauth2Config != "" {
				fmt.Println(oauth2Config)
			}
		default:
			fmt.Println("invalid type inputted: ", layerType)
		}
	},
}

func init() {
	serverCmd.AddCommand(serverDiscoverCmd)

	serverDiscoverCmd.Flags().StringP("issuer", "i", "", "Authorization server issuer to run discovery against.")
	serverDiscoverCmd.Flags().StringP("type", "t", "all", "Which top layer to use on top of OAuth 2")
	serverDiscoverCmd.Flags().BoolP("save", "s", false, "whether or not you want to save the returned JSON data")
	serverDiscoverCmd.Flags().BoolP("update", "u", false, "used to update an existing saved server")
}

func oidcDiscovery(issuer string, save bool, update bool) (string, bool) {
	discoveryUrl, err := oidc.BuildDiscoveryUrl(issuer)

	if err != nil {
		log.Fatalln(err)
	}

	oidcConfig, err := discovery.FetchDiscovery(discoveryUrl)

	if save || update {
		savedConfig, err := discovery.ParseMetaData(oidcConfig)

		if err == discovery.InvalidJSONErr {
			fmt.Println(discovery.InvalidJSONErr)
			return "", false
		} else {
			fmt.Println("Saved config: ", savedConfig)
		}

		slug, err := server.IssuerToSlug(issuer)
		if err != nil {
			log.Fatalln(err)
		}

		svr := server.Server{
			Slug:      slug,
			Type:      "oidc",
			CreatedAt: time.Now().Truncate(time.Second),
			Metadata:  savedConfig,
		}

		if save {
			err = server.Save(svr)
			fmt.Println(err)
		} else {
			err = server.Update(svr)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	if err == discovery.NotFoundErr {
		fmt.Println("OpenID Connect discovery not found!")
	} else if err != nil {
		fmt.Println(err)
		return "", false
	}

	return oidcConfig, true
}

func oauth2Discovery(issuer string, save bool, update bool) (string, bool) {
	discoveryUrl, err := oauth2.BuildDiscoveryUrl(issuer)

	if err != nil {
		log.Fatalln(err)
	}

	oauth2Config, err := discovery.FetchDiscovery(discoveryUrl)

	if save || update {
		savedConfig, err := discovery.ParseMetaData(oauth2Config)

		if err == discovery.InvalidJSONErr {
			fmt.Println(discovery.InvalidJSONErr)
			return "", false
		} else {
			fmt.Println("Saved config: ", savedConfig)
		}

		slug, err := server.IssuerToSlug(issuer)
		if err != nil {
			log.Fatalln(err)
		}

		svr := server.Server{
			Slug:      slug,
			Type:      "oauth2",
			CreatedAt: time.Now().Truncate(time.Second),
			Metadata:  savedConfig,
		}

		if save {
			server.Save(svr)
		} else {
			server.Update(svr)
		}
	}

	if err == discovery.NotFoundErr {
		fmt.Println("OAuth 2 Connect discovery not found!")
	} else if err != nil {
		fmt.Println(err)
		return "", false
	}

	return oauth2Config, true
}
