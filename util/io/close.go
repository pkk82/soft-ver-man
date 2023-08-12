package io

import (
	"github.com/pkk82/soft-ver-man/util/console"
	"io"
)

func CloseOrLog(closer io.Closer) {
	err := closer.Close()
	if err != nil {
		console.Error(err)
	}
}
