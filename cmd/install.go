/*
 * Copyright © 2024 Piotr Kozak <piotrkrzysztofkozak@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */
package cmd

import (
	"fmt"
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/software"
	"github.com/pkk82/soft-ver-man/util/console"
	"github.com/spf13/cobra"
)

func InstallCmd(name, longName string, options software.InstallOptions) *cobra.Command {
	return &cobra.Command{
		Use:     "install [version]",
		Aliases: []string{"i", "install"},
		Short:   "Install software package from software directory",
		Long:    fmt.Sprintf("Install %v from software directory", longName),
		Args:    VersionArg,
		Run: func(cmd *cobra.Command, args []string) {
			plugin := domain.GetPlugin(name)
			err := software.Install(plugin, FirstOrEmpty(args), options)
			if err != nil {
				console.Fatal(err)
			}
		},
	}
}
