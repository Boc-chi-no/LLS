package lfs

import (
	"linkshortener/lib/tool"
	"net/http"
	"strings"
)

type LlsFileSystem struct {
	Fs http.FileSystem
}

func (lfs LlsFileSystem) Open(name string) (http.File, error) {
	f, err := lfs.Fs.Open(name)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		index := tool.ConcatStrings(strings.TrimSuffix(name, "/"), "/index.html")
		if _, err := lfs.Fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}
	return f, nil
}
