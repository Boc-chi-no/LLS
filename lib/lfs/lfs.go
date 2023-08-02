package lfs

import (
	"linkshortener/lib/tool"
	"linkshortener/setting"
	"net/http"
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
		var indexName string
		if setting.Cfg.HTTP.DisableFilesDirEmbed {
			indexName = "index.html"
		} else {
			indexName = "/index.html"
		}

		index := tool.ConcatStrings(name, indexName)
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
