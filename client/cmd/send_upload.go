package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func Upload(url string, values map[string]io.Reader) (err error) {
    var b bytes.Buffer
	client := &http.Client{}
    w := multipart.NewWriter(&b)
    for key, r := range values {
        var fw io.Writer
        if x, ok := r.(io.Closer); ok {
            defer x.Close()
        }
        if x, ok := r.(*os.File); ok {
            if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
                return
            }
        } else {
            if fw, err = w.CreateFormField(key); err != nil {
                return
            }
        }
        if _, err = io.Copy(fw, r); err != nil {
            return err
        }

    }
    w.Close()

    req, err := http.NewRequest("POST", url, &b)
    if err != nil {
        return
    }
    req.Header.Set("Content-Type", w.FormDataContentType())
    res, err := client.Do(req)
    if err != nil {
        return
    }

    if res.StatusCode != http.StatusOK {
        err = fmt.Errorf("Bad status: %s", res.Status)
    }
    return
}

func mustOpen(f string) *os.File {
    r, err := os.Open(f)
    if err != nil {
        panic(err)
    }
    return r
}