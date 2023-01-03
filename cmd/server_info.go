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
	"github.com/futurama-dev/oauth-commander/server"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"log"
	"strings"
)

// infoCmd represents the info command
var serverInfoCmd = &cobra.Command{
	Use:   "info <slug|issuer>",
	Short: "Show detailed information for a server",
	Long: `Show all the configuration information for a server. The server can be specified either as the server
slug or the issuer.`,
	Run: func(cmd *cobra.Command, args []string) {
		servers := server.Load()
		slug_or_issuer := args[0]

		var server server.Server
		var ok bool

		if strings.HasPrefix(slug_or_issuer, "https://") {
			server, ok = servers.FindByIssuer(slug_or_issuer)
		} else {
			server, ok = servers.FindBySlug(slug_or_issuer)
		}

		if ok {
			info, err := yaml.Marshal(server)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(string(info))
		} else {
			fmt.Println("Server not found:", slug_or_issuer)
		}
	},
}

func init() {
	serverCmd.AddCommand(serverInfoCmd)
}
