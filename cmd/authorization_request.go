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
	"errors"
	"fmt"
	"github.com/futurama-dev/oauth-commander/authorization"
	"github.com/futurama-dev/oauth-commander/config"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

// authorizationRequestCmd represents the request command
var authorizationRequestCmd = &cobra.Command{
	Use:     "request",
	Aliases: []string{"req"},
	Short:   "Initiate an authorization request",
	RunE: func(cmd *cobra.Command, args []string) error {
		serverSlug := config.GetSelectedServer()
		if len(serverSlug) == 0 {
			return errors.New("no selected server")
		}

		clientSlug := config.GetSelectedClient()
		if len(clientSlug) == 0 {
			return errors.New("no selected client")
		}

		scope, err := cmd.Flags().GetStringArray("scope")
		if err != nil {
			return err
		}

		code, err := cmd.Flags().GetBool("code")
		if err != nil {
			return err
		}

		token, err := cmd.Flags().GetBool("token")
		if err != nil {
			return err
		}

		id_token, err := cmd.Flags().GetBool("id-token")
		if err != nil {
			return err
		}

		redirectUri, err := cmd.Flags().GetString("redirect-uri")
		if err != nil {
			return err
		}

		authUrl, err := authorization.GenerateAuthorizationRequestUrl(serverSlug, clientSlug, code, token, id_token, scope, redirectUri, verbose)
		if err != nil {
			return err
		}

		fmt.Println(authUrl)

		open, err := cmd.Flags().GetBool("open")
		if err != nil {
			return err
		}

		if open {
			err = browser.OpenURL(authUrl)
			if err != nil {
				return err
			}
		}

		listen, err := cmd.Flags().GetBool("listen")
		if err != nil {
			return err
		}

		if listen {

		}

		return nil
	},
}

func init() {
	authorizationCmd.AddCommand(authorizationRequestCmd)

	authorizationRequestCmd.Flags().StringArrayP("scope", "s", []string{}, "List of scopes to add to the request")
	authorizationRequestCmd.Flags().StringP("redirect-uri", "r", "", "The redirect URI to use with the request. Must be one of the configured ones for the client. Default to first one.")

	authorizationRequestCmd.Flags().BoolP("code", "c", true, "Add response type code")
	authorizationRequestCmd.Flags().BoolP("token", "t", false, "Add response type token")
	authorizationRequestCmd.Flags().BoolP("id-token", "i", false, "Add response type id_token")

	authorizationRequestCmd.Flags().BoolP("open", "o", false, "Open authorization request URL in default browser")
	authorizationRequestCmd.Flags().BoolP("listen", "l", false, "Start a web server and listen the response. The redirect URI must be localhost and handled by OAuth Commander")
}
