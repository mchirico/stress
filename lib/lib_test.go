package lib

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNetSwitching(t *testing.T) {
	local_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer local_server.Close()

	client, err := NewClientBindedToIP("127.0.0.1")
	if err != nil {
		t.Fail()
	}

	if _, err := client.Get("http://google.com"); err == nil {
		t.Fail()
	}

	if _, err := client.Get("https://google.com"); err == nil {
		t.Fail()
	}

	if _, err := client.Get(local_server.URL); err != nil {
		t.Fail()
	}
}

func TestReadFile(t *testing.T) {
	file := "../fixtures/list_of_urls"
	r, err := ReadFile(file)
	if len(r) != 3 {
		t.Fatalf("Can't read file:%s\n", err)
	}

}

func TestGet(t *testing.T) {
	file := "../fixtures/list_of_urls"
	records, _ := ReadFile(file)
	for _, record := range records {
		_, err := Get(record)
		if err != nil {
			t.Fatalf("Can't make url call")
		}
	}

	Process(records)
}
