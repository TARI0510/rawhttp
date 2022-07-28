package http

import (
	"net/http"
	"net/url"
	"strings"
)

func CdnReadHTTPRequest(raw, proxy string) *http.Response {
	req, err := readRequest(raw, "http")
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	if proxy != "" {
		proxyUri, e := url.Parse(proxy)
		if e != nil {
			panic(e)
		}
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyUri),
		}
	}

	for hk, _ := range req.Header {
		if strings.HasPrefix(strings.ToUpper(hk), "X-") {
			delete(req.Header, hk)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	return resp
}
