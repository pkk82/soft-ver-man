package java

import (
	"encoding/json"
	"github.com/pkk82/soft-ver-man/util/console"
	"io"
	"net/http"
)

type ExtendedPackage struct {
	Id     string `json:"package_uuid"`
	Sha256 string `json:"sha256_hash"`
}

func getExtendedPackage(packageId string) (ExtendedPackage, error) {
	url := PackagesAPIURL + "/" + packageId
	resp, err := http.Get(url)
	if err != nil {
		return ExtendedPackage{}, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			console.Fatal(err)
		}
	}(resp.Body)
	var extendedPackage ExtendedPackage
	err = json.NewDecoder(resp.Body).Decode(&extendedPackage)
	if err != nil {
		return ExtendedPackage{}, err
	}
	return extendedPackage, nil
}
