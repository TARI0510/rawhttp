package http

import (
	"bufio"
	"net/http"
	"net/url"
	"strings"
)

func readRequest(raw, scheme string) (*http.Request, error) {
	r, err := http.ReadRequest(bufio.NewReader(strings.NewReader(raw)))
	if err != nil {
		return nil, err
	}
	r.RequestURI, r.URL.Scheme, r.URL.Host = "", scheme, r.Host
	return r, nil
}

func ReadHTTPRequest(raw string, proxy string, excludeHeaders []string) *http.Response {
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

	for _, header := range excludeHeaders {
		for hk, _ := range req.Header {
			if strings.ToUpper(hk) == strings.ToUpper(header) {
				delete(req.Header, hk)
			}
		}
	}
	
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	return resp
}
