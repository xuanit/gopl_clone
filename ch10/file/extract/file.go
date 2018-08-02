package extract

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type ArchiveReader interface {
	io.Reader
	Next() error
}

type NewReaderFunc func(f *os.File) (ArchiveReader, error)

type Sign struct {
	ContentType string
	MagicNumber string
}

var readerCreators map[Sign]NewReaderFunc

func init() {
	readerCreators = make(map[Sign]NewReaderFunc)
}

func RegisterReader(s Sign, newReader NewReaderFunc) (err error) {
	if _, ok := readerCreators[s]; ok {
		err = fmt.Errorf("sign has been existing")
	}
	readerCreators[s] = newReader
	return
}

func Read(name string, handleFile func(reader io.Reader)) (err error) {
	file, err := os.Open(name)
	if err != nil {
		return fmt.Errorf("read file %s: %v", name, err)
	}
	defer file.Close()

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		return fmt.Errorf("read archive: %v", err)
	}
	file.Seek(0, 0)

	fileType := http.DetectContentType(buff)
	fmt.Println(fileType)
	s := Sign{ContentType: fileType, MagicNumber: string(buff[257:262])}
	shorthand := Sign{ContentType: fileType}
	newReader, ok := readerCreators[s]
	if !ok {
		newReader, ok = readerCreators[shorthand]
		if !ok {
			return fmt.Errorf("invalid format")
		}
	}

	archiveReader, err := newReader(file)
	if err != nil {
		return fmt.Errorf("read archive: %v", err)
	}

	for {
		err := archiveReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read file %s: %v", name, err)
		}
		handleFile(archiveReader)
	}

	return nil
}
