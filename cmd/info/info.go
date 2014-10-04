package info

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/kelseyhightower/pm/metadata"
)

func Run() {
	packageName := os.Args[2]
	f, err := os.Open(fmt.Sprintf("/opt/%s/metadata.json", packageName))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	m, err := metadata.New(f)
	if err != nil {
		log.Fatal(err)
	}
	data, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))
}
