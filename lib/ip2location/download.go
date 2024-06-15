package ip2location

import (
	"io"
	"net/http"
)

func GetOnline() ([]byte, error) {
	resp, err := http.Get("https://git.io/GeoLite2-City.mmdb")
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	return io.ReadAll(resp.Body)
}
