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
package java

import (
	"github.com/pkk82/soft-ver-man/cmd"
	"github.com/pkk82/soft-ver-man/software/java"
)

var Cmd = cmd.MainCmd(java.Name, java.LongName, java.Aliases)

var verifyChecksumFetch bool
var verifyChecksumInstall bool

func init() {
	cmd.RootCmd.AddCommand(Cmd)
	fetchCmd := cmd.FetchCmd(java.Name, java.LongName, &verifyChecksumFetch)
	fetchCmd.Flags().BoolVarP(&verifyChecksumFetch, "verify-checksum", "c", false, "Verify checksum of downloaded file")
	Cmd.AddCommand(fetchCmd)
	installCmd := cmd.InstallCmd(java.Name, java.LongName, cmd.InstallOptions{VerifyChecksum: &verifyChecksumInstall})
	installCmd.Flags().BoolVarP(&verifyChecksumInstall, "verify-checksum", "c", false, "Verify checksum of downloaded file")
	Cmd.AddCommand(installCmd)
	Cmd.AddCommand(cmd.UninstallCmd(java.Name, java.LongName))
	Cmd.AddCommand(cmd.InstalledCmd(java.Name))
	Cmd.AddCommand(cmd.AvailableCmd(java.Name, java.LongName))
}
