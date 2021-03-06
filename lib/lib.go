package lib

import (
	"fmt"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

/*
You can compile this on your Mac, then, copy to Linux

export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
go build -o url url.go


*/

func ProxyPullData(URL string) ([]byte, error) {
	PROXY_ADDR := "127.0.0.1:1337"
	dialer, err := proxy.SOCKS5("tcp", PROXY_ADDR, nil, &net.Dialer{
		Timeout:   60 * time.Second,
		KeepAlive: 30 * time.Second,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
		return nil, err
	}

	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}

	httpTransport.Dial = dialer.Dial

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't create request:", err)
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't GET page:", err)
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	return b, err
}

func NewClientBindedToIP(ip string) (*http.Client, error) {
	if trans, err := NewTransportBindedToIP(ip); err != nil {
		return nil, err
	} else {
		return &http.Client{Transport: trans}, nil
	}
}

func NewTransportBindedToIP(ip string) (*http.Transport, error) {
	ipAddr, err := net.ResolveIPAddr("ip", ip)
	if err != nil {
		return nil, err
	}

	trans := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			LocalAddr: &net.TCPAddr{IP: ipAddr.IP},
			Timeout:   3 * time.Second,
			KeepAlive: 3 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       9 * time.Second,
		TLSHandshakeTimeout:   4 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return trans, nil
}

type NetworkTransport struct {
	n *http.Transport
	c *http.Client
}

func InitNT() NetworkTransport {
	n := NetworkTransport{}
	var netTransport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 5 * time.Second,
		}).DialContext,
		MaxIdleConns:          10,
		IdleConnTimeout:       4 * time.Second,
		TLSHandshakeTimeout:   4 * time.Second,
		ExpectContinueTimeout: 4 * time.Second,
	}

	var netClient = &http.Client{
		Timeout:   time.Second * 1,
		Transport: netTransport,
	}

	n.n = netTransport
	n.c = netClient
	return n

}

func (n *NetworkTransport) Get(url string) ([]byte, error) {

	resp, err := n.c.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err

}

func (n *NetworkTransport) Process(records []string) {

	var wg sync.WaitGroup
	for i, record := range records {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			output, err := n.Get(url)
			if err != nil {
				fmt.Println(string(output), err)
			}

		}(record)
		fmt.Println(i, record)
	}
	wg.Wait()

}

func ReadFile(file string) ([]string, error) {
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	records := strings.Split(string(dat), "\n")
	return records, err
}
