package build

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"code.google.com/p/go.crypto/openpgp"
	"github.com/kelseyhightower/pm/metadata"
)

var (
	metadataFilePath   string
	binaryFilePath     string
	privateKeyFilePath string
)

var flagset = flag.NewFlagSet("build", flag.ExitOnError)

func init() {
	flagset.StringVar(&binaryFilePath, "b", "", "binary file path")
	flagset.StringVar(&metadataFilePath, "m", "metadata.json", "metadata file path")
	flagset.StringVar(&privateKeyFilePath, "p", "private.key", "private key file path")
}

func Run() {
	flagset.Parse(os.Args[2:])
	metadataFile, err := os.Open(metadataFilePath)
	if err != nil {
		log.Fatal(err)
	}
	m, err := metadata.New(metadataFile)
	if err != nil {
		log.Fatal(err)
	}
	err = createPackage(binaryFilePath, m)
	if err != nil {
		log.Fatal(err)
	}
}

func createPackage(binaryFilePath string, m *metadata.Metadata) error {
	// Write the tar.gz file.
	name, err := createTarGz(binaryFilePath, m)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}
	// Write the .sha256 file.
	checksumFile, err := os.Create(name + ".sha256")
	if err != nil {
		return err
	}
	defer checksumFile.Close()
	log.Println(name + ".sha256")
	checksumContents := fmt.Sprintf("%x  %s\n", sha256.Sum256(data), name)
	checksumFile.WriteString(checksumContents)
	// write the checksum of the binary file.
	binaryFileData, err := ioutil.ReadFile(binaryFilePath)
	if err != nil {
		return err
	}
	checksumContents = fmt.Sprintf("%x  %s\n", sha256.Sum256(binaryFileData), filepath.Base(binaryFilePath))
	checksumFile.WriteString(checksumContents)

	// Setup the openpgp keyrign. Right now we only support non-encrypted private
	// keys.
	f, err := os.Open(privateKeyFilePath)
	if err != nil {
		return err
	}
	es, err := openpgp.ReadArmoredKeyRing(f)
	if err != nil {
		return err
	}
	// Sign the tar.gz file.
	detachedSigPackageFile, err := os.Create(name + ".asc")
	if err != nil {
		return err
	}
	defer detachedSigPackageFile.Close()
	log.Println(name + ".asc")
	err = openpgp.ArmoredDetachSign(detachedSigPackageFile, es[0], bytes.NewReader(data), nil)
	if err != nil {
		return err
	}
	// Sign the .sha256 file.
	detachedSigDigestFile, err := os.Create(name + ".sha256.asc")
	if err != nil {
		return err
	}
	defer detachedSigDigestFile.Close()
	err = openpgp.ArmoredDetachSign(detachedSigDigestFile, es[0], bytes.NewReader([]byte(checksumContents)), nil)
	if err != nil {
		return err
	}
	log.Println(name + ".sha256.asc")
	return nil
}

func createTarGz(binaryFilePath string, m *metadata.Metadata) (string, error) {
	outputFile, err := os.Create(formatPackageName(m))
	if err != nil {
		return "", err
	}
	defer outputFile.Close()
	gw := gzip.NewWriter(outputFile)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	// Save the metadata file
	metadataBytes, err := json.MarshalIndent(&m, "", "  ")
	if err != nil {
		return "", err
	}
	metadataHeader := &tar.Header{
		Name:       "metadata.json",
		Size:       int64(len(metadataBytes)),
		ModTime:    time.Now(),
		AccessTime: time.Now(),
		ChangeTime: time.Now(),
		Mode:       0444,
	}
	if err := tw.WriteHeader(metadataHeader); err != nil {
		return "", err
	}
	if _, err := tw.Write(metadataBytes); err != nil {
		return "", err
	}
	// Save the binary file.
	binaryFile, err := os.Open(binaryFilePath)
	if err != nil {
		return "", err
	}
	defer binaryFile.Close()
	info, err := binaryFile.Stat()
	if err != nil {
		return "", err
	}
	hdr, err := tar.FileInfoHeader(info, "")
	hdr.Uid = 0
	hdr.Gid = 0
	hdr.Mode = 0555
	if err != nil {
		return "", err
	}
	err = tw.WriteHeader(hdr)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(tw, binaryFile)
	if err != nil {
		return "", err
	}
	log.Println(outputFile.Name())
	return outputFile.Name(), nil
}

func formatPackageName(m *metadata.Metadata) string {
	return fmt.Sprintf("%s-%s-%s-%s.tar.gz", m.Name, m.Tag, m.Platform, m.Architecture)
}
