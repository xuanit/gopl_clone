package filetar

import (
	"archive/tar"
	"io"
	"os"

	"gopl.io/ch10/file/extract"
)

type tarExtractor struct {
	io.Reader
}

func newReader(f *os.File) (extract.ArchiveReader, error) {
	return tarExtractor{Reader: tar.NewReader(f)}, nil
}

func (e tarExtractor) Next() error {
	_, err := e.Reader.(*tar.Reader).Next()
	return err
}

func init() {
	sign := extract.Sign{ContentType: "application/octet-stream", MagicNumber: "ustar"}
	extract.RegisterReader(sign, newReader)
}
