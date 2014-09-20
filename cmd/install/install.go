package install

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/kelseyhightower/pm/cmd"
	"github.com/kelseyhightower/pm/metadata"
	"github.com/kelseyhightower/pm/openpgp"
)

func Run() {
	packagePath := os.Args[2]
	if !strings.HasSuffix(packagePath, ".tar.gz") {
		log.Fatal("invalid package name: ", packagePath)
	}
	u, err := url.Parse(packagePath)
	if err != nil {
		log.Fatal(err)
	}
	if u.Scheme != "" {
		packagePath, err = cmd.DownloadPackage(packagePath)
		if err != nil {
			log.Fatal(err)
		}
	}
	if !openpgp.Verify(packagePath) {
		log.Fatal("not valid")
	}
	packageFile, err := os.Open(packagePath)
	if err != nil {
		log.Fatal(err)
	}
	gzipReader, err := gzip.NewReader(packageFile)
	if err != nil {
		log.Fatal(err)
	}
	tarReader := tar.NewReader(gzipReader)
	var m *metadata.Metadata
	var baseDir string
	for {
		hdr, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		if hdr.Name == "metadata.json" {
			data, err := ioutil.ReadAll(tarReader)
			if err != nil {
				log.Fatal(err)
			}
			m, err = metadata.New(bytes.NewReader(data))
			if err != nil {
				log.Fatalln(err)
			}
			baseDir = fmt.Sprintf("/opt/%s-%s", m.Name, m.Tag)
			err = os.MkdirAll(baseDir, 0755)
			if err != nil {
				log.Fatalln(err)
			}
			metadataPath := path.Join(baseDir, "metadata.json")
			f, err := os.Create(metadataPath)
			if err != nil {
				log.Fatalln(err)
			}
			defer f.Close()
			if _, err := io.Copy(f, bytes.NewReader(data)); err != nil {
				log.Fatalln(err)
			}
			err = os.Chmod(metadataPath, 0644)
			if err != nil {
				log.Fatalln(err)
			}
			log.Println(metadataPath)
			continue
		}
		if hdr.Name == m.Name {
			binDir := path.Join(baseDir, "/bin")
			binPath := path.Join(binDir, m.Name)
			err := os.MkdirAll(binDir, 0755)
			if err != nil {
				log.Fatalln(err)
			}
			f, err := os.Create(binPath)
			if err != nil {
				log.Fatalln(err)
			}
			defer f.Close()
			if _, err := io.Copy(f, tarReader); err != nil {
				log.Fatalln(err)
			}
			err = os.Chmod(binPath, 0555)
			if err != nil {
				log.Fatalln(err)
			}
			log.Println(binPath)
		}
	}
}
