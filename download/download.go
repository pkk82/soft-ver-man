package download

import (
	"fmt"
	"github.com/pkk82/soft-ver-man/console"
	"github.com/schollz/progressbar/v3"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func FetchFile(url, downloadDir, fileName string) {
	response, err := http.Get(url)
	if err != nil {
		console.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(response.Body)

	err = os.MkdirAll(downloadDir, os.ModePerm)
	if err != nil {
		console.Fatal(err)
	}

	file, err := os.Create(filepath.Join(downloadDir, fileName))
	if err != nil {
		console.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	bar := progressbar.DefaultBytes(response.ContentLength, "Downloading "+fileName)

	_, err = io.Copy(io.MultiWriter(file, bar), response.Body)
	if err != nil {
		console.Fatal(err)
	}
}
