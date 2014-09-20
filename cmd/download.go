package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

var cachePath = "/opt/pm/cache"

func DownloadPackage(rawUrl string) (string, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}
	fi, err := os.Stat(path.Join(cachePath, path.Base(u.Path)))
	if err == nil {
		localCopyTime := fi.ModTime()
		resp, err := http.Head(rawUrl)
		if err != nil {
			return "", err
		}
		lastModified := resp.Header.Get("Last-Modified")
		lastModifiedTime, err := time.Parse(time.RFC1123, lastModified)
		if err != nil {
			return "", err
		}
		if localCopyTime.After(lastModifiedTime) {
			log.Println("skipping download: local package newer than remote")
			return path.Join(cachePath, path.Base(u.Path)), nil
		}
	}
	packageName := u.Path
	files := []string{
		packageName,
		fmt.Sprintf("%s.sha256", packageName),
		fmt.Sprintf("%s.sha256.asc", packageName),
		fmt.Sprintf("%s.asc", packageName),
	}
	for _, f := range files {
		u.Path = f
		log.Printf("downloading %s", u)
		resp, err := http.Get(u.String())
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return "", errors.New("non 200 status code")
		}
		output, err := os.Create(path.Join(cachePath, path.Base(f)))
		if err != nil {
			return "", err
		}
		defer output.Close()
		if _, err := io.Copy(output, resp.Body); err != nil {
			return "", err
		}
	}
	return path.Join(cachePath, path.Base(packageName)), nil
}
