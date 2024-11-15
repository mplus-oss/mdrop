package main

import (
	"os"
	"path/filepath"
)

func ParseGlob(args []string) ([]string, error) {
    files := []string{}
    for _, arg := range args {
        globMatch, err := filepath.Glob(arg)
        if err != nil {
            return nil, err
        }

        for _, globFile := range globMatch {
            pwd, err := os.Getwd()
            if err != nil {
                return nil, err
            }
            fileInfo, err := os.Stat(filepath.Join(pwd, globFile))
            if err != nil {
                return nil, err
            }
            if fileInfo.IsDir() {
                continue
            }
            files = append(files, globFile)
        }
    }
    return files, nil
}
