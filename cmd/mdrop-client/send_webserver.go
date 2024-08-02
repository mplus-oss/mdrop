package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/schollz/progressbar/v3"
)

var filePath string = ""

func SendWebserver(localPort int, file string) error {
    filePath = file

    http.Handle("/receive", http.HandlerFunc(receiveSendWebserver))
    return http.ListenAndServe(":"+strconv.Itoa(localPort), nil)
}

func receiveSendWebserver(w http.ResponseWriter, request *http.Request) {
    file, err := os.Open(filePath)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    defer file.Close()

    fileInfo, err := file.Stat()
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    w.Header().Set("Transfer-Encoding", "identity")
    w.Header().Set(
        "Content-Length",
        strconv.FormatInt(fileInfo.Size(), 10),
    )
    w.Header().Set("Content-Type", "binary/octet-stream")
    w.Header().Set(
        "Content-Disposition",
        fmt.Sprintf("attachment; filename=\"%v\"", fileInfo.Name()),
    )

    bar := progressbar.DefaultBytes(fileInfo.Size(), fileInfo.Name())
    _, err = io.Copy(io.MultiWriter(bar, w), file)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    request.Close = true
}
