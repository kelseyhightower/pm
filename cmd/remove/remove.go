package remove

import (
	"log"
	"os"
	"path/filepath"
)

func Run() {
	packageName := os.Args[2]
	err := os.RemoveAll("/opt/" + packageName)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("/opt/" + packageName)
	files, err := filepath.Glob(filepath.Join("/opt/pm/cache", packageName+"*"))
	if err != nil {
		log.Fatal(err)
	}
	for _, name := range files {
		err := os.Remove(name)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(name)
	}
}
