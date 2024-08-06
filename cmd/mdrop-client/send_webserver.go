package main

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mplus-oss/mdrop/internal"
	"github.com/schollz/progressbar/v3"
)

var server *http.Server = &http.Server{}
var filePath string = ""
var isStillUsed bool = false

var senderErrorChan chan error = make(chan error)

func SendWebserver(localPort int, file string) (err error) {
	filePath = file
	server.Addr = ":"+strconv.Itoa(localPort)

	http.Handle("/receive", http.HandlerFunc(receiveSendWebserver))
	http.Handle("/checksum", http.HandlerFunc(checksumSendWebserver))

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			if strings.Contains(err.Error(), "Server closed") {
				return
			}
			senderErrorChan <- internal.CustomizeError("receiveWebserverFatal", err)
		}
	}()

	err = <-senderErrorChan

    shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
    defer shutdownRelease()

	fmt.Println("Gracefully shutdown server...")
	err = server.Shutdown(shutdownCtx)
    if err != nil {
		return err
    }
	return nil
}

func checksumSendWebserver(w http.ResponseWriter, request *http.Request) {
	fmt.Println("Receiver taking the checksum file.")
	file, err := os.Open(filePath)
	if err != nil {
		senderErrorChan <- internal.CustomizeError("checksumOpenFile", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		senderErrorChan <- internal.CustomizeError("checksumHashSum", err)
	}

	hashString := fmt.Sprintf("%x", hash.Sum(nil))
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, hashString)

	request.Close = true
}

func receiveSendWebserver(w http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if isStillUsed {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// This prevent from MITM after transfering file
	isStillUsed = true

	file, err := os.Open(filePath)
	if err != nil {
		senderErrorChan <- internal.CustomizeError("receiveOpenFile", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		senderErrorChan <- internal.CustomizeError("receiveOpenFileStat", err)
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
	w.Header().Set("X-Attachment-Name", fileInfo.Name())

	bar := progressbar.DefaultBytes(fileInfo.Size(), fileInfo.Name())
	_, err = io.Copy(io.MultiWriter(bar, w), file)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "broken pipe") {
			err = errors.New("Broken pipe from receiver because forced close or terminated.")
		}
		request.Close = true
		senderErrorChan <- internal.CustomizeError("receiveStreamFile", err)
	}

	request.Close = true
	senderErrorChan <- nil
}
