package filetar

import (
	"archive/zip"
	"fmt"
	"io"
	"os"

	"gopl.io/ch10/file/extract"
)

type zipExtractor struct {
	io.Reader
	zipReader *zip.Reader
	i         int
}

func newReader(file *os.File) (extract.ArchiveReader, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("read file %s: %v", fileInfo.Name(), err)
	}
	zipReader, err := zip.NewReader(file, fileInfo.Size())
	if err != nil {
		return nil, fmt.Errorf("read file %s: %v", fileInfo.Name(), err)
	}
	return &zipExtractor{zipReader: zipReader, i: -1}, nil
}

func (e *zipExtractor) Next() error {
	e.i++
	if e.i >= len(e.zipReader.File) {
		return io.EOF
	}
	reader, err := e.zipReader.File[e.i].Open()
	if err != nil {
		return fmt.Errorf("next: %v", err)
	}
	e.Reader = reader

	return nil
}

func init() {
	sign := extract.Sign{ContentType: "application/zip"}
	extract.RegisterReader(sign, newReader)
}
