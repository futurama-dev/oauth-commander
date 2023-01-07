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
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// docsCmd represents the docs command
var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Generate documentation for oauth-commander",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := cmd.Flags().GetString("dir")
		if err != nil {
			return err
		}

		if dir == "" {
			if dir, err = ioutil.TempDir("", "oauth-commander-docs-"); err != nil {
				return err
			}
		}

		markdown, err := cmd.Flags().GetBool("markdown")
		if err != nil {
			return err
		}

		man, err := cmd.Flags().GetBool("man")
		if err != nil {
			return err
		}

		yaml, err := cmd.Flags().GetBool("yaml")
		if err != nil {
			return err

		}

		rest, err := cmd.Flags().GetBool("rest")
		if err != nil {
			return err

		}

		return docsAction(dir, markdown, man, yaml, rest)
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)

	docsCmd.Flags().StringP("dir", "d", "", "Destination directory for the generated docs. Default: a temp folder.")
	docsCmd.Flags().BoolP("markdown", "m", true, "Generated Markdown documentation")
	docsCmd.Flags().BoolP("man", "p", false, "Generated man pages")
	docsCmd.Flags().BoolP("yaml", "y", false, "Generated YAML documentation")
	docsCmd.Flags().BoolP("rest", "r", false, "Generated ReStructured Text documentation")
}

func docsAction(dir string, markdown bool, man bool, yaml bool, rest bool) error {
	if markdown {
		if err := doc.GenMarkdownTree(rootCmd, dir); err != nil {
			return err
		}
		fmt.Println("Markdown documentation generated in:", dir)
	}

	if man {
		if err := doc.GenManTree(rootCmd, nil, dir); err != nil {
			return err
		}
		fmt.Println("man pages generated in:", dir)
	}

	if yaml {
		if err := doc.GenYamlTree(rootCmd, dir); err != nil {
			return err
		}
		fmt.Println("YAML format documentation generated in:", dir)
	}

	if rest {
		if err := doc.GenManTree(rootCmd, nil, dir); err != nil {
			return err
		}
		fmt.Println("ReStructured Text documentation generated in:", dir)
	}

	return nil
}
