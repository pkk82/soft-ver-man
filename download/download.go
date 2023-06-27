package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func FetchFile(url, downloadDir, fileName string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(response.Body)

	err = os.MkdirAll(downloadDir, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(downloadDir, fileName))
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
