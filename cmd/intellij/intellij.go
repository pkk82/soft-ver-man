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
package intellij

import (
	"github.com/pkk82/soft-ver-man/cmd"
	"github.com/pkk82/soft-ver-man/software"
	"github.com/pkk82/soft-ver-man/software/intellij"
)

var verifyChecksum bool
var main bool
var here bool
var Cmd = cmd.MainCmd(intellij.Name, intellij.LongName, intellij.Aliases)

func init() {
	cmd.RootCmd.AddCommand(Cmd)
	installCmd := cmd.InstallCmd(intellij.Name, intellij.LongName, software.InstallOptions{VerifyChecksum: &verifyChecksum, Main: &main, Here: &here})
	Cmd.AddCommand(installCmd)
	Cmd.AddCommand(cmd.UninstallCmd(intellij.Name, intellij.LongName))

	installCmd.Flags().BoolVarP(&verifyChecksum, "verify-checksum", "c", false, "Verify checksum of downloaded file")
	installCmd.Flags().BoolVarP(&main, "main", "m", false, "Make package main version (default: false)")
	installCmd.Flags().BoolVarP(&here, "here", "x", false, "Use package in the current directory via .direnv (default: false)")
}
