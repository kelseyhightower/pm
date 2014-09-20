package metadata

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type Metadata struct {
	Architecture string `json:"architecture"`
	Description  string `json:"description"`
	Name         string `json:"name"`
	Maintainer   string `json:"maintainer"`
	Platform     string `json:"platform"`
	SouceUrl     string `json:"sourceUrl"`
	Tag          string `json:"tag"`
}

func New(r io.Reader) (*Metadata, error) {
	var m *Metadata
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
