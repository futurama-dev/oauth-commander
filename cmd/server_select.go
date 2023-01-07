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
	"github.com/futurama-dev/oauth-commander/config"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

// selectCmd represents the select command
var serverSelectCmd = &cobra.Command{
	Use:   "select",
	Short: "Set the selected server",
	Run: func(cmd *cobra.Command, args []string) {
		slug, err := cmd.Flags().GetString("slug")
		if err != nil {
			log.Fatalln(err)
		}
		pathToFile := filepath.Join(config.ServerDir(), slug+".yaml")
		_, err = os.Stat(pathToFile)
		if err == os.ErrNotExist {
			log.Fatalln(err)
		} else if err == nil {
			err = config.SetSelectedServer(slug)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			log.Fatalln(err)
		}
	},
}

func init() {
	serverCmd.AddCommand(serverSelectCmd)
	serverSelectCmd.Flags().StringP("slug", "s", "", "select current server based on slug")
}
