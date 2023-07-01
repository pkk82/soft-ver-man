package node

import (
	"github.com/pkk82/soft-ver-man/download"
	"strings"
)

func fetchPGPKeys(softwareDownloadDir string) {

	for _, fingerprint := range getFingerprints() {
		download.FetchFile("https://keys.openpgp.org/vks/v1/by-fingerprint/"+fingerprint, softwareDownloadDir+"/node-pgp-keys", fingerprint)
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
