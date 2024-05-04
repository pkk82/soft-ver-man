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

package shell

import (
	"fmt"
	"github.com/pkk82/soft-ver-man/config"
	"github.com/pkk82/soft-ver-man/domain"
)

func bashToLoad(fileName string) string {
	return fmt.Sprintf("[[ -s \"$HOME/%v/%v\" ]] && source \"$HOME/%v/%v\"",
		config.HomeConfigDir, fileName, config.HomeConfigDir, fileName)
}

func makeRcName(name string) string {
	return fmt.Sprintf(".%vrc", name)
}

func exportVariable(name, value string) string {
	return fmt.Sprintf("export %v=\"%v\"", name, value)
}

func exportRefPathVariable(name, refVar, path string) string {
	return fmt.Sprintf("export %v=\"$%v/%v\"", name, refVar, path)
}

func exportHomeVariable(name, refVar, path string) string {
	return fmt.Sprintf("export %v_HOME=\"$%v/%v\"", name, refVar, path)
}

func exportHomeMajorVersionVariable(name string, version domain.Version, refVar, path string) string {
	return fmt.Sprintf("export %v_%v_HOME=\"$%v/%v\"", name, version.Major(), refVar, path)
}

func exportHomeMinorVersionVariable(name string, version domain.Version, refVar, path string) string {
	return fmt.Sprintf("export %v_%v_%v_HOME=\"$%v/%v\"", name, version.Major(), version.Minor(), refVar, path)
}
