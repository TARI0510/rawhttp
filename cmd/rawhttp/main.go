package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/TARI0510/rawhttp/pkg/http"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

var (
	proxy          string
	excludeHeaders string
	isQuiet        bool
	version        bool
)

func main() {
	flag.StringVar(&proxy, "p", "", "[optional] http proxy, example: http://127.0.0.1:8080")
	flag.StringVar(&excludeHeaders, "eh", "", "[optional] exclude header, separator by comma, example: X-From,X-other")
	flag.BoolVar(&isQuiet, "q", false, "[optional] quiet mode, doesn't output http response on screen")
	flag.BoolVar(&version, "v", false, "show version")
	if version {
		fmt.Println(http.Version)
		os.Exit(0)
	}
	flag.Parse()

	var httpRaw []byte

	httpRawReader := bufio.NewReader(os.Stdin)
	for {
		line, err := httpRawReader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if bytes.Equal(line, []byte("\n")) {
			break
		}
		httpRaw = append(httpRaw, line...)
	}

	raw := strings.Replace(string(httpRaw), "\n", "\r\n", -1)
	raw = raw + "\r\n\r\n"

	var resp http.Response
	excludeHeaderArr := strings.Split(excludeHeaders, ",")
	resp = http.ReadHTTPRequest(raw, proxy, excludeHeaderArr)
	fmt.Println()

	if isQuiet {
		return
	}
	resBytes, _ := ioutil.ReadAll(resp.Body)

	fmt.Printf("%s %s\n", resp.Proto, resp.Status)
	for hk, hvs := range resp.Header {
		fmt.Printf("%s: %s\n", hk, strings.Join(hvs, ","))
	}
	fmt.Printf("\n%s\n", string(resBytes))
}
