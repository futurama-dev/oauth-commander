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
	"github.com/futurama-dev/oauth-commander/authorization"
	"github.com/futurama-dev/oauth-commander/client"
	"github.com/futurama-dev/oauth-commander/oauth2"
	"github.com/spf13/cobra"
	"net/url"
)

// authorizationResponseCmd represents the response command
var authorizationResponseCmd = &cobra.Command{
	Use:     "response",
	Aliases: []string{"resp"},
	Short:   "Receive an authorization response",
	RunE: func(cmd *cobra.Command, args []string) error {
		responseUrlStr, err := cmd.Flags().GetString("url")
		if err != nil {
			return err
		}

		if len(responseUrlStr) == 0 {
			return errors.New("response url not specified")
		}

		responseUrl, err := url.Parse(responseUrlStr)
		if err != nil {
			return err
		}

		authResponse, err := authorization.ProcessResponseUrl(responseUrl)

		if err != nil {
			errorResp, ok := err.(oauth2.ErrorResponse)
			if ok {
				errorResp.Println()
			} else {
				return err
			}
		} else {
			err = client.ProcessResponse(authResponse)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	authorizationCmd.AddCommand(authorizationResponseCmd)

	authorizationResponseCmd.Flags().StringP("url", "u", "", "The full response URL")
}
