package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/mplus-oss/mdrop/internal"
	"github.com/schollz/progressbar/v3"
)

var filePath string = ""

func SendWebserver(localPort int, file string) error {
    filePath = file

    http.Handle("/receive", http.HandlerFunc(receiveSendWebserver))
    http.Handle("/checksum", http.HandlerFunc(checksumSendWebserver))

    return http.ListenAndServe(":"+strconv.Itoa(localPort), nil)
}

func checksumSendWebserver(w http.ResponseWriter, request *http.Request) {
    file, err := os.Open(filePath)
    if err != nil {
        internal.PrintErrorWithExit("checksumOpenFile", err, 1)
    }
    defer file.Close()

    hash := sha256.New()
    if _, err := io.Copy(hash, file); err != nil {
        internal.PrintErrorWithExit("checksumHashSum", err, 1)
    }

    w.Header().Set("Content-Type", "text/plain")
    fmt.Fprint(w, hex.EncodeToString(hash.Sum(nil)))

    request.Close = true
}

func receiveSendWebserver(w http.ResponseWriter, request *http.Request) {
    if request.Method != "POST" {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    file, err := os.Open(filePath)
    if err != nil {
        internal.PrintErrorWithExit("receiveOpenFile", err, 1)
    }
    defer file.Close()

    fileInfo, err := file.Stat()
    if err != nil {
        internal.PrintErrorWithExit("receiveOpenFileStat", err, 1)
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
        errMsg := err.Error()
        if strings.Contains(errMsg, "broken pipe") {
            err = errors.New("Broken pipe from receiver because forced close or terminated.")
        }
        internal.PrintErrorWithExit("receiveStreamFile", err, 1)
    }

    request.Close = true
}
