/*
Copyright © 2023 Piotr Kozak <piotrkrzysztofkozak@gmail.com>

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
package node

import (
	"fmt"
	"github.com/pkk82/soft-ver-man/console"
	"github.com/pkk82/soft-ver-man/download"
	"github.com/pkk82/soft-ver-man/version"
	"path/filepath"
)

func FetchVersion(inputVersion, softwareDownloadDir string, verify bool) {
	versions := getSupportedVersions()
	versionIds := make([]string, len(versions))
	for i, v := range versions {
		versionIds[i] = v.Id
	}
	foundVersion, index, err := version.FindVersion(inputVersion, versionIds)
	if err != nil {
		console.Fatal(err)
	}
	matchingVersion := versions[index]
	nodeDir := filepath.Join(softwareDownloadDir, "node")
	download.FetchFile(matchingVersion.DownloadLink(), nodeDir, matchingVersion.FileName)

	if verify {
		fetchPGPKeys(softwareDownloadDir)
		shaSumFile := fmt.Sprintf("node-%v-%v", foundVersion, ShaSumFileName)
		download.FetchFileSilently(matchingVersion.SumsLink(), nodeDir, shaSumFile)
		shaSumSigFile := fmt.Sprintf("node-%v-%v", foundVersion, ShaSumSigFileName)
		download.FetchFileSilently(matchingVersion.SumsSigLink(), nodeDir, shaSumSigFile)
	}

}
