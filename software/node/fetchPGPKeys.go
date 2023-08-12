package node

import (
	"github.com/pkk82/soft-ver-man/util/download"
	"strings"
)

func fetchPGPKeys(softwareDownloadDir string) []string {
	var paths = make([]string, 0)
	for _, fingerprint := range getFingerprints() {
		paths = append(paths, download.FetchFileSilently("https://keys.openpgp.org/vks/v1/by-fingerprint/"+fingerprint, softwareDownloadDir+"/node-pgp-keys", fingerprint))
	}
	return paths
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
