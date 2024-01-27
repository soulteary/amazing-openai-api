package fn

import (
	"compress/gzip"
	"io"
)

func Gunzip(r io.Reader) (reader io.ReadCloser, err error) {
	reader, err = gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return reader, nil
}
