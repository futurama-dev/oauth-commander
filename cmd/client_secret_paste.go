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
	"github.com/atotto/clipboard"
	"github.com/futurama-dev/oauth-commander/client"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// clientSecretPasteCmd represents the paste command
var clientSecretPasteCmd = &cobra.Command{
	Use:   "paste <client_slug>",
	Short: "Save the client secret by grabbing the value from the clipboard",
	Run: func(cmd *cobra.Command, args []string) {
		slug := args[0]
		clientX, ok := client.FindBySlug(slug)

		if ok {
			secret, err := clipboard.ReadAll()
			if err != nil {
				log.Fatalln(err)
			}

			if len(secret) == 0 {
				fmt.Println("Clipboard is empty")
				os.Exit(1)
			}

			err = clientX.SetSecret(secret)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			fmt.Println("Client not found:", slug)
		}
	},
}

func init() {
	clientSecretCmd.AddCommand(clientSecretPasteCmd)
}
