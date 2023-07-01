package node

import (
	"fmt"
	"github.com/pkk82/soft-ver-man/download"
	"os"
	"strings"
)

func fetchPGPKeys(softwareDownloadDir string) {

	for _, fingerprint := range getFingerprints() {

		err := download.FetchFile("https://keys.openpgp.org/vks/v1/by-fingerprint/"+fingerprint, softwareDownloadDir+"/node-pgp-keys", fingerprint)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

}

func getFingerprints() []string {
	result := make([]string, 0)
	for _, line := range strings.Split(gpgKeys, "\n") {
		lineSplit := strings.Split(line, "> ")
		if len(lineSplit) > 1 {
			result = append(result, lineSplit[1])
		}
	}
	return result
}
