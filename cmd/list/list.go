package list

import (
	"log"
	"path/filepath"
	"regexp"
)

var pattern = regexp.MustCompile("/opt/(.*)/metadata.json")

func Run() {
	paths, err := filepath.Glob("/opt/*/metadata.json")
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range paths {
		log.Println(pattern.FindStringSubmatch(p)[1])
	}
}
